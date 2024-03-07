package domain

type User struct {
	Email     string `json:"email"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

type ResetPasswordCodeReq struct {
	Email string `json:"email" validate:"required"`
}

type ValidateResetPasswordReq struct {
	Email string `json:"email" validate:"required"`
	Code  string `json:"code" validate:"required,len=6"`
}

type ResetPasswordReq struct {
	Email           string `json:"email" validate:"required"`
	NewPassword     string `json:"newPassword" validate:"required,gte=8,lte=20"`
	ConfirmPassword string `json:"confirmPassword" validate:"required,gte=8,lte=20,eqfield=NewPassword"`
}
