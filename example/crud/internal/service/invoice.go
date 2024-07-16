package service

import (
	"context"

	"github.com/cirius-go/gojen/example/crud/internal/dto"
	"github.com/cirius-go/gojen/example/crud/internal/repo/uow"
)

// Invoice represents the invoice service.
type Invoice struct {
	uow uow.UnitOfWork
}

// NewInvoice returns a new Invoice instance.
func NewInvoice(uow uow.UnitOfWork) *Invoice {
	return &Invoice{
		uow: uow,
	}
}

// Create creates a new invoice.
func (s *Invoice) Create(ctx context.Context, req *dto.CreateInvoiceReq) (*dto.CreateInvoiceResp, error) {
	panic("not implemented")
}

// Update updates a invoice.
func (s *Invoice) Update(ctx context.Context, req *dto.UpdateInvoiceReq) (*dto.UpdateInvoiceResp, error) {
	panic("not implemented")
}

// Get gets a invoice.
func (s *Invoice) Get(ctx context.Context, req *dto.GetInvoiceReq) (*dto.GetInvoiceResp, error) {
	panic("not implemented")
}

// List lists invoices.
func (s *Invoice) List(ctx context.Context, req *dto.ListInvoicesReq) (*dto.ListInvoicesResp, error) {
	panic("not implemented")
}

// Delete deletes a invoice.
func (s *Invoice) Delete(ctx context.Context, req *dto.DeleteInvoiceReq) (*dto.DeleteInvoiceResp, error) {
	panic("not implemented")
}
