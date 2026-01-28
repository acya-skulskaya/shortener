package shorturljsonfile

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"

	jsonModel "github.com/acya-skulskaya/shortener/internal/model/json"
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

func (c *FileReader) ReadFile() (list []jsonModel.URLList, error error) {
	// одиночное сканирование до следующей строки
	if !c.scanner.Scan() {
		return nil, c.scanner.Err()
	}
	// читаем данные из scanner
	data := c.scanner.Bytes()

	err := json.Unmarshal(data, &list)
	if err != nil {
		return nil, fmt.Errorf("could not read json: %w", err)
	}

	return list, nil
}

func (c *FileReader) Close() error {
	return c.file.Close()
}
