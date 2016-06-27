package kv

func NewKVChe() *KVChe {
	ch := make(chan map[string]interface{}, 1)
	ch <- make(map[string]interface{})
	return &KVChe{ch}
}

type KVChe struct {
	ch chan map[string]interface{}
}

func (kv *KVChe) Set(key string, val interface{}) {
	m := <-kv.ch
	m[key] = val
	kv.ch <- m
}

func (kv *KVChe) Get(key string) interface{} {
	m := <-kv.ch

	defer func() {
		select {
		case kv.ch <- m:
		default:
		}
	}()

	return m[key]
}
