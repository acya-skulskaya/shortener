package shorturljsonfile

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"

	jsonModel "github.com/acya-skulskaya/shortener/internal/model/json"
)

type FileWriter struct {
	file *os.File
	// добавляем Writer в FileWriter
	writer *bufio.Writer
}

func NewFileWriter(filename string) (*FileWriter, error) {
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return nil, fmt.Errorf("error opening file %s: %w", filename, err)
	}

	return &FileWriter{
		file: file,
		// создаём новый Writer
		writer: bufio.NewWriter(file),
	}, nil
}

func (p *FileWriter) WriteFile(row jsonModel.URLList) error {
	fileReader, err := NewFileReader(p.file.Name())
	if err != nil {
		return err
	}
	defer fileReader.Close()

	list, err := fileReader.ReadFile()
	if err != nil {
		return err
	}
	list = append(list, row)

	data, err := json.Marshal(list)
	if err != nil {
		return fmt.Errorf("could not encode json: %w", err)
	}

	// записываем событие в буфер
	if _, err := p.writer.Write(data); err != nil {
		return err
	}

	// записываем буфер в файл
	return p.writer.Flush()
}

func (p *FileWriter) WriteFileRows(rows []jsonModel.URLList) error {
	fileReader, err := NewFileReader(p.file.Name())
	if err != nil {
		return err
	}
	defer fileReader.Close()

	list, err := fileReader.ReadFile()
	if err != nil {
		return err
	}

	list = append(list, rows...)

	data, err := json.Marshal(list)
	if err != nil {
		return fmt.Errorf("could not encode json: %w", err)
	}

	// записываем событие в буфер
	if _, err := p.writer.Write(data); err != nil {
		return err
	}

	// записываем буфер в файл
	return p.writer.Flush()
}

func (p *FileWriter) OverwriteFile(rows []jsonModel.URLList) error {
	data, err := json.Marshal(rows)
	if err != nil {
		return fmt.Errorf("could not encode json: %w", err)
	}

	// записываем событие в буфер
	if _, err := p.writer.Write(data); err != nil {
		return err
	}

	// записываем буфер в файл
	return p.writer.Flush()
}
