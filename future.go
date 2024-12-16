package future

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
)

type tOrError[T any] struct {
	Value T
	Error error
}

type Future[T any] struct {
	resCh   chan *tOrError[T]
	resVal  atomic.Pointer[tOrError[T]]
	setOnce sync.Once
}

var ErrAlreadySet = errors.New("future already set")

func (f *Future[T]) Set(value T, err error) error {
	alreadySet := true
	f.setOnce.Do(func() {
		f.resCh <- &tOrError[T]{Value: value, Error: err}
		close(f.resCh)
		alreadySet = false
	})
	if alreadySet {
		return ErrAlreadySet
	}
	return nil
}

func (f *Future[T]) GetWithContext(ctx context.Context) (T, error) {
	select {
	case <-ctx.Done():
		var null T
		return null, ctx.Err()
	case res, ok := <-f.resCh:
		return f.onGet(res, ok)
	}
}

func (f *Future[T]) onGet(res *tOrError[T], ok bool) (T, error) {
	if ok {
		f.resVal.Swap(res)
	} else {
		for res = f.resVal.Load(); res == nil; res = f.resVal.Load() {
		}
	}
	return res.Value, res.Error
}

func (f *Future[T]) Get() (T, error) {
	res, ok := <-f.resCh
	return f.onGet(res, ok)
}

func New[T any]() *Future[T] {
	return &Future[T]{
		resCh: make(chan *tOrError[T], 1),
	}
}
