package nolock

import (
	"sync/atomic"
	"runtime"
)

func NewPtrLocker(ptr *int32) (lck *PtrLocker) {
	lck = &PtrLocker{ptr}
	return
}

type PtrLocker struct {
	ptr *int32
}

func (this *PtrLocker) WLock() {
	_WLock(this.ptr)
}
func (this *PtrLocker) WUnlock() {
	_WUnLock(this.ptr)
}
func (this *PtrLocker) Lock() {
	_RLock(this.ptr)
}
func (this *PtrLocker) Unlock() {
	_RUnLock(this.ptr)
}

/**
这个锁必须使用指针，不然atm发生了复制则锁会出错。
 */
type Locker struct {
	atm int32
}

func NewLocker() (lck *Locker) {
	lck = &Locker{}
	return
}
func (l *Locker) WLock() {
	_WLock(&l.atm)
}
func (l *Locker) WUnlock() {
	_WUnLock(&l.atm)
}
func (l *Locker) Lock() {
	_RLock(&l.atm)
}
func (l *Locker) Unlock() {
	_RUnLock(&l.atm)
}

func _WLock(addr *int32) {
	var lck, nlck int32
	for {
		lck = atomic.LoadInt32(addr)
		if lck != 0 { //=0 才意味没有读，也没有写
			//fmt.Println(`[put] write lock occupied`)
			runtime.Gosched()
		}
		nlck = lck - 1
		if atomic.CompareAndSwapInt32(addr, lck, nlck) {
			break
		} else {
			//fmt.Println(`[put] write lock changed`)
			runtime.Gosched()
		}
	}
}
func _WUnLock(addr *int32) {
	atomic.StoreInt32(addr, 0)
}
func _RLock(addr *int32) {
	var lck int32
	for {
		lck = atomic.LoadInt32(addr)
		if lck < 0 { //正在写入不能读，其它情况都能读
			//fmt.Println(`[get] read lock occupied`)
			runtime.Gosched()
		}
		atomic.AddInt32(addr, 1)
		break
	}
}
func _RUnLock(addr *int32) {
	atomic.AddInt32(addr, -1)
}
