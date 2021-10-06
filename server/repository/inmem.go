package repository

import (
	"sync"

	"github.com/xpzouying/feeds-app/server/feed"
)

type feedRepository struct {
	lock  sync.Mutex
	feeds []feed.Feed
}

func NewFeedRepository() feed.Repository {
	return &feedRepository{
		feeds: feed.SampleFeeds,
	}
}

func (r *feedRepository) ListFeeds(page, count int) ([]feed.Feed, error) {
	r.lock.Lock()
	defer r.lock.Unlock()

	var (
		begin = page * count
		end   = begin + count
	)

	feedsCount := len(r.feeds)
	if begin >= feedsCount {
		return []feed.Feed{}, nil
	}
	if end > feedsCount {
		end = feedsCount
	}

	return r.feeds[begin:end], nil
}
