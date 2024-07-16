package dto

type (
	// CreateInvoiceReq represents the request to create a new invoice.
	// swagger:model
	CreateInvoiceReq struct {
	}

	// CreateInvoiceResp represents the response to create a new invoice.
	// swagger:model
	CreateInvoiceResp struct {
	}
)

type (
	// UpdateInvoiceReq represents the request to update a invoice.
	// swagger:model
	UpdateInvoiceReq struct {
		ID string `param:"id"`
	}

	// UpdateInvoiceResp represents the response to update a invoice.
	// swagger:model
	UpdateInvoiceResp struct {
	}
)

type (
	// GetInvoiceReq represents the request to get a invoice.
	// swagger:model
	GetInvoiceReq struct {
		ID string `param:"id"`
	}

	// GetInvoiceResp represents the response to get a invoice.
	// swagger:model
	GetInvoiceResp struct {
	}
)

type (
	// ListInvoicesReq represents the request to list invoices.
	// swagger:model
	ListInvoicesReq struct {
	}

	// ListInvoicesResp represents the response to list invoices.
	// swagger:model
	ListInvoicesResp struct {
	}
)

type (
	// DeleteInvoiceReq represents the request to delete a invoice.
	// swagger:model
	DeleteInvoiceReq struct {
		ID string `param:"id"`
	}

	// DeleteInvoiceResp represents the response to delete a invoice.
	// swagger:model
	DeleteInvoiceResp struct {
	}
)
