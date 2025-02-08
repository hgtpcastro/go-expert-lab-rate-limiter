package limiter

import "context"

type Context struct {
	Limit     int64
	Remaining int64
	Reset     int64
	Reached   bool
}

type Limiter struct {
	Store Store
	Rate  Rate
}

func NewLimiter(store Store, rate Rate) *Limiter {
	return &Limiter{
		Store: store,
		Rate:  rate,
	}
}

func (l *Limiter) Get(ctx context.Context, key string) (Context, error) {
	return l.Store.Get(ctx, key, l.Rate)
}

func (l *Limiter) Peek(ctx context.Context, key string) (Context, error) {
	return l.Store.Peek(ctx, key, l.Rate)
}

func (l *Limiter) Reset(ctx context.Context, key string) (Context, error) {
	return l.Store.Reset(ctx, key, l.Rate)
}

func (l *Limiter) Inc(ctx context.Context, key string, count int64) (Context, error) {
	return l.Store.Inc(ctx, key, count, l.Rate)
}
