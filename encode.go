/**
 * 同样只支持1维跟2维的json关系
 * 返回json字符串时，比encoding/json的性能：
 * []int时性能比其快20％
 * map[string][]int时大概2倍
 * 其它3～5倍
 *
 * 尝试了直接传入io.Writer直接输出，发现比拼装成string再输出要慢N倍
 *
 * 关于支持其它格式，
 * 1.可以针对[]struct、map[string]struct、struct各写一套接口
 * 2.可以写个prebuilder跟protobuf一样给struct自动实现一套接口实现
 */
package json

import (
	// "fmt"
	// "io"
	"strconv"
)

func EncodeMapMap(input interface{}) string {
	var buf = make([]byte, 0, 512)
	buf = append(buf, '{')
	switch mp := input.(type) {
	case map[string]map[string]int:
		for key, val := range mp {
			buf = append(buf, '"')
			buf = append(buf, key...)
			buf = append(buf, '"', ':')
			buf = append(buf, encode_map_int(val)...)
			buf = append(buf, ',')
		}
	case map[string]map[string]int64:
		for key, val := range mp {
			buf = append(buf, '"')
			buf = append(buf, key...)
			buf = append(buf, '"', ':')
			buf = append(buf, encode_map_int64(val)...)
			buf = append(buf, ',')
		}
	case map[string]map[string]string:
		for key, val := range mp {
			buf = append(buf, '"')
			buf = append(buf, key...)
			buf = append(buf, '"', ':')
			buf = append(buf, encode_map_str(val)...)
			buf = append(buf, ',')
		}
	}
	buf[len(buf)-1] = '}'
	return string(buf)
}

func EncodeMapSlc(input interface{}) string {
	var buf = make([]byte, 0, 512)
	buf = append(buf, '{')
	switch mp := input.(type) {
	case map[string][]int:
		for key, val := range mp {
			buf = append(buf, '"')
			buf = append(buf, key...)
			buf = append(buf, '"', ':')
			buf = append(buf, encode_slc_int(val)...)
			buf = append(buf, ',')
		}
	case map[string][]int64:
		for key, val := range mp {
			buf = append(buf, '"')
			buf = append(buf, key...)
			buf = append(buf, '"', ':')
			buf = append(buf, encode_slc_int64(val)...)
			buf = append(buf, ',')
		}
	case map[string][]string:
		for key, val := range mp {
			buf = append(buf, '"')
			buf = append(buf, key...)
			buf = append(buf, '"', ':')
			buf = append(buf, encode_slc_str(val)...)
			buf = append(buf, ',')
		}
	}
	buf[len(buf)-1] = '}'
	return string(buf)
}

func EncodeSlcMap(input interface{}) string {
	var buf = make([]byte, 0, 512)
	buf = append(buf, '[')
	switch slc := input.(type) {
	case []map[string]int:
		for _, val := range slc {
			buf = append(buf, encode_map_int(val)...)
			buf = append(buf, ',')
		}
	case []map[string]int64:
		for _, val := range slc {
			buf = append(buf, encode_map_int64(val)...)
			buf = append(buf, ',')
		}
	case []map[string]string:
		for _, val := range slc {
			buf = append(buf, encode_map_str(val)...)
			buf = append(buf, ',')
		}
	}
	buf[len(buf)-1] = ']'
	return string(buf)
}

func EncodeSlcSlc(input interface{}) string {
	var buf = make([]byte, 0, 512)
	buf = append(buf, '[')
	switch slc := input.(type) {
	case [][]int:
		for _, val := range slc {
			buf = append(buf, encode_slc_int(val)...)
			buf = append(buf, ',')
		}
	case [][]int64:
		for _, val := range slc {
			buf = append(buf, encode_slc_int64(val)...)
			buf = append(buf, ',')
		}
	case [][]string:
		for _, val := range slc {
			buf = append(buf, encode_slc_str(val)...)
			buf = append(buf, ',')
		}
	}
	buf[len(buf)-1] = ']'
	return string(buf)
}

func encode_map_int(mp map[string]int) []byte {
	var buf = make([]byte, 0, len(mp)*5)
	buf = append(buf, '{')
	for key, val := range mp {
		buf = append(buf, '"')
		buf = append(buf, key...)
		buf = append(buf, '"', ':')
		buf = append(buf, strconv.FormatInt(int64(val), 10)...)
		buf = append(buf, ',')
	}
	buf[len(buf)-1] = '}'
	return buf
}

func encode_map_int64(mp map[string]int64) []byte {
	var buf = make([]byte, 0, len(mp)*5)
	buf = append(buf, '{')
	for key, val := range mp {
		buf = append(buf, '"')
		buf = append(buf, key...)
		buf = append(buf, '"', ':')
		buf = append(buf, strconv.FormatInt(val, 10)...)
		buf = append(buf, ',')
	}
	buf[len(buf)-1] = '}'
	return buf
}

func encode_map_str(mp map[string]string) []byte {
	var buf = make([]byte, 0, len(mp)*5)
	buf = append(buf, '{')
	for key, val := range mp {
		buf = append(buf, '"')
		buf = append(buf, key...)
		buf = append(buf, '"', ':', '"')
		buf = append(buf, val...)
		buf = append(buf, '"', ',')
	}
	buf[len(buf)-1] = '}'
	return buf
}

func EncodeMap(input interface{}) string {
	var buf []byte
	switch mp := input.(type) {
	case map[string]int:
		buf = encode_map_int(mp)
	case map[string]int64:
		buf = encode_map_int64(mp)
	case map[string]string:
		buf = encode_map_str(mp)
	}
	return string(buf)
}

func encode_slc_int(slc []int) []byte {
	var buf = make([]byte, 0, len(slc)*5)
	buf = append(buf, '[')
	for _, val := range slc {
		buf = append(buf, strconv.FormatInt(int64(val), 10)...)
		buf = append(buf, ',')
	}
	buf[len(buf)-1] = ']'
	return buf
}

func encode_slc_int64(slc []int64) []byte {
	var buf = make([]byte, 0, len(slc)*5)
	buf = append(buf, '[')
	for _, val := range slc {
		buf = append(buf, strconv.FormatInt(val, 10)...)
		buf = append(buf, ',')
	}
	buf[len(buf)-1] = ']'
	return buf
}

func encode_slc_str(slc []string) []byte {
	var buf = make([]byte, 0, len(slc)*5)
	buf = append(buf, '[')
	for _, val := range slc {
		buf = append(buf, '"')
		buf = append(buf, val...)
		buf = append(buf, '"', ',')
	}
	buf[len(buf)-1] = ']'
	return buf
}

func EncodeSlc(input interface{}) string {
	var buf []byte
	switch slc := input.(type) {
	case []int:
		buf = encode_slc_int(slc)
	case []int64:
		buf = encode_slc_int64(slc)
	case []string:
		buf = encode_slc_str(slc)
	}
	return string(buf)
}
