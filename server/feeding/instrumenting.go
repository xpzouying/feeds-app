package feeding

import (
	"time"

	"github.com/go-kit/kit/metrics"

	"github.com/xpzouying/feeds-app/server/feed"
)

type instrumentingMiddleware struct {
	reqCounter metrics.Counter
	next       Service
}

func InstrumentMiddleware(counter metrics.Counter) Middleware {
	return func(next Service) Service {
		return &instrumentingMiddleware{
			reqCounter: counter,
			next:       next,
		}
	}
}

func (im instrumentingMiddleware) ListFeeds(page, count int) (feeds []feed.Feed) {
	defer func(begin time.Time) {
		im.reqCounter.With("method", "list_feeds").Add(1)
	}(time.Now())

	feeds = im.next.ListFeeds(page, count)
	return
}
