package shorturljsonfile

import (
	"bufio"
	"encoding/json"
	"fmt"
	jsonModel "github.com/acya-skulskaya/shortener/internal/model/json"
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
		return nil, fmt.Errorf("error opening file %s: %w", filename, err)
	}

	return &FileReader{
		file: file,
		// создаём новый scanner
		scanner: bufio.NewScanner(file),
	}, nil
}

func (c *FileReader) ReadFile() (list []jsonModel.URLList, ids []string, error error) {
	// одиночное сканирование до следующей строки
	if !c.scanner.Scan() {
		return nil, nil, c.scanner.Err()
	}
	// читаем данные из scanner
	data := c.scanner.Bytes()

	err := json.Unmarshal(data, &list)
	if err != nil {
		return nil, nil, fmt.Errorf("could not read json: %w", err)
	}
	for _, item := range list {
		ids = append(ids, item.ID)
	}

	return list, ids, nil
}

func (c *FileReader) Close() error {
	return c.file.Close()
}
