package utils

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

var ()

func init() {
	// date := Date(`Y-m-d`, time.Now())
}

func Throw(err error, msg ...string) {
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
