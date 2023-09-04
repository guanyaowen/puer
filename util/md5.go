package util

import (
	"crypto/md5"
	"fmt"

	"github.com/guanyaowen/puer/util/bytesconv"
)

const hexTable = "0123456789abcdef"

func MD5(str string) string {
	src := md5.Sum(bytesconv.StringToBytes(str))
	var dst = make([]byte, 32)
	j := 0
	for _, v := range src {
		dst[j] = hexTable[v>>4]
		dst[j+1] = hexTable[v&0x0f]
		j += 2
	}
	return bytesconv.BytesToString(dst)
}

// Md5String 获取md5值
// 调用mis的接口鉴权用
func Md5String(data string) string {
	srcCode := md5.Sum(bytesconv.StringToBytes(data))
	return fmt.Sprintf("%x", srcCode)
}
