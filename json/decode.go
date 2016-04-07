/**
 * 只支持json的1维，2维对象或数组
 * 只支持string+string格式，不对int型转型，因为比较麻烦
 *
 * 性能大概是标准库的1.4倍到4倍
 * 都是string时是接近3~4倍
 * 需要str->int,比其快40％左右
 *
 * 关于其它格式的支持：
 * 可以写个prebuilder为其自动实现一套转换的接口实现
 */
package json

import (
	"strings"
)

var (
	Decode = &Decoder{}
)

func smartSplit(jn string) (ret []string, err error) {
	var parts = strings.Split(jn, `,`)
	var cnt_obj, cnt_arr = 0, 0
	var last = -1
	var complete = ``
	for i := 0; i < len(parts); i++ {
		// parts[i] = strings.TrimSpace(parts[i])
		cnt_obj += strings.Count(parts[i], `{`)
		cnt_obj -= strings.Count(parts[i], `}`)
		cnt_arr += strings.Count(parts[i], `[`)
		cnt_arr -= strings.Count(parts[i], `]`)
		if cnt_obj < 0 || cnt_arr < 0 {
			err = format_err
			return
		} else if cnt_obj == 0 && cnt_arr == 0 {
			if last == -1 {
				complete = parts[i]
			} else {
				complete = strings.Join(parts[last:i+1], `,`)
				last = -1
			}
			ret = append(ret, complete)
		} else { //现在是cnt_obj > 0 || cnt_arr > 0
			if last == -1 {
				last = i
			}
		}
	}
	return
}

func DecodeToSlc(jn string) (ret []string, err error) {
	return Decode.ToSlc(jn)
}

func DecodeToMap(jn string) (ret map[string]string, err error) {
	return Decode.ToMap(jn)
}

func DecodeToMapMap(jn string) (ret map[string]map[string]string, err error) {
	return Decode.ToMapMap(jn)
}

func DecodeToMapSlc(jn string) (ret map[string][]string, err error) {
	return Decode.ToMapSlc(jn)
}

func DecodeToSlcMap(jn string) (ret []map[string]string, err error) {
	return Decode.ToSlcMap(jn)
}

func DecodeToSlcSlc(jn string) (ret [][]string, err error) {
	return Decode.ToSlcSlc(jn)
}

func trim(str string, ch byte) (ret string) {
	var lt = uint32(0)
	var rt = uint32(len(str) - 1)
	for ; lt <= rt; lt++ {
		if str[lt] != ' ' && str[lt] != ch {
			break
		}
	}
	for ; rt >= lt; rt-- {
		if str[rt] != ' ' && str[rt] != ch {
			break
		}
	}
	ret = str[lt : rt+1]
	return
}

func valid(jn, lch, rch string) bool {
	var cnt = strings.Count(jn, lch)
	if cnt == 0 {
		return false
	}
	cnt -= strings.Count(jn, rch)
	if cnt != 0 {
		return false
	}
	return true
}

func strPos(str string, ch byte, offset int) int {
	if offset < 0 || offset >= len(str) {
		offset = 0
	}
	for i := offset; i < len(str); i++ {
		if str[i] == ch {
			return i
		}
	}
	return -1
}

func strRPos(str string, ch byte, offset int) int {
	if offset < 0 || offset >= len(str) {
		offset = len(str) - 1
	}
	for i := offset; i >= 0; i-- {
		if str[i] == ch {
			return i
		}
	}
	return -1
}

func strPairPos(str string, lch, rch byte, offset int) (int, int) {
	if offset < 0 || offset >= len(str) {
		offset = 0
	}
	var left, right = -1, -1
	var ch = lch
	var i = offset
	for ; i < len(str); i++ {
		if str[i] == ch { //找lch
			if left < 0 {
				left = i
				ch = rch
			} else { //找rch
				right = i
				break
			}
		}
	}
	if i == len(str) {
		return -1, -1
	}
	return left, right
}

/**

bottom

*/
