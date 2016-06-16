package kv

import (
	"github.com/pborman/uuid"
	"testing"
)

func BenchmarkKV_Set(b *testing.B) {
	obj := &struct{ somefield string }{`ALOXA`}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			KVStore.Set(uuid.New(), obj)
		}
	})
}

func BenchmarkKVR_Get(b *testing.B) {
	obj := &struct{ somefield string }{`ALOXA`}
	key := uuid.New()
	KVRStore.Set(key, obj)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			KVRStore.Get(key)
		}
	})
}

func BenchmarkKV_Get(b *testing.B) {
	obj := &struct{ somefield string }{`ALOXA`}
	key := uuid.New()
	KVStore.Set(key, obj)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			KVStore.Get(key)
		}
	})
}
