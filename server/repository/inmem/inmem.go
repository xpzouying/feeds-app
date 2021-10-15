package inmem

import (
	"errors"
	"sync"

	"github.com/xpzouying/feeds-app/server/feed"
	"github.com/xpzouying/feeds-app/server/user"
)

type feedRepository struct {
	lock  sync.Mutex
	feeds []feed.Feed

	// lastFeedID 记录一下最后一次分配的feed id
	lastFeedID int
}

func NewFeedRepository() feed.Repository {
	sampleFeeds := feed.SampleFeeds
	maxID := 0
	for _, feed := range sampleFeeds {
		if feed.ID > maxID {
			maxID = feed.ID
		}
	}

	return &feedRepository{
		feeds:      feed.SampleFeeds,
		lastFeedID: maxID,
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

func (r *feedRepository) PostFeed(uid int, text string) (feed.Feed, error) {
	newFeed := feed.Feed{
		ID:       r.getNewFeedID(),
		AuthorID: uid,
		Text:     text,
	}

	r.lock.Lock()
	r.feeds = append(r.feeds, newFeed)
	r.lock.Unlock()
	return newFeed, nil
}

func (r *feedRepository) getNewFeedID() (feedID int) {
	r.lock.Lock()
	defer r.lock.Unlock()

	feedID = r.lastFeedID + 1
	r.lastFeedID++
	return
}

type userRepository struct {
	lock  sync.Mutex
	users map[int]user.User

	lastUid int // the next new uid will be: lastUid+1
}

func NewUserRepository() user.Repository {
	var (
		sampleUsers = user.SampleUsers
		maxUid      = 0
	)
	for _, user := range sampleUsers {
		if user.Uid > maxUid {
			maxUid = user.Uid
		}
	}

	return &userRepository{
		users:   sampleUsers,
		lastUid: maxUid,
	}
}

// Get a user by uid.
func (r *userRepository) Get(uid int) (user.User, error) {
	r.lock.Lock()
	defer r.lock.Unlock()

	if u, ok := r.users[uid]; ok {
		return u, nil
	} else {
		return user.User{}, errors.New("user not exists")
	}
}

// Create a user, and return this user model.
func (r *userRepository) Create(name, avatar string) (user.User, error) {
	newUser := user.User{
		Uid:    r.getNewUserID(),
		Name:   name,
		Avatar: avatar,
	}

	r.lock.Lock()
	r.users[newUser.Uid] = newUser
	r.lock.Unlock()

	return newUser, nil
}

func (r *userRepository) getNewUserID() (uid int) {
	r.lock.Lock()
	defer r.lock.Unlock()

	uid = r.lastUid + 1
	r.lastUid++
	return
}
