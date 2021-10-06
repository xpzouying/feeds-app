package user

// User is a user model.
type User struct {
	Uid    int
	Name   string
	Avatar string
}

type Repository interface {
	// Get a user by uid.
	Get(uid int) error

	// Create a user, and return this user model.
	Create(name, avatar string) (User, error)
}
