package utils

import (
	"fmt"
	"os"
)

// 判断文件是否存在
func FileExists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			fmt.Printf("[FileExists] error(%s)", err)
			return false
		}
	}
	return true
}
