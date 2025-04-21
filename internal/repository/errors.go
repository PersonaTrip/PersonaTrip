package repository

import "errors"

// 常见错误定义
var (
	ErrNotFound = errors.New("resource not found")
)
