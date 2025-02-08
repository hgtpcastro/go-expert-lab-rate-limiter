package limiter

import "context"

type Store interface {
	Get(ctx context.Context, key string, rate Rate) (Context, error)
	Peek(ctx context.Context, key string, rate Rate) (Context, error)
	Reset(ctx context.Context, key string, rate Rate) (Context, error)
	Inc(ctx context.Context, key string, count int64, rate Rate) (Context, error)
}

type StoreOptions struct {
	Prefix string
}
