package kv

import "sync"

type KV struct {
	l sync.Mutex
	s map[string]interface{}
}

type KVR struct {
	l sync.RWMutex
	s map[string]interface{}
}

func (k *KV) Set(key string, val interface{}) {
	k.l.Lock()
	k.s[key] = val
	k.l.Unlock()
}

func (k *KVR) Set(key string, val interface{}) {
	k.l.Lock()
	k.s[key] = val
	k.l.Unlock()
}

func (k *KV) Get(key string) interface{} {
	k.l.Lock()
	defer k.l.Unlock()
	return k.s[key]
}

func (k *KVR) Get(key string) interface{} {
	k.l.RLock()
	defer k.l.RUnlock()
	return k.s[key]
}

var KVStore = &KV{s: make(map[string]interface{}, 0)}

var KVRStore = &KVR{s: make(map[string]interface{}, 0)}
