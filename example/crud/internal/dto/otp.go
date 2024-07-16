package dto

type (
	// CreateOtpReq represents the request to create a new otp.
	// swagger:model
	CreateOtpReq struct {
	}

	// CreateOtpResp represents the response to create a new otp.
	// swagger:model
	CreateOtpResp struct {
	}
)

type (
	// UpdateOtpReq represents the request to update a otp.
	// swagger:model
	UpdateOtpReq struct {
		ID string `param:"id"`
	}

	// UpdateOtpResp represents the response to update a otp.
	// swagger:model
	UpdateOtpResp struct {
	}
)

type (
	// GetOtpReq represents the request to get a otp.
	// swagger:model
	GetOtpReq struct {
		ID string `param:"id"`
	}

	// GetOtpResp represents the response to get a otp.
	// swagger:model
	GetOtpResp struct {
	}
)

type (
	// ListOtpsReq represents the request to list otps.
	// swagger:model
	ListOtpsReq struct {
	}

	// ListOtpsResp represents the response to list otps.
	// swagger:model
	ListOtpsResp struct {
	}
)

type (
	// DeleteOtpReq represents the request to delete a otp.
	// swagger:model
	DeleteOtpReq struct {
		ID string `param:"id"`
	}

	// DeleteOtpResp represents the response to delete a otp.
	// swagger:model
	DeleteOtpResp struct {
	}
)
