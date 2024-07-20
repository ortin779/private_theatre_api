package models

type LoginParams struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (loginParams LoginParams) Validate() map[string]string {
	errs := make(map[string]string)

	if !isEmailValid(loginParams.Email) {
		errs["email"] = "invalid email format"
	}
	if loginParams.Password == "" {
		errs["password"] = "password can not be empty"
	}
	return errs
}

type LoginResponse struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}
