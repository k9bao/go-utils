// Package stringutil 字符串常用函数的封装
package stringutil

import (
	"context"
	"crypto/sha1"
	"fmt"
	"math/rand"
	"strconv"

	
	"github.com/google/uuid"
)

// letterBytes 参考 https://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-go
const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

// RandomLengthString generate a given length string
func RandomLengthString(n int) string {
	if n <= 0 {
		logs.Log.Warn( "len is %v", n)
		return ""
	}
	b := make([]byte, n)
	// A rand.Int63() generates 63 random bits, enough for letterIdxMax letters!
	for i, cache, remain := n-1, rand.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = rand.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}

// GetUUID 获取 UUID，独立提取函数，方便单元测试
func GetUUID() string {
	return uuid.New().String()
}
