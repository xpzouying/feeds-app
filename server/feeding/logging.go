package feeding

import (
	"time"

	"github.com/go-kit/log"

	"github.com/xpzouying/feeds-app/server/feed"
)

type loggingMiddleware struct {
	logger log.Logger
	next   Service
}

func LoggingMiddleware(logger log.Logger) Middleware {

	return func(next Service) Service {
		return &loggingMiddleware{logger, next}
	}
}

func (lm loggingMiddleware) ListFeeds(page, count int) (feeds []feed.Feed) {
	defer func(begin time.Time) {
		lm.logger.Log(
			"method", "list_feeds",
			"page", page,
			"count", count,
			"feeds_ret_len", len(feeds),
			"time_used", time.Since(begin),
		)
	}(time.Now())

	feeds = lm.next.ListFeeds(page, count)
	return
}