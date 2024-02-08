package memcache

import (
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
)

const (
	B = 1 << (iota * 10)
	KB
	MB
	GB
	TB
)

func ParseSize(size string) (int64, string) {

	re, _ := regexp.Compile("[0-9]+")

	//单位
	unit := string(re.ReplaceAll([]byte(size), []byte("")))

	//数字部分
	num, _ := strconv.ParseInt(strings.Replace(size, unit, "", 1), 10, 64)

	//转换大写单位
	var byteSize int64 = 0
	unit = strings.ToUpper(unit)
	switch unit {
	case "B":
		byteSize = num
	case "KB":
		byteSize = num * KB
	case "MB":
		byteSize = num * MB
	case "GB":
		byteSize = num * GB
	case "TB":
		byteSize = num * TB
	default:
		byteSize = 0
	}

	//默认值
	if num == 0 {
		log.Println("ParseSize 仅支持 B KB MB GB TB")
		unit = "MB"
		num = 100
		byteSize = num * MB
	}
	return byteSize, fmt.Sprintf("%d%s", num, unit)
}

func GetValueSize(val interface{}) int64 {
	bytes, _ := json.Marshal(val)
	size := int64(len(bytes))
	//size := unsafe.Sizeof(val)
	fmt.Println(size)
	return 0
}
