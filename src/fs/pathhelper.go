package fs

import (
	"fmt"
	"io"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"errors"
)

// GetParentAbsDir 获取所在目录的绝对路径
func GetParentAbsDir(path string) (string, error) {
	return filepath.Abs(filepath.Dir(path))
}

// IsDir 是否是目录
func IsDir(file string) bool {
	fi, err := os.Stat(file)
	if err != nil {
		return false
	}
	return fi.IsDir()
}

// IsFile 是否是文件
func IsFile(file string) bool {
	fi, err := os.Stat(file)
	if err != nil {
		return false
	}
	return !fi.IsDir()
}

// GetFileName /usr/local/test.txt -> test
func GetFileName(inputPath string) string {
	filenameWithSuffix := path.Base(inputPath)
	fileSuffix := path.Ext(inputPath)
	return strings.TrimSuffix(filenameWithSuffix, fileSuffix)
}

// FileNameAppend /usr/local/test.txt -> /usr/local/test$(appText).txt
func FileNameAppend(inputPath, appText, defaultExt string) string {
	if inputPath == "" {
		return appText
	}
	newFilename := fmt.Sprintf("%v%s", GetFileName(inputPath), appText)

	ext := filepath.Ext(inputPath)
	if ext == "" {
		ext = defaultExt
	}
	outputPath := filepath.Join(filepath.Dir(inputPath), newFilename+ext)
	return outputPath
}

// LocalCopy 本地拷贝
func LocalCopy(src, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}
	if src == dst {
		return sourceFileStat.Size(), nil
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	srcHandle, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer srcHandle.Close()

	dstHandle, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer dstHandle.Close()
	nBytes, err := io.Copy(dstHandle, srcHandle)
	return nBytes, err
}

// GetNameWithNewExt 替换新后缀
func GetNameWithNewExt(u string, ext string) string {
	oldExt := path.Ext(u)
	if oldExt == ext {
		return u
	} else if len(oldExt) == 0 {
		return u + ext
	}
	dir := path.Dir(u)
	return path.Join(dir, GetFileName(u)+ext)
}

// AddNonceParam add nonce str to url path
func AddNonceParam(u string, t time.Time) (string, error) {
	if u == "" {
		return "", errors.New("url should not be empty!")
	}
	url, err := url.Parse(u)
	if err != nil {
		return "", err
	}
	query := url.Query()
	query.Set("noncestr", fmt.Sprintf("%v", t.Unix()))
	url.RawQuery = query.Encode()
	return url.String(), nil
}

// NormalizeURL 格式化 URL 加上格式的 scheme
func NormalizeURL(url string) string {
	if url == "" {
		return url
	}
	if strings.HasPrefix(url, "//") {
		url = "http:" + url
	} else if !strings.HasPrefix(url, "/") && !strings.HasPrefix(url, "http") {
		url = "http://" + url
	}
	return url
}
