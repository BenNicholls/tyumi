package util

import "testing"

type testCacheData struct {
	Name  string
	Value int
}

func TestDataCacheRegisterAndGet(t *testing.T) {
	var cache DataCache[testCacheData, uint32]

	id1 := cache.RegisterDataType(testCacheData{Name: "first", Value: 1})
	id2 := cache.RegisterDataType(testCacheData{Name: "second", Value: 2})

	if id1 != 0 || id2 != 1 {
		t.Errorf("RegisterDataType returned ids (%d, %d), want (0, 1)", id1, id2)
	}

	got1 := cache.GetData(id1)
	if got1.Name != "first" || got1.Value != 1 {
		t.Errorf("GetData(%d) = %+v, want {first 1}", id1, got1)
	}

	got2 := cache.GetData(id2)
	if got2.Name != "second" || got2.Value != 2 {
		t.Errorf("GetData(%d) = %+v, want {second 2}", id2, got2)
	}
}

func TestDataCacheGetInvalidType(t *testing.T) {
	var cache DataCache[testCacheData, uint32]
	cache.RegisterDataType(testCacheData{Name: "only", Value: 1})

	got := cache.GetData(uint32(99))
	want := testCacheData{}
	if got != want {
		t.Errorf("GetData with invalid type = %+v, want zero value", got)
	}
}

func TestDataCacheReplaceData(t *testing.T) {
	var cache DataCache[testCacheData, uint32]
	id := cache.RegisterDataType(testCacheData{Name: "original", Value: 1})

	cache.ReplaceData(id, testCacheData{Name: "replaced", Value: 99})

	got := cache.GetData(id)
	if got.Name != "replaced" || got.Value != 99 {
		t.Errorf("GetData after ReplaceData = %+v, want {replaced 99}", got)
	}
}

func TestDataCacheReplaceInvalidType(t *testing.T) {
	var cache DataCache[testCacheData, uint32]
	cache.RegisterDataType(testCacheData{Name: "only", Value: 1})

	// should not panic
	cache.ReplaceData(uint32(99), testCacheData{Name: "nope", Value: 0})

	got := cache.GetData(uint32(0))
	if got.Name != "only" {
		t.Errorf("ReplaceData with invalid type should not affect valid entries, got %+v", got)
	}
}
