package service

// CustomerAPI service.
type CustomerAPI struct {
}

// Please select option 2,3 to append service method for service 'CustomerAPI'.
// +gojen:append-template=service

// Create new 'CustomerAPI'.
func (s *CustomerAPI) Create(ctx context.Context, req *dto.CreateCustomerAPIReq) (*dto.CreateCustomerAPIRes, error) {
  panic("not implemented")
}

// List all 'CustomerApis'.
func (s *CustomerAPI) List(ctx context.Context, req *dto.ListCustomerApisReq) (*dto.ListCustomerApisRes, error) {
  panic("not implemented")
}

// Get one of 'CustomerApis'.
func (s *CustomerAPI) Get(ctx context.Context, req *dto.GetCustomerAPIReq) (*dto.GetCustomerAPIRes, error) {
  panic("not implemented")
}

// Update one of 'CustomerApis'.
func (s *CustomerAPI) Update(ctx context.Context, req *dto.UpdateCustomerAPIReq) (*dto.UpdateCustomerAPIRes, error) {
  panic("not implemented")
}

// Delete one of 'CustomerApis'.
func (s *CustomerAPI) Delete(ctx context.Context, req *dto.DeleteCustomerAPIReq) (*dto.DeleteCustomerAPIRes, error) {
  panic("not implemented")
}


