package stdlib

import (
	"net/http"
	"strings"

	limiter "github.com/hgtpcastro/go-expert-lab-rate-limiter"
)

type Option interface {
	apply(*Middleware)
}

type option func(*Middleware)

func (o option) apply(m *Middleware) {
	o(m)
}

type ErrorHandler func(w http.ResponseWriter, r *http.Request, err error)

func WithErrorHandler(h ErrorHandler) Option {
	return option(func(m *Middleware) {
		m.OnError = h
	})
}

func WithDefaultErrorHandler(w http.ResponseWriter, r *http.Request, err error) {
	panic(err)
}

type LimitReachedHandler func(w http.ResponseWriter, r *http.Request)

func WithLimitReachedHandler(h LimitReachedHandler) Option {
	return option(func(m *Middleware) {
		m.OnLimitReached = h
	})
}

func WithDefaultLimitReachedHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "you have reached the maximum number of requests or actions allowed within a certain time frame", http.StatusTooManyRequests)
}

type KeyGetter func(r *http.Request) string

func WithKeyGetter(h KeyGetter) Option {
	return option(func(m *Middleware) {
		m.KeyGetter = h
	})
}

func WithIPKeyGetter(l *limiter.Limiter) func(r *http.Request) string {
	return func(r *http.Request) string {
		if strings.TrimSpace(l.GetToken(r)) != "" {
			return ""
		} else {
			return l.GetIP(r).String()
		}
	}
}

func WithTokenKeyGetter(l *limiter.Limiter) func(r *http.Request) string {
	return func(r *http.Request) string {
		return l.GetToken(r)
	}
}

func WithTokenAndIPKeyGetter(l *limiter.Limiter) func(r *http.Request) string {
	return func(r *http.Request) string {
		apiKey := l.GetToken(r)

		if apiKey != "" {
			// TODO: personalizar limit por token
			return apiKey
		}

		// TODO: personalizar limit por IP
		return limiter.GetIP(r).String()
	}
}
