package requests

type LoginMetaMaskRequest struct {
	PublicAddress string `json:"public_address" validate:"required"`
	Signature     string `json:"signature" validate:"required"`
}
