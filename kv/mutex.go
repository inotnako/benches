package kv

import (
	"sync"
)

type KV struct {
	l sync.Mutex
	s map[string]interface{}
}

type KVR struct {
	l sync.RWMutex
	s map[string]interface{}
}

type d struct {
	k string
	v interface{}
}

func newSpy(key string) *spy {
	return &spy{
		k: key,
		v: make(chan interface{}),
	}
}

type spy struct {
	k string
	v chan interface{}
}

type KVC struct {
	in  chan *d
	s   map[string]interface{}
	out chan *spy
}

func NewKVC() *KVC {
	store := &KVC{
		s:   make(map[string]interface{}, 0),
		in:  make(chan *d),
		out: make(chan *spy),
	}

	go func() {
		for {
			select {
			case data := <-store.in:
				store.s[data.k] = data.v
			case s := <-store.out:
				s.v <- store.s[s.k]
			}
		}
	}()

	return store
}

func (k *KVC) Set(key string, val interface{}) {
	k.in <- &d{key, val}
}

func (k *KVC) Get(key string) interface{} {
	transfer := newSpy(key)
	k.out <- transfer
	return <-transfer.v
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

var KVChannelStore = NewKVC()
