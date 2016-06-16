package kv

import (
	"github.com/pborman/uuid"
	"testing"
)

type KeyValueStore interface {
	Set(key string, v interface{})
	Get(key string) interface{}
}

var fake = &struct{ somefield string }{`ALOXA`}

func generateData(count int, store KeyValueStore) {
	for i := 0; i < count; i++ {
		store.Set(uuid.New(), fake)
	}
}

func TestKVC_Get(t *testing.T) {
	obj := &struct{ somefield string }{`ALOXA`}

	KVChannelStore.Set(`key`, obj)
	KVChannelStore.Set(`key-1`, `data-1`)
	KVChannelStore.Set(`key-2`, `data-2`)

	data := KVChannelStore.Get(`key`)
	if data == nil {
		t.Fatalf(`method get must return data by exists key`)
	}

	data = KVChannelStore.Get(`key-1`)
	if data == nil || data.(string) != `data-1` {
		t.Fatalf(`method get must return data by exists key`)
	}

	data = KVChannelStore.Get(`key-2`)
	if data == nil || data.(string) != `data-2` {
		t.Fatalf(`method get must return data by exists key`)
	}
}

func BenchmarkKV_Set(b *testing.B) {
	obj := &struct{ somefield string }{`ALOXA`}
	generateData(1000000, KVStore)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			KVStore.Set(`somekey`, obj)
		}
	})
}

func BenchmarkKV_Get(b *testing.B) {
	obj := &struct{ somefield string }{`ALOXA`}
	key := uuid.New()
	generateData(1000000, KVRStore)
	KVStore.Set(key, obj)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			KVStore.Get(key)
		}
	})
}

func BenchmarkKV_Channel_Set(b *testing.B) {
	obj := &struct{ somefield string }{`ALOXA`}
	generateData(1000000, KVChannelStore)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			KVChannelStore.Set(`somekey`, obj)
		}
	})
}

func BenchmarkKV_Channel_Get(b *testing.B) {
	obj := &struct{ somefield string }{`ALOXA`}
	key := uuid.New()
	generateData(1000000, KVChannelStore)
	KVChannelStore.Set(key, obj)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			KVChannelStore.Get(key)
		}
	})
}

func BenchmarkKVR_Get(b *testing.B) {
	obj := &struct{ somefield string }{`ALOXA`}
	key := uuid.New()
	generateData(1000000, KVRStore)
	KVRStore.Set(key, obj)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			KVRStore.Get(key)
		}
	})
}

func BenchmarkMap_Successively_Set(b *testing.B) {
	obj := &struct{ somefield string }{`ALOXA`}
	store := make(map[int]interface{})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		store[i] = obj
	}
}
