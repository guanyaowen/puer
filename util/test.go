package util

import (
	"fmt"
	"runtime"
	"time"
)

// MeasureFuncExecTime 函数执行时间
func MeasureFuncExecTime() func() {
	pc, _, _, _ := runtime.Caller(1)
	funcName := runtime.FuncForPC(pc).Name()
	start := time.Now().UnixNano()
	return func() {
		fmt.Printf("function %s exec cost %+v\n", funcName, time.Duration(time.Now().UnixNano()-start).String())
	}
}
