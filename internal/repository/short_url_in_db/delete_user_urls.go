package shorturlindb

import (
	"context"
	"database/sql"
	"errors"
	errorsInternal "github.com/acya-skulskaya/shortener/internal/errors"
	"github.com/acya-skulskaya/shortener/internal/logger"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
	"sync"
)

func checkIDsInDB(db *sql.DB, ctx context.Context, list []string, userID string) (err error) {
	stmt, err := db.PrepareContext(ctx,
		"SELECT user_id, is_deleted FROM short_urls WHERE id = $1")
	if err != nil {
		logger.Log.Debug("could not prepare statement",
			zap.Error(err),
		)
		return err
	}
	defer stmt.Close()

	for _, id := range list {
		row := stmt.QueryRowContext(ctx, id)

		var rowUserID string
		var rowIsDeleted bool
		err = row.Scan(&rowUserID, &rowIsDeleted) // разбираем результат
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return errorsInternal.ErrIDNotFound
			}

			logger.Log.Debug("could not get row from db",
				zap.Error(err),
			)
			return err
		}

		if rowUserID != userID {
			return errorsInternal.ErrUserIDUnauthorized
		}
		if rowIsDeleted {
			return errorsInternal.ErrIDDeleted
		}
	}

	return nil
}

func (repo *InDBShortURLRepository) DeleteUserUrls(ctx context.Context, list []string, userID string) (err error) {
	err = checkIDsInDB(repo.DB, ctx, list, userID)
	if err != nil {
		return err
	}

	tx, err := repo.DB.Begin()
	if err != nil {
		logger.Log.Debug("could not start transaction",
			zap.Error(err),
		)
		return err
	}

	stmt, err := tx.PrepareContext(ctx,
		"UPDATE short_urls SET is_deleted = true WHERE id = $1")
	if err != nil {
		logger.Log.Debug("could not prepare statement",
			zap.Error(err),
		)
		return err
	}
	defer stmt.Close()

	// сигнальный канал для завершения горутин
	doneChannel := make(chan struct{})
	// закрываем его при завершении программы
	defer close(doneChannel)

	// канал с данными
	inputChannel := generator(doneChannel, list)

	// получаем слайс каналов из 10 рабочих add
	channels := fanOut(doneChannel, inputChannel, ctx, stmt)

	// а теперь объединяем десять каналов в один
	resultChannel := fanIn(doneChannel, channels...)

	// выводим результаты расчетов из канала
	for err = range resultChannel {
		if err != nil {
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

// generator функция из предыдущего примера, делает то же, что и делала
func generator(doneCh chan struct{}, list []string) chan string {
	inputCh := make(chan string)

	go func() {
		defer close(inputCh)

		for _, data := range list {
			select {
			case <-doneCh:
				return
			case inputCh <- data:
			}
		}
	}()

	return inputCh
}

// fanIn объединяет несколько каналов resultChs в один.
func fanIn(doneCh chan struct{}, resultChs ...chan error) chan error {
	// конечный выходной канал в который отправляем данные из всех каналов из слайса, назовём его результирующим
	finalCh := make(chan error)

	// понадобится для ожидания всех горутин
	var wg sync.WaitGroup

	// перебираем все входящие каналы
	for _, ch := range resultChs {
		// в горутину передавать переменную цикла нельзя, поэтому делаем так
		chClosure := ch

		// инкрементируем счётчик горутин, которые нужно подождать
		wg.Add(1)

		go func() {
			// откладываем сообщение о том, что горутина завершилась
			defer wg.Done()

			// получаем данные из канала
			for data := range chClosure {
				select {
				// выходим из горутины, если канал закрылся
				case <-doneCh:
					return
				// если не закрылся, отправляем данные в конечный выходной канал
				case finalCh <- data:
				}
			}
		}()
	}

	go func() {
		// ждём завершения всех горутин
		wg.Wait()
		// когда все горутины завершились, закрываем результирующий канал
		close(finalCh)
	}()

	// возвращаем результирующий канал
	return finalCh
}

// fanOut принимает канал данных, порождает 10 горутин
func fanOut(doneCh chan struct{}, inputCh chan string, ctx context.Context, stmt *sql.Stmt) []chan error {
	// количество горутин add
	numWorkers := 10
	// каналы, в которые отправляются результаты
	channels := make([]chan error, numWorkers)

	for i := 0; i < numWorkers; i++ {
		// получаем канал из горутины add
		updateResultCh := update(doneCh, inputCh, ctx, stmt)
		// отправляем его в слайс каналов
		channels[i] = updateResultCh
	}

	// возвращаем слайс каналов
	return channels
}

func update(doneCh chan struct{}, inputCh chan string, ctx context.Context, stmt *sql.Stmt) chan error {
	updateErr := make(chan error)

	go func() {
		defer close(updateErr)

		for id := range inputCh {
			_, err := stmt.ExecContext(ctx, id)

			select {
			case <-doneCh:
				return
			case updateErr <- err:
			}
		}
	}()

	return updateErr
}
