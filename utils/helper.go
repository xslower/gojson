package utils

import (
	"io"
	"os"
	"strings"
	"time"
	"math/rand"
)

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

func StdDate(t time.Time) (date string) {
	date = t.Format(`2006-01-02`)
	return
}
func StdDateTime(t time.Time) (datetime string) {
	datetime = t.Format(`2006-01-02 15:04:05`)
	return
}
func StdTime(t time.Time) (time string) {
	time = t.Format(`15:04:05`)
	return
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

// convert camel string to underlined, XxYy to xx_yy
func ToUnderline(cameled string) string {
	buf := make([]byte, 0, len(cameled)+5)
	for i := 0; i < len(cameled); i++ {
		c := cameled[i]
		if c >= 'A' && c <= 'Z' {
			c += 32
			if i > 0 {
				buf = append(buf, '_')
			}
		}
		buf = append(buf, c)
	}
	return string(buf)
}

func IsUpperCase(c byte) bool {
	if c >= 'A' && c <= 'Z' {
		return true
	}
	return false
}

// convert underlined string to camel, aa_bb to AaBb
func ToCamel(underlined string) string {
	buf := make([]byte, 0, len(underlined))
	upper := true
	for i := 0; i < len(underlined); i++ {
		c := underlined[i]
		if c == '_' {
			upper = true
			continue
		}
		if upper {
			upper = false
			if c >= 'a' && c <= 'z' {
				c -= 32
			}
		}
		buf = append(buf, c)
	}
	return string(buf)
}

var (
	_rand = rand.New(rand.NewSource(time.Now().Unix()))
	_dict = []byte(`_0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ`)
)

func Rand(max int) int {
	return _rand.Intn(max)
}

func RandBytes(ln int) (bs []byte) {
	bs = make([]byte, ln)
	for i := 0; i < ln; i++ {
		r := Rand(len(_dict))
		bs[i] = _dict[r]
	}
	return
}
