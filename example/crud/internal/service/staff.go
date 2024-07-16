package service

import (
	"context"

	"github.com/cirius-go/gojen/example/crud/internal/dto"
	"github.com/cirius-go/gojen/example/crud/internal/repo/uow"
)

// Staff represents the staff service.
type Staff struct {
	uow uow.UnitOfWork
}

// NewStaff returns a new Staff instance.
func NewStaff(uow uow.UnitOfWork) *Staff {
	return &Staff{
		uow: uow,
	}
}

// Create creates a new staff.
func (s *Staff) Create(ctx context.Context, req *dto.CreateStaffReq) (*dto.CreateStaffResp, error) {
	panic("not implemented")
}

// Update updates a staff.
func (s *Staff) Update(ctx context.Context, req *dto.UpdateStaffReq) (*dto.UpdateStaffResp, error) {
	panic("not implemented")
}

// Get gets a staff.
func (s *Staff) Get(ctx context.Context, req *dto.GetStaffReq) (*dto.GetStaffResp, error) {
	panic("not implemented")
}

// List lists staffs.
func (s *Staff) List(ctx context.Context, req *dto.ListStaffsReq) (*dto.ListStaffsResp, error) {
	panic("not implemented")
}

// Delete deletes a staff.
func (s *Staff) Delete(ctx context.Context, req *dto.DeleteStaffReq) (*dto.DeleteStaffResp, error) {
	panic("not implemented")
}
