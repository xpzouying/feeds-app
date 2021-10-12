package feeding

import (
	"time"

	"github.com/go-kit/kit/metrics"
)

type instrumentingMiddleware struct {
	reqCounter metrics.Counter
	reqLatency metrics.Histogram
	next       Service
}

func InstrumentMiddleware(counter metrics.Counter, latency metrics.Histogram) Middleware {
	return func(next Service) Service {
		return &instrumentingMiddleware{counter, latency, next}
	}
}

func (im instrumentingMiddleware) ListFeeds(page, count int) (feeds []Feed) {
	defer func(begin time.Time) {
		im.reqCounter.With("method", "list_feeds").Add(1)
		im.reqLatency.With("method", "list_feeds").Observe(float64(time.Since(begin).Milliseconds()))
	}(time.Now())

	feeds = im.next.ListFeeds(page, count)
	return
}

func (im instrumentingMiddleware) PostFeed(uid int, text string) (Feed, error) {
	defer func(begin time.Time) {
		im.reqCounter.With("method", "post_feed").Add(1)
		im.reqLatency.With("method", "post_feed").Observe(float64(time.Since(begin).Milliseconds()))
	}(time.Now())

	return im.next.PostFeed(uid, text)
}
