package utils

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"
)

var ()

func init() {
	// date := Date(`Y-m-d`, time.Now())
}

func Throw(err error, msg ...interface{}) {
	if err == nil {
		return
	}
	// for i, data := range msg {
	// 	m := ``
	// 	if len(data) > 100 {
	// 		m = string(data[len(data)-100:])
	// 	} else {
	// 		m = string(data)
	// 	}
	// 	msg[i] = m
	// }

	panic(err.Error() + fmt.Sprint(msg...))

}

//以标准的日期格式输入，获得时间字符串
func Date(format string, t time.Time) string {
	format = strings.Replace(format, `Y`, `2006`, -1)
	format = strings.Replace(format, `m`, `01`, -1)
	format = strings.Replace(format, `d`, `02`, -1)
	format = strings.Replace(format, `H`, `15`, -1)
	format = strings.Replace(format, `i`, `04`, -1)
	format = strings.Replace(format, `s`, `05`, -1)
	return t.Format(format)
}

/**
 * BKDR Hash Function
 * 把字符串hash到64位整数
 */
func BKDRHash(str string) uint64 {
	var seed uint64 = 131 // 31 131 1313 13131 131313 etc..
	var hash uint64 = 0

	for i := 0; i < len(str); i++ {
		var ui = uint64(str[i])
		hash = hash*seed + ui
	}
	return hash
}

// 复制文件
func CopyFile(src, dst string) (w int64, err error) {
	srcFile, err := os.Open(src)
	if err != nil {
		return
	}
	defer srcFile.Close()
	dstFile, err := os.Create(dst)
	if err != nil {
		return
	}
	defer dstFile.Close()
	return io.Copy(dstFile, srcFile)
}

//判断文件目录是否存在
func PathExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}

const (
	TYPE_NOT_SUPPORT uint8 = iota
	TYPE_STRING
	TYPE_INT
	TYPE_UINT
	TYPE_FLOAT
	TYPE_BOOL
	TYPE_NULL
)

//not support int8/uint8 please use int16/uint16
func InterfaceToString(ifc interface{}) (value string, stype uint8) {
	if ifc == nil {
		return ``, TYPE_NULL
	}
	switch val := ifc.(type) {
	case int16:
		value, stype = strconv.FormatInt(int64(val), 10), TYPE_INT
	case int32:
		value, stype = strconv.FormatInt(int64(val), 10), TYPE_INT
	case int64:
		value, stype = strconv.FormatInt(int64(val), 10), TYPE_INT
	case int:
		value, stype = strconv.FormatInt(int64(val), 10), TYPE_INT
	case uint16:
		value, stype = strconv.FormatUint(uint64(val), 10), TYPE_UINT
	case uint32:
		value, stype = strconv.FormatUint(uint64(val), 10), TYPE_UINT
	case uint64:
		value, stype = strconv.FormatUint(uint64(val), 10), TYPE_UINT
	case uint:
		value, stype = strconv.FormatUint(uint64(val), 10), TYPE_UINT
	case float32:
		value, stype = strconv.FormatFloat(float64(val), 'g', -1, 32), TYPE_FLOAT
	case float64:
		value, stype = strconv.FormatFloat(val, 'g', -1, 64), TYPE_FLOAT
	case bool:
		ret := `0`
		if val {
			ret = `1`
		}
		value, stype = ret, TYPE_BOOL
	case string:
		value, stype = val, TYPE_STRING
	case []byte:
		value, stype = string(val), TYPE_STRING
	case time.Time:
		value, stype = ``, TYPE_STRING
	default:
		value, stype = ``, TYPE_NOT_SUPPORT
	}
	return
}

//not support int8/uint8 please use int16/uint16
func InterfaceToInt(ifc interface{}) (value int, stype uint8) {
	if ifc == nil {
		return 0, TYPE_NULL
	}
	switch val := ifc.(type) {
	case int16:
		value, stype = int(val), TYPE_INT
	case int32:
		value, stype = int(val), TYPE_INT
	case int64:
		value, stype = int(val), TYPE_INT
	case int:
		value, stype = val, TYPE_INT
	case uint16:
		value, stype = int(val), TYPE_UINT
	case uint32:
		value, stype = int(val), TYPE_UINT
	case uint64:
		value, stype = int(val), TYPE_UINT
	case uint:
		value, stype = int(val), TYPE_UINT
	case float32:
		value, stype = int(val), TYPE_FLOAT
	case float64:
		value, stype = int(val), TYPE_FLOAT
	case bool:
		ret := 0
		if val {
			ret = 1
		}
		value, stype = ret, TYPE_BOOL
	case string:
		ret, _ := strconv.Atoi(val)
		value, stype = ret, TYPE_STRING
	case []byte:
		ret, _ := strconv.Atoi(string(val))
		value, stype = ret, TYPE_STRING
	case time.Time:
		value, stype = 0, TYPE_STRING
	default:
		value, stype = 0, TYPE_NOT_SUPPORT
	}
	return
}
