package jsons

import (
	"fmt"

	"github.com/guanyaowen/puer/util/bytesconv"
	jsoniter "github.com/json-iterator/go"
)

// 滴滴开源的一个100%兼容原生json包的序列化库
// https://github.com/json-iterator/go
//
// 					ns/op		allocation bytes	allocation times
// std decode		35510 ns/op	1960 B/op			99 allocs/op
// jsoniter decode	5623 ns/op	160 B/op			3 allocs/op
// std encode		2213 ns/op	712 B/op			5 allocs/op
// jsoniter encode	837 ns/op	384 B/op			4 allocs/op

var json = jsoniter.ConfigCompatibleWithStandardLibrary

func Marshal(v interface{}) ([]byte, error) {
	return json.Marshal(&v)
}

func Unmarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

func JsonConv(v, target interface{}) error {
	bytes, err := json.Marshal(&v)
	if err != nil {
		return fmt.Errorf("entitys：json marshal err: %w", err)
	}
	err = json.Unmarshal(bytes, &target)
	if err != nil {
		return fmt.Errorf("entitys：json unmarshal err: %w", err)
	}
	return nil
}

func Valid(data []byte) bool {
	if len(data) == 0 {
		return false
	}
	return json.Valid(data)
}

// ToJson 对象转json格式
func ToJson(data interface{}) string {
	bytes, err := Marshal(data)
	if err != nil {
		return ""
	}
	return bytesconv.BytesToString(bytes)
}
