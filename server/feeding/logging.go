package feeding

import (
	"time"

	"github.com/go-kit/log"
)

type loggingMiddleware struct {
	logger log.Logger
	next   Service
}

func WithLoggingMiddleware(logger log.Logger) Middleware {

	return func(next Service) Service {
		return &loggingMiddleware{logger, next}
	}
}

func (lm loggingMiddleware) ListFeeds(page, count int) (feeds []Feed) {
	defer func(begin time.Time) {
		lm.logger.Log(
			"method", "list_feeds",
			"page", page,
			"count", count,
			"feeds_ret_len", len(feeds),
			"time_used", time.Since(begin).Milliseconds(),
		)
	}(time.Now())

	feeds = lm.next.ListFeeds(page, count)
	return
}

func (lm loggingMiddleware) PostFeed(uid int, text string) (Feed, error) {
	defer func(begin time.Time) {
		lm.logger.Log(
			"method", "post_feed",
			"uid", uid,
			"text", text,
			"time_used", time.Since(begin).Milliseconds(),
		)
	}(time.Now())

	return lm.next.PostFeed(uid, text)
}
