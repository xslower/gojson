package nolock

import (
	"runtime"
	"sync/atomic"
)

func NewPool(newObj func() interface{}, ln uint16) (p *Pool) {
	p = &Pool{newObj: newObj,
		slots:        make([]interface{}, ln),
		slotLock:     make([]int32, ln),
		len:          ln}
	var i uint16 = 0
	for ; i < ln; i++ {
		p.slots[i] = newObj()
	}
	return
}

/**
连接池之类的基类
 */
type Pool struct {
	newObj   func() interface{}
	slots    []interface{}
	slotLock []int32
	selfLock int32
	len      uint16
}
//暂时不支持调整池大小，麻烦
func (this *Pool) Extend(ln uint16) {

}
func (this *Pool) Borrow() (idx uint16, i interface{}) {
	var slot_lck, new_lck int32
	var addr *int32
	for {
		var i uint16 = 0
		for ; i < this.len; i++ {
			addr = &this.slotLock[i]
			slot_lck = atomic.LoadInt32(addr)
			if slot_lck > 0 {
				continue
			} else {
				new_lck = slot_lck + 1
				if atomic.CompareAndSwapInt32(addr, slot_lck, new_lck) {
					return i, this.slots[i]
				} else {
					continue
				}
			}
		}
		runtime.Gosched()
	}
}

func (this *Pool) GiveBack(idx uint16) {
	atomic.StoreInt32(&this.slotLock[idx], 0)
}
