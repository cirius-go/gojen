package dto

type (
	// CreateStaffReq represents the request to create a new staff.
	// swagger:model
	CreateStaffReq struct {
	}

	// CreateStaffResp represents the response to create a new staff.
	// swagger:model
	CreateStaffResp struct {
	}
)

type (
	// UpdateStaffReq represents the request to update a staff.
	// swagger:model
	UpdateStaffReq struct {
		ID string `param:"id"`
	}

	// UpdateStaffResp represents the response to update a staff.
	// swagger:model
	UpdateStaffResp struct {
	}
)

type (
	// GetStaffReq represents the request to get a staff.
	// swagger:model
	GetStaffReq struct {
		ID string `param:"id"`
	}

	// GetStaffResp represents the response to get a staff.
	// swagger:model
	GetStaffResp struct {
	}
)

type (
	// ListStaffsReq represents the request to list staffs.
	// swagger:model
	ListStaffsReq struct {
	}

	// ListStaffsResp represents the response to list staffs.
	// swagger:model
	ListStaffsResp struct {
	}
)

type (
	// DeleteStaffReq represents the request to delete a staff.
	// swagger:model
	DeleteStaffReq struct {
		ID string `param:"id"`
	}

	// DeleteStaffResp represents the response to delete a staff.
	// swagger:model
	DeleteStaffResp struct {
	}
)
