package cmap

import (
	"bytes"
	"fmt"
	"sync"
	"sync/atomic"
)

// Bucket 代表并发安全的散列桶的接口。
type Bucket interface {
	// Put 会放入一个键-元素对。
	// 第一个返回值表示是否新增了键-元素对。
	// 若在调用此方法前已经锁定lock，则不要把lock传入！否则必须传入对应的lock！
	Put(p Pair, lock sync.Locker) (bool, error)
	// Get 会获取指定键的键-元素对。
	Get(key string) Pair
	// GetFirstPair 会返回第一个键-元素对。
	GetFirstPair() Pair
	// Delete 会删除指定的键-元素对。
	// 若在调用此方法前已经锁定lock，则不要把lock传入！否则必须传入对应的lock！
	Delete(key string, lock sync.Locker) bool
	// Clear 会清空当前散列桶。
	// 若在调用此方法前已经锁定lock，则不要把lock传入！否则必须传入对应的lock！
	Clear(lock sync.Locker)
	// Size 会返回当前散列桶的尺寸。
	Size() uint64
	// String 会返回当前散列桶的字符串表示形式。
	String() string
}

type bucket struct {
	firstValue atomic.Value
	size       uint64
}

// 占位符。
// 由于原子值不能存储nil，所以当散列桶空时用此符占位。
var placeholder Pair = &pair{}

// newBucket 会创建一个Bucket类型的实例。简单工厂模式
func newBucket() Bucket {
	b := &bucket{}
	b.firstValue.Store(placeholder)

	return b
}

func (b *bucket) GetFirstPair() Pair {
	if v := b.firstValue.Load(); v == nil {
		return nil
	} else if p, ok := v.(Pair); !ok || p == placeholder {
		return nil
	} else {
		return p
	}
}

func (b *bucket) Put(p Pair, lock sync.Locker) (bool, error) {
	if p == nil {
		return false, newIllegalParameterError("pair is nil")
	}

	if lock != nil {
		lock.Lock()
		defer lock.Unlock()
	}

	firstPair := b.GetFirstPair()
	if firstPair == nil {
		b.firstValue.Store(p)
		atomic.AddUint64(&b.size, 1)
		return true, nil
	}

	var target Pair
	key := p.Key()
	for v := firstPair; v != nil; v = v.Next() {
		if v.Key() == key {
			target = v
			break
		}
	}

	if target != nil {
		target.SetElement(p.Element())
		return false, nil
	}

	p.SetNext(firstPair)
	b.firstValue.Store(p)
	atomic.AddUint64(&b.size, 1)

	return true, nil
}

// 将被删除的元素前面的元素都存储起来构成一个新链表，然后将删除元素的下一个元素接上去
func (b *bucket) Delete(key string, lock sync.Locker) bool {
	if lock != nil {
		lock.Lock()
		defer lock.Unlock()
	}

	firstPair := b.GetFirstPair()
	if firstPair == nil {
		return false
	}

	var prevPairs []Pair
	var target Pair
	var breakpoint Pair

	for v := firstPair; v != nil; v = v.Next() {
		if v.Key() == key {
			target = v
			breakpoint = v.Next()
			break
		}
		prevPairs = append(prevPairs, v)
	}

	if target == nil {
		return false
	}

	newFirstPair := breakpoint
	for i := len(prevPairs) - 1; i >= 0; i-- {
		pairCopy := prevPairs[i].Copy()
		pairCopy.SetNext(newFirstPair)
		newFirstPair = pairCopy
	}

	if newFirstPair != nil {
		b.firstValue.Store(newFirstPair)
	} else {
		b.firstValue.Store(placeholder)
	}

	atomic.AddUint64(&b.size, ^uint64(0))
	return true
}

func (b *bucket) Get(key string) Pair {
	firstPair := b.GetFirstPair()
	for v := firstPair; v != nil; v = v.Next() {
		if v.Key() == key {
			return v
		}
	}

	return nil
}

func (b *bucket) Clear(lock sync.Locker) {
	if lock != nil {
		lock.Lock()
		defer lock.Unlock()
	}

	b.firstValue.Store(placeholder)
	atomic.StoreUint64(&b.size, 0)
}

func (b *bucket) Size() uint64 {
	return atomic.LoadUint64(&b.size)
}

func (b *bucket) String() string {
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("Len: %d, [ ", b.size))

	for v := b.GetFirstPair(); v != nil; v = v.Next() {
		buf.WriteString(v.String())
		buf.WriteString(" ")
	}

	buf.WriteString(" ]")
	return buf.String()
}
