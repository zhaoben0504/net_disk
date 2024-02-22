package dto

type LoginRequest struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refreshToken"`
}

type UserDetailsRequest struct {
	Indentity string `json:"indentity"`
}

type UserDetailsResponse struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type EmailCodeRequest struct {
	Email string `json:"email"`
}

type EmailCodeResponse struct {
	Code string `json:"code"`
	Msg  string `json:"msg"`
}

type UserRegisterRequest struct {
	Name     string `json:"name"`
	Password string `json:"password"`
	Email    string `json:"email"`
	Code     string `json:"code"`
}

type UserRegisterResponse struct {
	Msg string `json:"msg"`
}
