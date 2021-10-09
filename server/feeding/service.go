package feeding

import (
	"github.com/xpzouying/feeds-app/server/feed"
	"github.com/xpzouying/feeds-app/server/user"
)

type Middleware func(Service) Service

type Service interface {
	ListFeeds(page, count int) []feed.Feed

	PostFeed(uid int, text string) (feed.Feed, error)
}

type service struct {
	feedRepo feed.Repository
	userRepo user.Repository
}

func NewService(feedRepo feed.Repository, userRepo user.Repository) Service {
	return &service{feedRepo, userRepo}
}

func (s *service) ListFeeds(page, count int) []feed.Feed {
	feeds, err := s.feedRepo.ListFeeds(page, count)
	if err != nil {
		return []feed.Feed{}
	}
	return feeds
}

func (s *service) PostFeed(uid int, text string) (feed.Feed, error) {
	_, err := s.userRepo.Get(uid)
	if err != nil {
		return feed.Feed{}, err
	}

	return s.feedRepo.PostFeed(uid, text)
}
