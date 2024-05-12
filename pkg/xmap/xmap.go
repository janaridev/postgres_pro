package xmap

import (
	"sync"
)

type XMap[K comparable, V any] struct {
	mp sync.Map
}

func (x *XMap[K, V]) Load(key K) (V, bool) {
	loaded, ok := x.mp.Load(key)
	if !ok {
		var zeroV V
		return zeroV, false
	}

	return loaded.(V), true
}

func (x *XMap[K, V]) Delete(key K) {
	x.mp.Delete(key)
}

func (x *XMap[K, V]) Store(key K, value V) {
	x.mp.Store(key, value)
}
