package service

import (
	"context"

	"github.com/cirius-go/gojen/example/crud/internal/dto"
	"github.com/cirius-go/gojen/example/crud/internal/repo/uow"
)

// User represents the user service.
type User struct {
	uow uow.UnitOfWork
}

// NewUser returns a new User instance.
func NewUser(uow uow.UnitOfWork) *User {
	return &User{
		uow: uow,
	}
}

// Create creates a new user.
func (s *User) Create(ctx context.Context, req *dto.CreateUserReq) (*dto.CreateUserResp, error) {
	panic("not implemented")
}

// Update updates a user.
func (s *User) Update(ctx context.Context, req *dto.UpdateUserReq) (*dto.UpdateUserResp, error) {
	panic("not implemented")
}

// Get gets a user.
func (s *User) Get(ctx context.Context, req *dto.GetUserReq) (*dto.GetUserResp, error) {
	panic("not implemented")
}

// List lists users.
func (s *User) List(ctx context.Context, req *dto.ListUsersReq) (*dto.ListUsersResp, error) {
	panic("not implemented")
}

// Delete deletes a user.
func (s *User) Delete(ctx context.Context, req *dto.DeleteUserReq) (*dto.DeleteUserResp, error) {
	panic("not implemented")
}
