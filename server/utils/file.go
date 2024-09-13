package utils

import (
	"io"
	"os"
)

// IsFileExist 判断文件是否存在
//
// 参数：
// filePath: 文件路径
//
// # IsFileExist 判断文件是否存在
//
// 参数：
// filePath string - 待判断的文件路径
//
// 返回值：
// bool - 如果文件存在返回true，否则返回false
func IsFileExist(filePath string) bool {
	_, err := os.Stat(filePath)
	return err == nil
}

// FileReadString 从指定文件路径中读取内容并返回字符串和错误
//
// 参数：
// filePath string - 文件路径
//
// 返回值：
// string - 读取到的文件内容
// error - 如果发生错误，则返回非nil的错误信息
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

// FileSaveString 将字符串 data 保存到指定文件路径 filePath 中，返回可能发生的错误
//
// 参数：
// filePath string - 要保存的文件路径
// data string - 要保存的字符串
//
// 返回值：
// error - 保存过程中可能发生的错误，如果保存成功则返回 nil
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
