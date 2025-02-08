package stdlib

import (
	"net/http"
	"strconv"
	"strings"

	limiter "github.com/hgtpcastro/go-expert-lab-rate-limiter"
)

type Middleware struct {
	Limiter        *limiter.Limiter
	OnError        ErrorHandler
	OnLimitReached LimitReachedHandler
	KeyGetter      KeyGetter
}

func NewMiddleware(limiter *limiter.Limiter, options ...Option) *Middleware {
	middleware := &Middleware{
		Limiter:        limiter,
		OnError:        WithDefaultErrorHandler,
		OnLimitReached: WithDefaultLimitReachedHandler,
		KeyGetter:      WithIPKeyGetter(limiter),
	}

	for _, option := range options {
		option.apply(middleware)
	}

	return middleware
}

func (middleware *Middleware) Handler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := middleware.KeyGetter(r)

		if strings.TrimSpace(key) == "" {
			h.ServeHTTP(w, r)
			return
		}

		context, err := middleware.Limiter.Get(r.Context(), key)
		if err != nil {
			middleware.OnError(w, r, err)
			return
		}

		w.Header().Add("X-RateLimit-Limit", strconv.FormatInt(context.Limit, 10))
		w.Header().Add("X-RateLimit-Remaining", strconv.FormatInt(context.Remaining, 10))
		w.Header().Add("X-RateLimit-Reset", strconv.FormatInt(context.Reset, 10))

		if context.Reached {
			middleware.OnLimitReached(w, r)
			return
		}

		//h.ServeHTTP(w, r)
	})
}
