package api

// User represents the user API.
type User struct {
	svc UserSvc
}

// NewUser returns a new User instance.
func NewUser(svc UserSvc) *User {
	return &User{
		svc: svc,
	}
}
