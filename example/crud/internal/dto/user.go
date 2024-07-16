package dto

type (
	// CreateUserReq represents the request to create a new user.
	// swagger:model
	CreateUserReq struct {
	}

	// CreateUserResp represents the response to create a new user.
	// swagger:model
	CreateUserResp struct {
	}
)

type (
	// UpdateUserReq represents the request to update a user.
	// swagger:model
	UpdateUserReq struct {
		ID string `param:"id"`
	}

	// UpdateUserResp represents the response to update a user.
	// swagger:model
	UpdateUserResp struct {
	}
)

type (
	// GetUserReq represents the request to get a user.
	// swagger:model
	GetUserReq struct {
		ID string `param:"id"`
	}

	// GetUserResp represents the response to get a user.
	// swagger:model
	GetUserResp struct {
	}
)

type (
	// ListUsersReq represents the request to list users.
	// swagger:model
	ListUsersReq struct {
	}

	// ListUsersResp represents the response to list users.
	// swagger:model
	ListUsersResp struct {
	}
)

type (
	// DeleteUserReq represents the request to delete a user.
	// swagger:model
	DeleteUserReq struct {
		ID string `param:"id"`
	}

	// DeleteUserResp represents the response to delete a user.
	// swagger:model
	DeleteUserResp struct {
	}
)
