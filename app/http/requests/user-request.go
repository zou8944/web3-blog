package requests

type CreateUserRequest struct {
	UniqueName    string `json:"unique_name" validate:"required"`
	PublicAddress string `json:"public_address" validate:"required"`
}

type UpdateUserRequest struct {
	UniqueName string `json:"unique_name" validate:"required"`
}
