package kv

import (
	"golang.org/x/net/context"
)

func NewKVCtx() *KVCtx {
	ch := make(chan context.Context, 1)
	ch <- context.TODO()
	return &KVCtx{ch}
}

type KVCtx struct {
	ch chan context.Context
}

func (kv *KVCtx) Set(key string, val interface{}) {
	kv.ch <- context.WithValue(<-kv.ch, key, val)
}

func (kv *KVCtx) Get(key string) interface{} {
	ctx, ok := <-kv.ch
	if !ok {
		return nil
	}

	defer func() {
		select {
		case kv.ch <- ctx:
		default:
		}
	}()

	return ctx.Value(key)
}
