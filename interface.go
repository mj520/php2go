package php2go

import (
	"encoding/json"
	"strconv"
)

//GetInterfaceToString interface 转 string
func GetInterfaceToString(value interface{}) string {
	var key string
	if value == nil {
		return key
	}
	switch value.(type) {
	case float64:
		ft := value.(float64)
		key = strconv.FormatFloat(ft, 'f', -1, 64)
	case float32:
		ft := value.(float32)
		key = strconv.FormatFloat(float64(ft), 'f', -1, 64)
	case int:
		it := value.(int)
		key = strconv.Itoa(it)
	case uint:
		it := value.(uint)
		key = strconv.Itoa(int(it))
	case int8:
		it := value.(int8)
		key = strconv.Itoa(int(it))
	case uint8:
		it := value.(uint8)
		key = strconv.Itoa(int(it))
	case int16:
		it := value.(int16)
		key = strconv.Itoa(int(it))
	case uint16:
		it := value.(uint16)
		key = strconv.Itoa(int(it))
	case int32:
		it := value.(int32)
		key = strconv.Itoa(int(it))
	case uint32:
		it := value.(uint32)
		key = strconv.Itoa(int(it))
	case int64:
		it := value.(int64)
		key = strconv.FormatInt(it, 10)
	case uint64:
		it := value.(uint64)
		key = strconv.FormatUint(it, 10)
	case string:
		key = value.(string)
	case []byte:
		key = string(value.([]byte))
	default:
		newValue, _ := json.Marshal(value)
		key = string(newValue)
	}
	return key
}

//GetInterfaceToInt interface 转 int
func GetInterfaceToInt(value interface{}) int {
	var it int
	switch value.(type) {
	case uint:
		it = int(value.(uint))
		break
	case int8:
		it = int(value.(int8))
		break
	case uint8:
		it = int(value.(uint8))
		break
	case int16:
		it = int(value.(int16))
		break
	case uint16:
		it = int(value.(uint16))
		break
	case int32:
		it = int(value.(int32))
		break
	case uint32:
		it = int(value.(uint32))
		break
	case int64:
		it = int(value.(int64))
		break
	case uint64:
		it = int(value.(uint64))
		break
	case float32:
		it = int(value.(float32))
		break
	case float64:
		it = int(value.(float64))
		break
	case string:
		it, _ = strconv.Atoi(value.(string))
		break
	default:
		it = value.(int)
		break
	}
	return it
}

//GetInterfaceToFloat interface 转 float64
func GetInterfaceToFloat(value interface{}) float64 {
	var it float64
	switch value.(type) {
	case uint:
		it = float64(value.(uint))
		break
	case int8:
		it = float64(value.(int8))
		break
	case uint8:
		it = float64(value.(uint8))
		break
	case int16:
		it = float64(value.(int16))
		break
	case uint16:
		it = float64(value.(uint16))
		break
	case int32:
		it = float64(value.(int32))
		break
	case uint32:
		it = float64(value.(uint32))
		break
	case int64:
		it = float64(value.(int64))
		break
	case uint64:
		it = float64(value.(uint64))
		break
	case float32:
		it = float64(value.(float32))
		break
	case string:
		it, _ = strconv.ParseFloat(value.(string), 64)
		break
	default:
		it = value.(float64)
		break
	}
	return it
}
