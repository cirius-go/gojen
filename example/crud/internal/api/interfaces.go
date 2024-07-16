package api

import (
	"context"

	"github.com/cirius-go/gojen/example/crud/internal/dto"
)

// UserSvc represents the user service.
type UserSvc interface {
	Create(ctx context.Context, req *dto.CreateUserReq) (*dto.CreateUserResp, error)
	Update(ctx context.Context, req *dto.UpdateUserReq) (*dto.UpdateUserResp, error)
	Get(ctx context.Context, req *dto.GetUserReq) (*dto.GetUserResp, error)
	List(ctx context.Context, req *dto.ListUsersReq) (*dto.ListUsersResp, error)
	Delete(ctx context.Context, req *dto.DeleteUserReq) (*dto.DeleteUserResp, error)
}

// StaffSvc represents the staff service.
type StaffSvc interface {
	Create(ctx context.Context, req *dto.CreateStaffReq) (*dto.CreateStaffResp, error)
	Update(ctx context.Context, req *dto.UpdateStaffReq) (*dto.UpdateStaffResp, error)
	Get(ctx context.Context, req *dto.GetStaffReq) (*dto.GetStaffResp, error)
	List(ctx context.Context, req *dto.ListStaffsReq) (*dto.ListStaffsResp, error)
	Delete(ctx context.Context, req *dto.DeleteStaffReq) (*dto.DeleteStaffResp, error)
}

// InvoiceSvc represents the invoice service.
type InvoiceSvc interface {
	Create(ctx context.Context, req *dto.CreateInvoiceReq) (*dto.CreateInvoiceResp, error)
	Update(ctx context.Context, req *dto.UpdateInvoiceReq) (*dto.UpdateInvoiceResp, error)
	Get(ctx context.Context, req *dto.GetInvoiceReq) (*dto.GetInvoiceResp, error)
	List(ctx context.Context, req *dto.ListInvoicesReq) (*dto.ListInvoicesResp, error)
	Delete(ctx context.Context, req *dto.DeleteInvoiceReq) (*dto.DeleteInvoiceResp, error)
}

// OtpSvc represents the otp service.
type OtpSvc interface {
	Create(ctx context.Context, req *dto.CreateOtpReq) (*dto.CreateOtpResp, error)
	Update(ctx context.Context, req *dto.UpdateOtpReq) (*dto.UpdateOtpResp, error)
	Get(ctx context.Context, req *dto.GetOtpReq) (*dto.GetOtpResp, error)
	List(ctx context.Context, req *dto.ListOtpsReq) (*dto.ListOtpsResp, error)
	Delete(ctx context.Context, req *dto.DeleteOtpReq) (*dto.DeleteOtpResp, error)
}
