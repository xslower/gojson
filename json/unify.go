package json

import (
	// `strings`
	"errors"
	"strconv"
)

var (
	format_err   = errors.New(`the json string is not valid!`)
	not_a_object = errors.New(`the json is not a object!`)
	not_a_array  = errors.New(`the json is not a array!`)
	not_obj_obj  = errors.New(`the json format is not {object{object}}!`)
	not_obj_arr  = errors.New(`the json format is not {object[array]}!`)
	not_arr_arr  = errors.New(`the json format is not [array[array]]!`)
	not_arr_obj  = errors.New(`the json format is not [array{object}]!`)
)

func Encode(input interface{}) string {
	var buf = make([]byte, 0, 512)
	switch slc := input.(type) {
	case int:
		return strconv.FormatInt(int64(slc), 10)
	case int64:
		return strconv.FormatInt(slc, 10)
	case string:
		buf = append(buf, '"')
		buf = append(buf, slc...)
		buf = append(buf, '"')
	case map[string]int:
		buf = encode_map_int(slc)
	case map[string]int64:
		buf = encode_map_int64(slc)
	case map[string]string:
		buf = encode_map_str(slc)
	case []int:
		buf = encode_slc_int(slc)
	case []int64:
		buf = encode_slc_int64(slc)
	case []string:
		buf = encode_slc_str(slc)
	case map[string]map[string]int:
		buf = append(buf, '{')
		for key, val := range slc {
			buf = append(buf, '"')
			buf = append(buf, key...)
			buf = append(buf, '"', ':')
			buf = append(buf, encode_map_int(val)...)
			buf = append(buf, ',')
		}
		buf[len(buf)-1] = '}'
	case map[string]map[string]int64:
		buf = append(buf, '{')
		for key, val := range slc {
			buf = append(buf, '"')
			buf = append(buf, key...)
			buf = append(buf, '"', ':')
			buf = append(buf, encode_map_int64(val)...)
			buf = append(buf, ',')
		}
		buf[len(buf)-1] = '}'
	case map[string]map[string]string:
		buf = append(buf, '{')
		for key, val := range slc {
			buf = append(buf, '"')
			buf = append(buf, key...)
			buf = append(buf, '"', ':')
			buf = append(buf, encode_map_str(val)...)
			buf = append(buf, ',')
		}
		buf[len(buf)-1] = '}'
	case map[string][]int:
		buf = append(buf, '{')
		for key, val := range slc {
			buf = append(buf, '"')
			buf = append(buf, key...)
			buf = append(buf, '"', ':')
			buf = append(buf, encode_slc_int(val)...)
			buf = append(buf, ',')
		}
		buf[len(buf)-1] = '}'
	case map[string][]int64:
		buf = append(buf, '{')
		for key, val := range slc {
			buf = append(buf, '"')
			buf = append(buf, key...)
			buf = append(buf, '"', ':')
			buf = append(buf, encode_slc_int64(val)...)
			buf = append(buf, ',')
		}
		buf[len(buf)-1] = '}'
	case map[string][]string:
		buf = append(buf, '{')
		for key, val := range slc {
			buf = append(buf, '"')
			buf = append(buf, key...)
			buf = append(buf, '"', ':')
			buf = append(buf, encode_slc_str(val)...)
			buf = append(buf, ',')
		}
		buf[len(buf)-1] = '}'
	case []map[string]int:
		buf = append(buf, '[')
		for _, val := range slc {
			buf = append(buf, encode_map_int(val)...)
			buf = append(buf, ',')
		}
		buf[len(buf)-1] = ']'
	case []map[string]int64:
		buf = append(buf, '[')
		for _, val := range slc {
			buf = append(buf, encode_map_int64(val)...)
			buf = append(buf, ',')
		}
		buf[len(buf)-1] = ']'
	case []map[string]string:
		buf = append(buf, '[')
		for _, val := range slc {
			buf = append(buf, encode_map_str(val)...)
			buf = append(buf, ',')
		}
		buf[len(buf)-1] = ']'
	case [][]int:
		buf = append(buf, '[')
		for _, val := range slc {
			buf = append(buf, encode_slc_int(val)...)
			buf = append(buf, ',')
		}
		buf[len(buf)-1] = ']'
	case [][]int64:
		buf = append(buf, '[')
		for _, val := range slc {
			buf = append(buf, encode_slc_int64(val)...)
			buf = append(buf, ',')
		}
		buf[len(buf)-1] = ']'
	case [][]string:
		buf = append(buf, '[')
		for _, val := range slc {
			buf = append(buf, encode_slc_str(val)...)
			buf = append(buf, ',')
		}
		buf[len(buf)-1] = ']'
	//下面三种情况是为了增加适应面而放弃了性能
	case Stringor:
		buf = slc.Stringify()
	default:
		panic(`not support type!`)
	}
	return string(buf)
}

/**
 * 所有自定义类型要实现encode就要实现下面的接口。
 * TODO 写个prebuilder为所有需要encode的自定义类型实现下面的接口
 */
type Stringor interface {
	Stringify() (ret []byte)
}
