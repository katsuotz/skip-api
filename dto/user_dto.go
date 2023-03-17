package dto

type RegistrasiUserRequest struct {
	Phone    string `json:"phone" form:"phone"`
	Username string `json:"username" form:"username"`
	Email    string `json:"email" form:"email" `
	Password string `json:"password" form:"password"`
}
