package shorturljsonfile

import (
	"bufio"
	"encoding/json"
	"github.com/acya-skulskaya/shortener/internal/model"
	"os"
)

type FileReader struct {
	file *os.File
	// заменяем FileReader на Scanner
	scanner *bufio.Scanner
}

func NewFileReader(filename string) (*FileReader, error) {
	file, err := os.OpenFile(filename, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}

	return &FileReader{
		file: file,
		// создаём новый scanner
		scanner: bufio.NewScanner(file),
	}, nil
}

func (c *FileReader) ReadFile() ([]model.URLList, error) {
	// одиночное сканирование до следующей строки
	if !c.scanner.Scan() {
		return nil, c.scanner.Err()
	}
	// читаем данные из scanner
	data := c.scanner.Bytes()

	var list []model.URLList
	err := json.Unmarshal(data, &list)
	if err != nil {
		return nil, err
	}

	return list, nil
}

func (c *FileReader) Close() error {
	return c.file.Close()
}
