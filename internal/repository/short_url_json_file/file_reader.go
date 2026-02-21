package shorturljsonfile

import (
	"encoding/json"
	"fmt"

	"github.com/acya-skulskaya/shortener/internal/helpers"
	jsonModel "github.com/acya-skulskaya/shortener/internal/model/json"
)

type FileReader struct {
	filename string
}

func NewFileReader(filename string) *FileReader {
	return &FileReader{
		filename: filename,
	}
}

func (c *FileReader) ReadFile() (list []jsonModel.URLList, error error) {
	data, err := helpers.Scan(c.filename)
	if err != nil {
		return nil, fmt.Errorf("could not scan file %s: %w", c.filename, err)
	}

	if len(data) == 0 {
		return nil, nil
	}

	err = json.Unmarshal(data, &list)
	if err != nil {
		return nil, fmt.Errorf("could not read json: %w", err)
	}

	return list, nil
}
