// Package hash 提供MD5/SHA1等封装函数
package hash

import (
	"crypto/md5"
	"crypto/sha1"
	"fmt"
	"io"
	"os"
)

// FMD5 file md5 sum
func FMD5(localFilePath string) (string, error) {
	file, err := os.Open(localFilePath)
	if err != nil {
		return "", err
	}
	defer file.Close()
	md5h := md5.New()
	if _, err := io.Copy(md5h, file); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", md5h.Sum(nil)), nil
}

// FSHA1 file FSHA1 sum
func FSHA1(localFilePath string) (string, error) {
	file, err := os.Open(localFilePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	sha1h := sha1.New()
	if _, err := io.Copy(sha1h, file); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", sha1h.Sum(nil)), nil
}

// MD5 对字符串进行MD5哈希
func MD5(data string) (string, error) {
	t := md5.New()
	_, err := io.WriteString(t, data)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", t.Sum(nil)), nil
}

// Md5String calculate md5 for string
func Md5String(s string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(s)))
}

// SHA1 对字符串进行SHA1哈希
func SHA1(data string) (string, error) {
	t := sha1.New()
	_, err := io.WriteString(t, data)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", t.Sum(nil)), nil
}
