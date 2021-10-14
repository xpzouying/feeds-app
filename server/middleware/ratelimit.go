package middleware

import (
	"time"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/ratelimit"
	"golang.org/x/time/rate"
)

func WithRateLimiter(ep endpoint.Endpoint) endpoint.Endpoint {
	limiter := rate.NewLimiter(rate.Every(1*time.Second), 1)

	return ratelimit.NewErroringLimiter(limiter)(ep)
}
