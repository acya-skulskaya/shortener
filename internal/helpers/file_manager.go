package helpers

import (
	"bufio"
	"fmt"
	"os"
)

func Scan(filename string) ([]byte, error) {
	file, err := os.OpenFile(filename, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return nil, fmt.Errorf("error opening file %s: %w", filename, err)
	}
	defer file.Close()

	fileinfo, _ := file.Stat()
	bs := make([]byte, fileinfo.Size())

	_, err = file.Read(bs)
	if err != nil {
		return nil, fmt.Errorf("error reading file %s: %w", filename, err)
	}

	return bs, nil
}

func Write(filename string, data []byte) error {
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return fmt.Errorf("error opening file %s: %w", filename, err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)

	// записываем событие в буфер
	if _, err = writer.Write(data); err != nil {
		return fmt.Errorf("error writing file %s: %w", filename, err)
	}

	// записываем буфер в файл
	return writer.Flush()
}
