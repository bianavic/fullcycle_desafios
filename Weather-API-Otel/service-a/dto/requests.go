package dto

type CEPRequest struct {
	CEP string `json:"cep" validate:"required,len=8,numeric"`
}
