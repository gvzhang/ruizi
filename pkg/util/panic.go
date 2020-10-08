package util

import (
	"fmt"
	"ruizi/pkg/logger"
	"runtime"

	"go.uber.org/zap"
)

// PanicToError Panic转换为error
func PanicToError(f func()) (err error) {
	defer func() {
		if r := recover(); r != nil {
			logger.Logger.DPanic("panic", zap.Any("panic", r), zap.String("stack", GetStackInfo()))
		}
	}()
	f()
	return
}

// RecoverPanic 恢复panic
func RecoverPanic() {
	err := recover()
	if err != nil {
		logger.Logger.DPanic("panic", zap.Any("panic", err), zap.String("stack", GetStackInfo()))
	}
}

// GetStackInfo 获取Panic堆栈信息
func GetStackInfo() string {
	buf := make([]byte, 4096)
	n := runtime.Stack(buf, false)
	return fmt.Sprintf("%s", buf[:n])
}
