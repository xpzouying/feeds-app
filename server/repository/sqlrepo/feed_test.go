package sqlrepo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	mysqlOption = Option{
		User:     "zy",
		Password: "123456",

		Host: "docker.zy.local",
		Port: 3306,
	}
)

func TestNewFeedRepo(t *testing.T) {
	_, err := NewFeedRepo(mysqlOption)
	assert.NoError(t, err)
}

func TestListFeeds(t *testing.T) {
	repo, err := NewFeedRepo(mysqlOption)
	assert.NoError(t, err)

	feeds, err := repo.ListFeeds(0, 10)
	assert.NoError(t, err)

	t.Logf("list feeds: %+v", feeds)
}

func TestPostFeed(t *testing.T) {
	repo, err := NewFeedRepo(mysqlOption)
	assert.NoError(t, err)

	var (
		uid      = 1
		feedText = "hello world"
	)
	newFeed, err := repo.PostFeed(uid, feedText)
	assert.NoError(t, err)

	{
		defer removeFeedByID(t, newFeed.ID)
	}

	{
		assert.NotZero(t, newFeed.ID)
		assert.Equal(t, uid, newFeed.AuthorID)
		assert.Equal(t, feedText, newFeed.Text)
	}
}

func removeFeedByID(t *testing.T, feedID int) error {
	db, err := newMysqlDB(mysqlOption)
	assert.NoError(t, err)
	defer db.Close()

	result, err := db.Exec("DELETE FROM feeds where id = ?", feedID)
	assert.NoError(t, err)

	affectedRows, _ := result.RowsAffected()
	t.Logf("deleted row_id=%d affected_rows=%d", feedID, affectedRows)
	return nil
}
