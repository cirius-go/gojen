package service

import (
	"context"

	"github.com/cirius-go/gojen/example/crud/internal/dto"
	"github.com/cirius-go/gojen/example/crud/internal/repo/uow"
)

// Otp represents the otp service.
type Otp struct {
	uow uow.UnitOfWork
}

// NewOtp returns a new Otp instance.
func NewOtp(uow uow.UnitOfWork) *Otp {
	return &Otp{
		uow: uow,
	}
}

// Create creates a new otp.
func (s *Otp) Create(ctx context.Context, req *dto.CreateOtpReq) (*dto.CreateOtpResp, error) {
	panic("not implemented")
}

// Update updates a otp.
func (s *Otp) Update(ctx context.Context, req *dto.UpdateOtpReq) (*dto.UpdateOtpResp, error) {
	panic("not implemented")
}

// Get gets a otp.
func (s *Otp) Get(ctx context.Context, req *dto.GetOtpReq) (*dto.GetOtpResp, error) {
	panic("not implemented")
}

// List lists otps.
func (s *Otp) List(ctx context.Context, req *dto.ListOtpsReq) (*dto.ListOtpsResp, error) {
	panic("not implemented")
}

// Delete deletes a otp.
func (s *Otp) Delete(ctx context.Context, req *dto.DeleteOtpReq) (*dto.DeleteOtpResp, error) {
	panic("not implemented")
}
