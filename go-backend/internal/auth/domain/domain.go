package domain

type SignupReq struct {
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,gte=8,lte=20"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

type LoginReq struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,gte=8,lte=20"`
}

type AuthResp struct {
	AccessToken string `json:"accessToken"`
	Type        string `json:"type"`
}
