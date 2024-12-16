package future

import (
	"context"
	"errors"
	"testing"
)

func TestFuture(t *testing.T) {
	f := New[int]()
	t.Run("get context", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		val, err := f.GetWithContext(ctx)
		if val != 0 || !errors.Is(err, context.Canceled) {
			t.Errorf("expected 0, context.Canceled; got %d, %v", val, err)
		}
	})

	t.Run("set", func(t *testing.T) {
		err := f.Set(42, nil)
		if err != nil {
			t.Errorf("expected nil; got %v", err)
		}
	})

	t.Run("set twice", func(t *testing.T) {
		err := f.Set(42, nil)
		if !errors.Is(err, ErrAlreadySet) {
			t.Errorf("expected ErrAlreadySet; got %v", err)
		}
	})

	t.Run("get", func(t *testing.T) {
		val, err := f.Get()
		if val != 42 || err != nil {
			t.Errorf("expected 42, nil; got %d, %v", val, err)
		}
	})

	t.Run("get twice", func(t *testing.T) {
		val, err := f.Get()
		if val != 42 || err != nil {
			t.Errorf("expected 42, nil; got %d, %v", val, err)
		}
	})

	t.Run("get context normal", func(t *testing.T) {
		ctx := context.Background()
		val, err := f.GetWithContext(ctx)
		if val != 42 || err != nil {
			t.Errorf("expected 42, nil; got %d, %v", val, err)
		}
	})
}
