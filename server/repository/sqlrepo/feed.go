package sqlrepo

import (
	"fmt"

	"github.com/xpzouying/feeds-app/server/feed"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type feedRepo struct {
	db *sqlx.DB
}

type Option struct {
	User     string
	Password string

	Host string
	Port int
}

// NewFeedRepo 新建一个sql的链接。
// user:password@tcp(localhost:5555)/dbname?loc=Local&tls=skip-verify&autocommit=true
func NewFeedRepo(opt Option) (feed.Repository, error) {

	db, err := newMysqlDB(opt)
	if err != nil {
		return nil, err
	}

	return &feedRepo{
		db,
	}, nil
}

func newMysqlDB(opt Option) (*sqlx.DB, error) {
	return sqlx.Connect("mysql", formatToSNI(opt))
}

func formatToSNI(opt Option) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/zydb?loc=Local", opt.User, opt.Password, opt.Host, opt.Port)
}

func (repo *feedRepo) ListFeeds(page, count int) (feeds []feed.Feed, err error) {

	var (
		offset = page * count
		limit  = count
	)

	err = repo.db.Select(&feeds,
		"SELECT * FROM feeds WHERE status = ? LIMIT ? OFFSET ?",
		feed.StatusNormal, limit, offset)
	return
}

func (repo *feedRepo) PostFeed(uid int, text string) (feed.Feed, error) {

	result, err := repo.db.Exec(
		"INSERT INTO feeds(uid, text, status) VALUES (?, ?, ?)",
		uid, text, feed.StatusNormal)
	if err != nil {
		return feed.Feed{}, err
	}

	recordID, err := result.LastInsertId()
	if err != nil {
		return feed.Feed{}, err
	}

	var record feed.Feed
	err = repo.db.Get(&record, "SELECT * FROM feeds WHERE id = ?", recordID)

	return record, err
}
