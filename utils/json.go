package utils

import (
	"encoding/json"
	"log"
)

// ObjToString 结构体转换string
func ObjToString(params interface{}) string {
	if data, err := json.Marshal(params); err != nil {
		return ""
	} else {
		return string(data)
	}
}

// ObjToByteList 结构体转换byte数组
func ObjToByteList(params interface{}) []byte {
	if data, err := json.Marshal(params); err != nil {
		return nil
	} else {
		return data
	}
}

// ByteListToObj byte数组转换结构体
func ByteListToObj(params []byte, obj interface{}) {
	if err := json.Unmarshal(params, obj); err != nil {
		log.Println(err)
	}
}
