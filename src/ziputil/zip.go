// Package ziputil zip安装函数封装
package ziputil

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// Zip 打包成zip文件
func Zip(srcDir string, zipFilename string) error {
	os.RemoveAll(zipFilename) // 预防：旧文件无法覆盖

	zipfile, err := os.Create(zipFilename)
	if err != nil {
		return err
	}
	defer zipfile.Close()

	archive := zip.NewWriter(zipfile) // 打开：zip文件
	defer archive.Close()

	return filepath.Walk(srcDir, func(path string, info os.FileInfo, err error) error { // 遍历路径信息

		if err != nil {
			return err
		}

		if path == srcDir { // 如果是源路径，提前进行下一个遍历
			return nil
		}

		// 忽略 zip 包本身
		if path == zipFilename {
			return nil
		}
		header, err := zip.FileInfoHeader(info) // 获取：文件头信息
		if err != nil {
			return err
		}
		header.Name = strings.TrimPrefix(path, srcDir+string(filepath.Separator))

		if info.IsDir() {
			header.Name += string(filepath.Separator)
		} else {
			header.Method = zip.Deflate // 设置：zip的文件压缩算法
		}

		writer, err := archive.CreateHeader(header) // 创建：压缩包头部信息
		if err != nil {
			return err
		}
		if !info.IsDir() {
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()
			_, err = io.Copy(writer, file)
			if err != nil {
				return err
			}
		}
		return nil
	})
}
