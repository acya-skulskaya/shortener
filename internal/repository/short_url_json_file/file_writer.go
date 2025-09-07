package shorturljsonfile

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/acya-skulskaya/shortener/internal/model"
	"os"
)

type FileWriter struct {
	file *os.File
	// добавляем Writer в FileWriter
	writer *bufio.Writer
}

func NewFileWriter(filename string) (*FileWriter, error) {
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}

	return &FileWriter{
		file: file,
		// создаём новый Writer
		writer: bufio.NewWriter(file),
	}, nil
}

func (p *FileWriter) WriteFile(row model.URLList) error {
	fileReader, err := NewFileReader(p.file.Name())
	if err != nil {
		return err
	}
	defer fileReader.Close()

	list, err := fileReader.ReadFile()
	if err != nil {
		fmt.Println(err)
		return err
	}

	list = append(list, row)

	data, _ := json.Marshal(list)

	// записываем событие в буфер
	if _, err := p.writer.Write(data); err != nil {
		return err
	}

	// записываем буфер в файл
	return p.writer.Flush()
}
