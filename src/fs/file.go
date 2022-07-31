// Package fs 提供文件系统的常用函数封装，比如是否是文件...
package fs

import (
	"os"
	"fmt"

	"errors"
)

// GetFileSizeNoErr 获取文件大小
func GetFileSizeNoErr(filename string) int64 {
	fi, err := os.Stat(filename)
	if err == nil {
		return fi.Size()
	}
	return 0
}

// ReadFileWithSize 从文件里指定位置读取指定大小
func ReadFileWithSize(fileURL string, pos int64, size int) ([]byte, error) {
	file, err := os.Open(fileURL)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	newPos, err := file.Seek(pos, 0)
	if err != nil {
		return nil, err
	}

	if newPos != pos {
		return nil, errors.New(fmt.Sprintf("failed to %s seek to %v, actual is %v", fileURL, pos, newPos))
	}

	buf := make([]byte, size)
	n, err := file.Read(buf)
	if err != nil {
		return nil, err
	}

	if n < size {
		return nil, errors.New(fmt.Sprintf("failed to read %s, want %v but got %v", fileURL, size, n))
	}

	return buf, nil
}
