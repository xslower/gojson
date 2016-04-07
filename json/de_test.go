package json

import (
	// "runtime"
	// "runtime/debug"
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"
)

func placeHolder() {
	_ = strings.Index(`abc`, `a`)
	_ = strconv.Itoa(1)
	_ = time.Now()
	_ = os.O_WRONLY
	_ = bytes.MinRead
	// _ = filepath.Ext(``)
	// _, _ = ioutil.TempDir(`/tmp/`, `haha`)
	// var a = []int{1, 2}
	// sort.Ints(a)
	// _ = reflect.TypeOf(a)
	// _, _ = json.Marshal(`a`)
	// _ = errors.New(``)
}

func TestDM(t *testing.T) {
	var jn = `{"a": "b", "b":"c" ,"c":"d", "z":"ccccc"}`
	var m, err = Decode.ToMap(jn)
	if err != nil {
		echo(err)
		t.Fail()
	}
	if len(m) != 4 || m[`c`] != `d` || m[`b`] != `c` {
		echo(m)
		t.Fail()
	}
}

func TestDS(t *testing.T) {
	var jn = `[9,8,3,1,83,1,3]`
	var m, err = Decode.ToSlc(jn)
	throw(err)
	if len(m) != 7 || m[1] != `8` || m[6] != `3` {
		echo(m)
		t.Fail()
	}
}

func TestDMM(t *testing.T) {
	var jn = `{"a" :{"a": 1}, "b":{"b":2}, "c":{"c" :3} }`
	var m, err = Decode.ToMapMap(jn)
	throw(err)
	if len(m) != 3 || m[`c`][`c`] != `3` || m[`b`][`b`] != `2` {
		echo(m)
		t.Fail()
	}
}

func TestDMS(t *testing.T) {
	var jn = `{"a":[1,2,3,4,5,6, 7,8,9,10, 11,12],"b":[ "4" ,5]}`
	var m, err = Decode.ToMapSlc(jn)
	throw(err)
	if len(m) != 2 || len(m[`a`]) != 12 || m[`b`][0] != `4` {
		echo(m)
		t.Fail()
	}
}

func TestDSM(t *testing.T) {
	var jn = `[{"a": "b"}, {"b" :2}, {"c":3} ] `
	var s, err = Decode.ToSlcMap(jn)
	throw(err)
	if len(s) != 3 || s[0][`a`] != `b` || s[2][`c`] != `3` {
		echo(s)
		t.Fail()
	}
}

func TestDSS(t *testing.T) {
	var jn = `[[1,2,3],[2,3,4],[3,4,5]]`
	_ = "breakpoint"
	var s, err = Decode.ToSlcSlc(jn)
	throw(err)
	if len(s) != 3 || len(s[1]) != 3 || s[2][2] != `5` {
		echo(s)
		t.Fail()
	}
}

func TestES(t *testing.T) {
	var input = []int{3, 5, 8, 123, 532, 908}
	var str = Encode(input)
	if len(str) != 19 {
		echo(str)
		t.Fail()
	}
}

func TestBenchMarkDecode(t *testing.T) {
	// return
	var total = 5000
	// total = 1
	// var s = `{"a": "bjjj", "c": "dhhh", "ehi": "foo", "g9": "iiih"}`
	var s = `{"a":{"a": "你好"}, "b":{"b": "hahaha"}, "c": {"c":"iii"}}`
	// var s = `{"a":{"a":1}, "b":{"b":2}, "c":{"c":3}}`
	// var s = `{"99":[1,2,3,4,5,6,7,8,9,10,11,12,1000,888888,9999999,7654321],"77":[4,5]}`
	// var s = `[{"a": "bjjj"}, {"c": "dhhh"}, {"ehi": "foo"}, {"g9": "iiih"}]`
	// var s = `[["1","2","3"], ["4","5","6"], ["7","8","9"]]`
	start := time.Now().UnixNano()
	for i := 0; i < total; i++ {
		// var m = map[int][]int{}
		_ = "breakpoint"
		ms, err := Decode.ToMapMap(s)
		// echo(ms)ss
		_ = ms
		throw(err)
		// for k, v := range ms {
		// 	kk, _ := strconv.Atoi(k)
		// 	for _, i := range v {
		// 		it, _ := strconv.Atoi(i)
		// 		m[kk] = append(m[kk], it)
		// 	}
		// }
		// _ = m[`a`]
		// echo(m, len(m))
	}
	// return
	end := time.Now().UnixNano()
	echo(end - start)
	start = time.Now().UnixNano()
	for i := 0; i < total; i++ {
		// var mm = map[int][]int{}
		// var m = map[string][]int{}
		var m = map[string]map[string]string{}
		// var m = []map[string]string{}
		// var m = [][]string{}
		err := json.Unmarshal([]byte(s), &m)
		throw(err)
		// for k, v := range m {
		// 	kk, _ := strconv.Atoi(k)
		// 	mm[kk] = v
		// }
		// echo(m)
	}
	end = time.Now().UnixNano()
	echo(end - start)
}

func TestBenchMarkEncode(t *testing.T) {
	return
	var total = 20000
	// var m = []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 1000, 888888, 9999999, 7654321}
	var m = []string{"a", "b", "c", "d", `e`, `f`, `g`, `h`, `i`, `k`}
	// var m = map[string]string{"a": "bjjj", "c": "dhhh", "ehi": "foo", "g9": "iiih"}
	// var m = map[string][]int{"a": {1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 1000, 888888}, "b": {9999999, 7654321}, "c": {8, 8, 8}}
	// var m = map[string][]string{"a": {"a", "b", "c", "d", "e"}, "b": {"z", "x", "Y", "u", "v", "w"}}
	// var m = map[string]map[string]string{"a": {"b": "bjjj"}, "c": {"b": "dhhh"}, "ehi": {"c": "foo"}, "g9": {"d": "iiih"}}
	// var m = map[string]map[string]int{"a": {"b": 8888}, "b": {"c": 7777}, "c": {"d": 4444}, "d": {"c": 3333}, "e": {"e": 8765}}

	// var w = os.Stderr
	var str string
	var start = time.Now().UnixNano()
	for i := 0; i < total; i++ {
		// str = encode_slc_str(m)
		str = Encode(m)
		// str = EncodeMap(m)
		// str = EncodeMapMap(m)
		// fmt.Fprint(w, str)
		// str = EncodeSlc(m)
		// echo(str)
		// echo(string(w.Bytes()))
	}
	var end = time.Now().UnixNano()
	echo(end - start)
	_ = str
	// echo(str)
	start = time.Now().UnixNano()
	for i := 0; i < total; i++ {
		// str = EncodeMapMap(m)
		// fmt.Fprint(w, str)
		b, _ := json.Marshal(m)
		_ = b
	}
	// echo(str)
	end = time.Now().UnixNano()
	echo(end - start)
}

func throw(err error) {
	if err != nil {
		panic(err)
	}
}
func echo(i ...interface{}) {
	fmt.Println(i...)
}
