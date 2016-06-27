package kv

func NewKVChm() *KVChm {
	chr := make(chan struct{}, 1)
	chw := make(chan struct{}, 1)
	chr <- struct{}{}
	chw <- struct{}{}
	return &KVChm{make(map[string]interface{}), chr, chw}
}

type KVChm struct {
	m   map[string]interface{}
	chr chan struct{}
	chw chan struct{}
}

func (kv *KVChm) Set(key string, val interface{}) {
	<-kv.chw
	kv.m[key] = val
	kv.chw <- struct{}{}
}

func (kv *KVChm) Get(key string) interface{} {
	<-kv.chr
	val := kv.m[key]
	kv.chr <- struct{}{}
	return val
}
