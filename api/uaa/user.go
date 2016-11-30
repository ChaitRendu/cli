package uaa

// User represents an UAA user account.
type User struct {
	ID string
}

// NewUser creates a new UAA user account with the provided password.
func (client *Client) NewUser(username string, password string) (User, error) {
	return User{}, nil
}
