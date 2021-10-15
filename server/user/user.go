package user

// User is a user model.
type User struct {
	Uid    int    `json:"uid" db:"uid"`
	Name   string `json:"name" db:"name"`
	Avatar string `json:"avatar" db:"avatar"`
}

type Repository interface {
	// Get a user by uid.
	Get(uid int) (User, error)

	// Create a user, and return this user model.
	Create(name, avatar string) (User, error)
}
