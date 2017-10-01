package utils


func NewBytes(n int) (b *Bytes) {
	b = &Bytes{buf: make([]byte, 0, n)}
	return
}

type Bytes struct {
	buf []byte
}

/**
下面两个是兼容go的bytes.writer接口用的。
 */
func (b *Bytes) WriteByte(bt byte) {
	b.buf = append(b.buf, bt)
}

func (b *Bytes) Write(bts []byte) {
	b.buf = append(b.buf, bts...)
}

func (b *Bytes) Add(bt byte) {
	b.buf = append(b.buf, bt)
}

func (b *Bytes) AddBytes(bts []byte) {
	b.buf = append(b.buf, bts...)
}

func (b *Bytes) WriteString(str string) {
	b.buf = append(b.buf, str...)
}

func (b *Bytes) Change(idx int, bt byte) bool {
	if idx >= len(b.buf) {
		return false
	}
	b.buf[idx] = bt
	return true
}

func (b *Bytes) ChangeBlock(start int, bts []byte) bool {
	if start >= len(b.buf)-len(bts) {
		return false
	}
	copy(b.buf[start:], bts)
	return true
}

func (b *Bytes) Len() int {
	return len(b.buf)
}

func (b *Bytes) Bytes() []byte {
	return b.buf
}

func (b *Bytes) String() (str string) {
	str = string(b.buf)
	return
}
