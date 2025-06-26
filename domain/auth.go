package domain

type AuthRequest struct {
	Method string
	Path   string
	Token  string
}

type UserRequestDTO struct {
	Email           string `json:"email" validate:"required"`
	Username        string `json:"username" validate:"required"`
	Password        string `json:"password" `
	ConfirmPassword string `json:"confirm_password" validate:"required"`
	PhoneNumber     string `json:"phone_number" validate:"required"`
	FullName        string `json:"full_name" validate:"required"`
	Gender          string `json:"gender" validate:"required"`
}
