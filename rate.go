package limiter

import "time"

type Rate struct {
	Period time.Duration
	Limit  int64
}

func NewRate(limit int64, period int) Rate {
	rate := Rate{
		Limit:  limit,
		Period: time.Duration(period) * time.Second,
	}

	return rate
}
