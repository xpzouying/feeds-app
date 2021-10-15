package feed

type Status uint8

const (
	StatusDeleted = iota
	StatusNormal
)

type Feed struct {
	ID       int    `json:"id" db:"id"`
	AuthorID int    `json:"author_id" db:"uid"`
	Text     string `json:"text" db:"text"`
	Status   Status `json:"status" db:"status"`
}

type Repository interface {
	ListFeeds(page, count int) ([]Feed, error)

	PostFeed(uid int, text string) (Feed, error)
}
