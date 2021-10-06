package feed

type Feed struct {
	ID       int
	AutherID int
	Text     string
}

type Repository interface {
	ListFeeds(page, count int) ([]Feed, error)
}
