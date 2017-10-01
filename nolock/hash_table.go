package nolock

import (
	"github.com/resure-tech/lib/goutils/utils"
	"time"
)

func NewBaseHashTable(pow uint8) (ht *BaseHashTable) {
	num := uint64(1 << pow)
	ht = &BaseHashTable{
		mod:   num - 1,
		slots: make([]interface{}, num),
		lock:  make([]int32, num),
	}
	return
}

/**
基础的无锁ht
 */
type BaseHashTable struct {
	mod   uint64 //mod=2^N-1
	slots []interface{}
	lock  []int32
}

//Del = Put(key, nil)
func (this *BaseHashTable) PutByKey(key uint64, val interface{}) {
	//取余，
	key = key & this.mod
	addr := &this.lock[key]
	_WLock(addr)
	this.slots[key] = val
	_WUnLock(addr)
	return
}

func (this *BaseHashTable) GetByKey(key uint64) (val interface{}) {
	key = key & this.mod
	addr := &this.lock[key]
	_RLock(addr)
	val = this.slots[key]
	_RUnLock(addr)
	return
}

func (this *BaseHashTable) Len() (uint64) {
	return this.mod + 1
}

func (this *BaseHashTable) Put(key string, val interface{}) {
	idx := utils.BKDRHash(key)
	this.PutByKey(idx, val)
}

func (this *BaseHashTable) Get(key string) (val interface{}) {
	idx := utils.BKDRHash(key)
	val = this.GetByKey(idx)
	return
}

const (
	//默认过期时间
	DEFAULT_TIMEOUT int64 = 3600 * 24 * 30
	//最大遍历步数
	TRAVERSE_NUM uint64 = 1 << 5
	MAX_IDX      uint64 = 1<<64 - 1
)

type Slot struct {
	key  string
	val  interface{}
	time int64
}

func NewHashTable(pow uint8) (ht *HashTable) {
	num := uint64(1 << pow)
	ht = &HashTable{
		mod:   num - 1,
		slots: make([]*Slot, num),
		lock:  make([]int32, num),
	}
	return
}

/**
BaseHashTable加强
增加功能：
1.冲突时，遍历后N位找空位
2.带过期时间

防止锁外操作的方案：
写操作时，全部新创建一个Slot插入。
这样读操作时，Slot内部的值永远不会被修改。
 */
type HashTable struct {
	mod   uint64 //mod=2^N-1
	slots []*Slot
	lock  []int32
}

//Del = Put(key, nil)
func (this *HashTable) _put(idx uint64, slot *Slot) {
	//取余，
	idx = idx & this.mod
	addr := &this.lock[idx]
	_WLock(addr)
	this.slots[idx] = slot
	_WUnLock(addr)
}

func (this *HashTable) _get(key uint64) (val *Slot) {
	key = key & this.mod
	addr := &this.lock[key]
	_RLock(addr)
	val = this.slots[key]
	_RUnLock(addr)
	return
}

//是否已过期
func (this *HashTable) expired(tm int64) bool {
	if tm < 0 { //已删除
		return true
	} else if tm == 0 {
		return false
	}
	now := time.Now().Unix()
	if tm < now { //已过期
		return true
	}
	return false
}

func (this *HashTable) validForPut(slot *Slot, key string) bool {
	if slot == nil {
		return true
	}
	if key == slot.key {
		return true
	}
	return this.expired(slot.time)
}

func (this *HashTable) validForGet(slot *Slot, key string) bool {
	if slot == nil {
		return false
	}
	if key != slot.key {
		return false
	}
	return !this.expired(slot.time)
}

func (this *HashTable) findIdx(dst_idx uint64, key string, valid func(*Slot, string) bool) (idx uint64, slot *Slot) {
	//dst_idx := utils.BKDRHash(key)
	var i uint64 = 0
	for ; i < TRAVERSE_NUM; i++ {
		idx = dst_idx + i
		slot = this._get(idx)
		if valid(slot, key) {
			return
		}
	}
	return 0, nil
}

//@timeout =-1 已删除; =0 永久有效
func (this *HashTable) put(key string, val interface{}, timeout int64) {
	//fmt.Println(`put`)
	var stamp = timeout
	if timeout > 0 {
		stamp += time.Now().Unix()
	}
	slot := &Slot{key: key, val: val, time: stamp}
	dst_idx := utils.BKDRHash(key)
	idx, org_slot := this.findIdx(dst_idx, key, this.validForPut)
	if org_slot != nil {
		this._put(idx, slot)
	} else {
		//如果遍历一圈没发现空白、或无效数据，则强行插到初始位置
		this._put(dst_idx, slot)
	}
}

/**
在指定位置的后cnt个位置遍历查找
 */
func (this *HashTable) get(key string) (val interface{}) {
	//fmt.Println(`get`)
	idx := utils.BKDRHash(key)
	_, slot := this.findIdx(idx, key, this.validForGet)
	if slot != nil {
		val = slot.val
	}
	return
}

//func (this *HashTable) PutByIdx(idx uint64, val interface{}, timeOut int64) {
//	var stamp = timeOut
//	if timeOut > 0 {
//		stamp += time.Now().Unix()
//	}
//	slot := &Slot{key: ``, val: val, time: stamp}
//	this._put(idx, slot)
//}
//
//func (this *HashTable) GetByIdx(idx uint64) (val interface{}) {
//	slot := this._get(idx)
//	if slot != nil {
//		val = slot.val
//	}
//	return
//}

//@timeout过期时间，单位秒
func (this *HashTable) Put(key string, val interface{}, timeOut int64) {
	//idx := utils.BKDRHash(key)
	this.put(key, val, timeOut)
}

func (this *HashTable) Get(key string) (val interface{}) {
	//idx := utils.BKDRHash(key)
	val = this.get(key)
	return
}

func (this *HashTable) Len() (uint64) {
	return this.mod + 1
}

func (this *HashTable) Del(key string) {
	dst_idx := utils.BKDRHash(key)
	//请求删除，必然key是已经存在的
	//所以要使用get方法找到idx
	idx, slot := this.findIdx(dst_idx, key, this.validForGet)
	if slot == nil {
		return
	}
	this._put(idx, nil)
}

func (this *HashTable) Exist(key string) bool {
	dst_idx := utils.BKDRHash(key)
	//判断是否存在，必须使用valid get
	_, slot := this.findIdx(dst_idx, key, this.validForGet)
	if slot == nil {
		return false
	}
	return true
}

func (this *HashTable) Update(key string, val interface{}) {
	dst_idx := utils.BKDRHash(key)
	//请求更新，必然key是已经存在的
	//所以要使用get方法找到idx
	idx, slot := this.findIdx(dst_idx, key, this.validForGet)
	if slot == nil {
		return
	}
	new_slot := &Slot{key, val, slot.time}
	//虽然slot是指针，但不能直接操作slot，防止其它线程读写slot
	this._put(idx, new_slot)
}

func (this *HashTable) GetLockAndAddr(key string) (lk *PtrLocker, val *interface{}) {
	dst_idx := utils.BKDRHash(key)
	//请求加锁，必然key是已经存在的，否则直接写即可，无需加锁
	//所以要使用get方法找到idx
	idx, slot := this.findIdx(dst_idx, key, this.validForGet)
	if slot == nil {
		return
	}
	idx = idx & this.mod
	addr := &this.lock[idx]
	lk = NewPtrLocker(addr)
	val = &slot.val
	return
}

//func (this *HashTable) Lock(idx uint64) {
//	idx = idx & this.mod
//	addr := &this.lock[idx]
//	_WLock(addr)
//}
//
//func (this *HashTable) Unlock(idx uint64) {
//	idx = idx & this.mod
//	addr := &this.lock[idx]
//	_WUnLock(addr)
//}
