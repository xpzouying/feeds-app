package feeding

import "github.com/xpzouying/feeds-app/server/feed"

type Service interface {
	ListFeeds(page, count int) []feed.Feed
}

type service struct {
	feedRepo feed.Repository
}

func NewService(feedRepo feed.Repository) Service {
	return &service{feedRepo}
}

func (s *service) ListFeeds(page, count int) []feed.Feed {
	feeds, err := s.feedRepo.ListFeeds(page, count)
	if err != nil {
		return []feed.Feed{}
	}
	return feeds
}
