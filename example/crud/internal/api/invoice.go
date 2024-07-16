package api

// Invoice represents the invoice API.
type Invoice struct {
	svc InvoiceSvc
}

// NewInvoice returns a new Invoice instance.
func NewInvoice(svc InvoiceSvc) *Invoice {
	return &Invoice{
		svc: svc,
	}
}
