// Package contextutil 提供context常用函数封装
package contextutil

import "context"

// CheckCtx 检测 ctx 是否结束
func CheckCtx(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		return nil
	}
}
