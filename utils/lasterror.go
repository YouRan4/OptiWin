package utils

import "sync"

var (
	lastErrMu sync.Mutex
	lastErr   string
)

func SetLastError(msg string) {
	lastErrMu.Lock()
	lastErr = msg
	lastErrMu.Unlock()
}

// TakeLastError 取出并清空：一次失败只对应一次提示，不会粘到下一次操作上
func TakeLastError() string {
	lastErrMu.Lock()
	defer lastErrMu.Unlock()
	msg := lastErr
	lastErr = ""
	return msg
}
