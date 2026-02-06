package shorturljsonfile

import (
	"encoding/json"
	"fmt"

	"github.com/acya-skulskaya/shortener/internal/helpers"
	jsonModel "github.com/acya-skulskaya/shortener/internal/model/json"
)

type FileWriter struct {
	filename string
}

func NewFileWriter(filename string) *FileWriter {
	return &FileWriter{
		filename: filename,
	}
}

func (p *FileWriter) WriteFile(row jsonModel.URLList) error {
	fileReader := NewFileReader(p.filename)

	list, err := fileReader.ReadFile()
	if err != nil {
		return fmt.Errorf("could not read file: %w", err)
	}
	list = append(list, row)

	data, err := json.Marshal(list)
	if err != nil {
		return fmt.Errorf("could not encode json: %w", err)
	}

	if err = helpers.Write(p.filename, data); err != nil {
		return fmt.Errorf("could not write file: %w", err)
	}

	return nil
}

func (p *FileWriter) WriteFileRows(rows []jsonModel.URLList) error {
	fileReader := NewFileReader(p.filename)

	list, err := fileReader.ReadFile()
	if err != nil {
		return fmt.Errorf("could not read file: %w", err)
	}

	list = append(list, rows...)

	data, err := json.Marshal(list)
	if err != nil {
		return fmt.Errorf("could not encode json: %w", err)
	}

	if err = helpers.Write(p.filename, data); err != nil {
		return fmt.Errorf("could not write file rows: %w", err)
	}

	return nil
}

func (p *FileWriter) OverwriteFile(rows []jsonModel.URLList) error {
	data, err := json.Marshal(rows)
	if err != nil {
		return fmt.Errorf("could not encode json: %w", err)
	}

	if err = helpers.Write(p.filename, data); err != nil {
		return fmt.Errorf("could not overwrite file: %w", err)
	}

	return helpers.Write(p.filename, data)
}
