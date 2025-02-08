package common

import (
	"time"

	limiter "github.com/hgtpcastro/go-expert-lab-rate-limiter"
)

func GetContextFromState(
	now time.Time,
	rate limiter.Rate,
	expiration time.Time,
	count int64,
) limiter.Context {
	limit := rate.Limit
	remaining := int64(0)
	reached := true

	if count <= limit {
		remaining = limit - count
		reached = false
	}

	reset := expiration.Unix()

	return limiter.Context{
		Limit:     limit,
		Remaining: remaining,
		Reset:     reset,
		Reached:   reached,
	}
}
