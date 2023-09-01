package utils

import (
	"io"
	"os"
)

func IsFileExist(filePath string) bool {
	_, err := os.Stat(filePath)
	return err == nil
}

func FileReadString(filePath string) (string, error) {
	f, err := os.Open(filePath)
	defer func(file *os.File) {
		// ignore file close error
		file.Close()
	}(f)
	if err != nil {
		return "", err
	}

	var result []byte
	readBuff := make([]byte, 1024*4)
	for {
		n, err := f.Read(readBuff)
		if err != nil {
			if err == io.EOF {
				if n != 0 {
					result = append(result, readBuff[:n]...)
				}
				break
			}
			return "", err
		}
		result = append(result, readBuff[:n]...)
	}
	return string(result), nil
}
func FileSaveString(filePath string, data string) error {
	f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}

	data_length := len(data)
	send_count := 0
	for {
		n, err := f.WriteString(data[send_count:])
		if err != nil {
			return err
		}
		if n+send_count >= data_length {
			send_count += n
			break
		}
	}
	return nil
}
