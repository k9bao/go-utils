// Package algorithm 提供常用算法封装
package algorithm

import (
	"math"
	"strconv"
)

// MaxInt64 取最大值
func MaxInt64(a, b int64) int64 {
	if a <= b {
		return b
	}
	return a
}

// MinInt64 取最小值
func MinInt64(a, b int64) int64 {
	if a >= b {
		return b
	}
	return a
}

// MaxInt 取最大值
func MaxInt(a, b int) int {
	if a <= b {
		return b
	}
	return a
}

// MinInt 取最小值
func MinInt(a, b int) int {
	if a >= b {
		return b
	}
	return a
}

// ParseFloat 字符串转float64
func ParseFloat(s string, defaultVal float64) float64 {
	if out, err := strconv.ParseFloat(s, 64); err == nil {
		return out
	}
	return defaultVal
}

// ParseFloatDecimals 保留指定小数位
func ParseFloatDecimals(s string, defaultVal float64, decimals int) float64 {
	if out, err := strconv.ParseFloat(s, 64); err == nil {
		return math.Trunc(out*math.Pow10(decimals)+0.5) / math.Pow10(decimals)
	}

	return defaultVal
}

// ParseInt 字符串到int64
func ParseInt(s string, defaultVal int64) int64 {
	if out, err := strconv.ParseInt(s, 10, 64); err == nil {
		return out
	}
	return defaultVal
}
