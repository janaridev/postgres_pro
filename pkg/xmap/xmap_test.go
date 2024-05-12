package xmap

import (
	"testing"
)

func TestXMap_Load(t *testing.T) {
	x := &XMap[int, string]{}
	x.Store(1, "value")

	val, ok := x.Load(1)
	if !ok {
		t.Errorf("Load() returned false for existing key")
	}
	if val != "value" {
		t.Errorf("Load() returned incorrect value")
	}

	_, ok = x.Load(2)
	if ok {
		t.Errorf("Load() returned true for non-existing key")
	}
}

func TestXMap_Delete(t *testing.T) {
	x := &XMap[int, string]{}
	x.Store(1, "value")

	x.Delete(1)

	_, ok := x.Load(1)
	if ok {
		t.Errorf("Delete() failed to remove key")
	}
}

func TestXMap_Store(t *testing.T) {
	x := &XMap[int, string]{}
	x.Store(1, "value")

	val, ok := x.Load(1)
	if !ok {
		t.Errorf("Store() failed to store value")
	}
	if val != "value" {
		t.Errorf("Store() stored incorrect value")
	}
}
