package shorturlindb

import (
	"context"
	"database/sql"
	"sync"
	"time"

	"github.com/acya-skulskaya/shortener/internal/logger"
	jsonModel "github.com/acya-skulskaya/shortener/internal/model/json"
	"go.uber.org/zap"
)

func (repo *InDBShortURLRepository) DeleteUserUrls(ctx context.Context, list []string, userID string) {
	inputCh := make(chan jsonModel.URLList, len(list))
	go func() {
		defer close(inputCh)
		for _, shortID := range list {
			select {
			case <-ctx.Done():
				return
			case inputCh <- jsonModel.URLList{ID: shortID, UserID: userID}:
			}
		}
	}()

	batchChannels := fanOut(ctx, inputCh)
	batchesCh := fanIn(ctx, batchChannels...)
	deleteBatch(ctx, repo.DB, batchesCh)
}

func fanOut(ctx context.Context, inputCh chan jsonModel.URLList) []chan jsonModel.URLListBatch {
	const numWorkers = 4
	channels := make([]chan jsonModel.URLListBatch, numWorkers)
	batchSize := 50
	for i := 0; i < numWorkers; i++ {
		outCh := batchWorker(ctx, inputCh, batchSize)
		channels[i] = outCh
	}
	return channels
}

func batchWorker(
	ctx context.Context,
	inCh chan jsonModel.URLList,
	batchSize int,
) chan jsonModel.URLListBatch {
	outCh := make(chan jsonModel.URLListBatch)

	go func() {
		defer close(outCh)
		batch := make(jsonModel.URLListBatch, 0, batchSize)
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				flush(ctx, outCh, &batch, batchSize)
				return
			case urlToDelete, ok := <-inCh:
				if !ok {
					flush(ctx, outCh, &batch, batchSize)
					return
				}
				batch = append(batch, urlToDelete)
				if len(batch) >= batchSize {
					flush(ctx, outCh, &batch, batchSize)
				}
			case <-ticker.C:
				flush(ctx, outCh, &batch, batchSize)
			}
		}
	}()

	return outCh
}

func flush(
	ctx context.Context,
	outputCh chan jsonModel.URLListBatch,
	batch *jsonModel.URLListBatch,
	batchSize int,
) {
	if len(*batch) > 0 {
		select {
		case <-ctx.Done():
		case outputCh <- *batch:
			*batch = make(jsonModel.URLListBatch, 0, batchSize)
		}
	}
}

func fanIn(ctx context.Context, channels ...chan jsonModel.URLListBatch) chan jsonModel.URLListBatch {
	finalCh := make(chan jsonModel.URLListBatch)
	var wg sync.WaitGroup

	for _, ch := range channels {
		wg.Add(1)
		go func(inCh chan jsonModel.URLListBatch) {
			defer wg.Done()

			for batch := range inCh {
				select {
				case <-ctx.Done():
					return
				case finalCh <- batch:
				}
			}
		}(ch)
	}

	go func() {
		wg.Wait()
		close(finalCh)
	}()

	return finalCh
}

func deleteBatch(ctx context.Context, db *sql.DB, batchesCh chan jsonModel.URLListBatch) {
	for batch := range batchesCh {
		func(batch jsonModel.URLListBatch) {
			ids := make([]string, len(batch))
			idsToUser := make(map[string]string)
			for i, j := range batch {
				ids[i] = j.ID
				idsToUser[j.ID] = j.UserID
			}

			tx, err := db.Begin()
			if err != nil {
				logger.Log.Debug("could not start transaction",
					zap.Error(err),
				)
				return
			}
			defer tx.Rollback()

			rows, err := db.QueryContext(ctx, "SELECT id, user_id from short_urls WHERE id = ANY($1) AND is_deleted = false FOR UPDATE", ids)
			if err != nil {
				logger.Log.Debug("could not query from db", zap.Error(err))
				return
			}
			defer rows.Close()

			var idsToDelete []string
			for rows.Next() {
				shortURL := jsonModel.URLList{}
				err = rows.Scan(&shortURL.ID, &shortURL.UserID)
				if err != nil {
					logger.Log.Debug("could not scan row", zap.Error(err))
					return
				}

				if idsToUser[shortURL.ID] == shortURL.UserID {
					idsToDelete = append(idsToDelete, shortURL.ID)
				}
			}

			err = rows.Err()
			if err != nil {
				logger.Log.Debug("error getting user urls", zap.Error(err))
				return
			}

			if len(idsToDelete) == 0 {
				return
			}

			_, err = tx.ExecContext(ctx,
				"UPDATE short_urls SET is_deleted = true WHERE id = ANY($1)", idsToDelete)
			if err != nil {
				logger.Log.Debug("could not execute update query",
					zap.Error(err),
				)
				return
			}

			err = tx.Commit()
			if err != nil {
				logger.Log.Debug("could not commit transaction",
					zap.Error(err),
				)
			}

			logger.Log.Info("batch was deleted",
				zap.Error(err),
			)
		}(batch)
	}
}
