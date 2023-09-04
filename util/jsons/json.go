package jsons

import (
	"encoding/json"

	"github.com/guanyaowen/puer/util/bytesconv"
)

// ToJson 对象转json格式
func ToJson(data interface{}) string {
	bytes, err := json.Marshal(data)
	if err != nil {
		return ""
	}
	return bytesconv.BytesToString(bytes)
}
