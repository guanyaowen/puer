package util

import (
	"fmt"
	"runtime"
)

func NoPanic(f func() error) func() error {
	return func() (err error) {
		defer func() {
			if re := recover(); re != nil {
				buf := make([]byte, 4096)
				runtime.Stack(buf, false)
				err = fmt.Errorf("panic: %v, %s", re, FilterNull(buf))
			}
		}()

		return f()
	}
}
