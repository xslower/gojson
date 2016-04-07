package json

import (
	"strings"
)

type Decoder struct {
}

func (d *Decoder) preparse(jn string, lch, rch byte) (ret string, err error) {
	ret = strings.TrimSpace(jn)
	if ret[0] != lch || ret[len(ret)-1] != rch {
		err = format_err
		return
	}
	ret = ret[1 : len(ret)-1]
	return
}

func (d *Decoder) getKey(str string) (key string, err error) {
	var left = strPos(str, '"', 0)
	if left < 0 {
		err = not_a_object
		return
	}
	left++
	var right = strPos(str, '"', left)
	if right < 0 {
		err = not_a_object
		return
	}
	key = str[left:right]
	return
}

func (d *Decoder) to_slc(jn string) (ret []string, err error) {
	var cnt = strings.Count(jn, `,`)
	ret = make([]string, cnt+1)
	var last = 0
	var n = 0
	for i := 0; i < len(jn); i++ {
		if jn[i] == ',' {
			ret[n] = trim(jn[last:i], '"')
			n++
			last = i + 1
		}
	}
	ret[n] = trim(jn[last:], '"')
	return
	// ret = strings.Split(jn, `,`)
	// for i, _ := range ret { //没有检查引号是否对称
	// 	ret[i] = trim(ret[i], '"')
	// }
}

func (d *Decoder) ToSlc(jn string) (ret []string, err error) {
	jn, err = d.preparse(jn, '[', ']')
	if err != nil {
		return
	}
	ret, err = d.to_slc(jn)
	return
}

func (d *Decoder) to_map(jn string) (ret map[string]string, err error) {
	ret = map[string]string{}
	var last = 0
	var k string
	for i := 0; i < len(jn); i++ {
		if jn[i] == ':' {
			k, err = d.getKey(jn[last:i])
			if err != nil {
				return
			}
			last = i + 1
		} else if jn[i] == ',' {
			ret[k] = trim(jn[last:i], '"')
			last = i + 1
		}
	}
	ret[k] = trim(jn[last:], '"')
	return
	// var parts = strings.Split(jn, `,`)
	// ret = make(map[string]string)
	// var k, v string
	// for _, str := range parts {
	// 	k, v, err = getKeyValue(str)
	// 	if err != nil {
	// 		return
	// 	}
	// 	ret[k] = v
	// }
}

func (d *Decoder) ToMap(jn string) (ret map[string]string, err error) {
	jn, err = d.preparse(jn, '{', '}')
	if err != nil {
		return
	}
	ret, err = d.to_map(jn)
	return
}

func (d *Decoder) ToMapMap(jn string) (ret map[string]map[string]string, err error) {
	jn, err = d.preparse(jn, '{', '}')
	if err != nil {
		return
	}
	ret = map[string]map[string]string{}
	var last = 0
	var pair = 0
	var k string
	for i := 0; i < len(jn); i++ {
		if jn[i] == '{' {
			k, err = d.getKey(jn[last:i])
			if err != nil {
				return
			}
			last = i + 1
			pair++
		} else if jn[i] == '}' {
			ret[k], err = d.to_map(jn[last:i])
			if err != nil {
				return
			}
			last = i + 1
			pair--
		}
	}
	if pair != 0 { //{}数量是否匹配
		err = format_err
		return
	}
	return
	// var parts = strings.Split(jn, `}`)
	// parts = parts[:len(parts)-1] //最后一个为空所以跳过
	// for _, pt := range parts {
	// 	pos := strPos(pt, '{', 0)
	// 	if pos < 0 {
	// 		return nil, format_err
	// 	}
	// 	keys := strings.Split(pt[:pos], `"`)
	// 	if len(keys) != 3 {
	// 		return nil, format_err
	// 	}
	// 	ret[keys[1]], err = d.to_map(pt[pos+1:])
	// 	if err != nil {
	// 		return
	// 	}
	// }
}

func (d *Decoder) ToMapSlc(jn string) (ret map[string][]string, err error) {
	jn, err = d.preparse(jn, '{', '}')
	if err != nil {
		return
	}
	ret = map[string][]string{}
	var last = 0
	var pair = 0
	var k string
	for i := 0; i < len(jn); i++ {
		if jn[i] == '[' {
			k, err = d.getKey(jn[last:i])
			if err != nil {
				return
			}
			last = i + 1
			pair++
		} else if jn[i] == ']' {
			ret[k], err = d.to_slc(jn[last:i])
			if err != nil {
				return
			}
			last = i + 1
			pair--
		}
	}
	if pair != 0 {
		err = format_err
		return
	}
	return
	// var parts = strings.Split(jn, `]`)
	// parts = parts[:len(parts)-1]
	// for _, pt := range parts {
	// 	pos := strings.IndexByte(pt, '[')
	// 	if pos < 0 {
	// 		return nil, format_err
	// 	}
	// 	keys := strings.Split(pt[:pos], `"`)
	// 	if len(keys) != 3 {
	// 		return nil, format_err
	// 	}
	// 	ret[keys[1]], err = d.to_slc(pt[pos+1:])
	// 	if err != nil {
	// 		return
	// 	}
	// }
}

func (d *Decoder) ToSlcMap(jn string) (ret []map[string]string, err error) {
	jn, err = d.preparse(jn, '[', ']')
	if err != nil {
		return
	}
	var cnt = strings.Count(jn, `}`)
	ret = make([]map[string]string, cnt)
	var n = 0
	var left = 0
	for i := 0; i < len(jn); i++ {
		if jn[i] == '{' {
			left = i + 1
			cnt--
		} else if jn[i] == '}' {
			ret[n], err = d.to_map(jn[left:i])
			if err != nil {
				return
			}
			n++
		}
	}
	if cnt != 0 {
		err = format_err
		return
	}
	return
	// var parts = strings.Split(jn, `}`)
	// parts = parts[:len(parts)-1]
	// ret = make([]map[string]string, len(parts))
	// for i, pt := range parts {
	// 	pos := strings.IndexByte(pt, '{')
	// 	if pos < 0 {
	// 		return nil, format_err
	// 	}
	// 	ret[i], err = d.to_map(pt[pos+1:])
	// 	if err != nil {
	// 		return
	// 	}
	// }
}

func (d *Decoder) ToSlcSlc(jn string) (ret [][]string, err error) {
	jn, err = d.preparse(jn, '[', ']')
	if err != nil {
		return
	}
	var cnt = strings.Count(jn, `]`)
	ret = make([][]string, cnt)
	var n = 0
	var left = 0
	for i := 0; i < len(jn); i++ {
		if jn[i] == '[' {
			left = i + 1
			cnt--
		} else if jn[i] == ']' {
			ret[n], err = d.to_slc(jn[left:i])
			if err != nil {
				return
			}
			n++
		}
	}
	if cnt != 0 {
		err = format_err
		return
	}
	return
	// var parts = strings.Split(jn, `]`)
	// parts = parts[:len(parts)-1]
	// ret = make([][]string, len(parts))
	// for i, pt := range parts {
	// 	pos := strings.IndexByte(pt, '[')
	// 	if pos < 0 {
	// 		return nil, format_err
	// 	}
	// 	ret[i], err = d.to_slc(pt[pos+1:])
	// 	if err != nil {
	// 		return
	// 	}
	// }
	// return
}

/**
 * bottom
 */
