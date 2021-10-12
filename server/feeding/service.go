package feeding

import (
	"github.com/xpzouying/feeds-app/server/feed"
	"github.com/xpzouying/feeds-app/server/user"
)

type Middleware func(Service) Service

type Service interface {
	ListFeeds(page, count int) []Feed

	PostFeed(uid int, text string) (Feed, error)
}

type Feed struct {
	feed.Feed

	AuthorName   string `json:"author_name"`
	AuthorAvatar string `json:"author_avatar"`
}

type service struct {
	feedRepo feed.Repository
	userRepo user.Repository
}

func NewService(feedRepo feed.Repository, userRepo user.Repository) Service {
	return &service{feedRepo, userRepo}
}

func (s *service) ListFeeds(page, count int) []Feed {

	feedsList, err := s.feedRepo.ListFeeds(page, count)
	if err != nil {
		return []Feed{}
	}

	feeds := make([]Feed, 0, len(feedsList))
	for _, f := range feedsList {
		author, err := s.userRepo.Get(f.AuthorID)
		if err != nil {
			continue
		}

		feeds = append(feeds, Feed{
			Feed:         f,
			AuthorName:   author.Name,
			AuthorAvatar: author.Avatar,
		})
	}

	return feeds
}

func (s *service) PostFeed(uid int, text string) (Feed, error) {
	author, err := s.userRepo.Get(uid)
	if err != nil {
		return Feed{}, err
	}

	newFeed, err := s.feedRepo.PostFeed(uid, text)
	if err != nil {
		return Feed{}, err
	}

	return Feed{
		Feed:         newFeed,
		AuthorName:   author.Name,
		AuthorAvatar: author.Avatar,
	}, nil
}
