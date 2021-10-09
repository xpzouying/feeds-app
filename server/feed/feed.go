package feed

type Feed struct {
	ID       int    `json:"id"`
	AuthorID int    `json:"author_id"`
	Text     string `json:"text"`
}

type Repository interface {
	ListFeeds(page, count int) ([]Feed, error)

	PostFeed(uid int, text string) (Feed, error)
}
