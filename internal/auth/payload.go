package auth

type InitAuthRequest struct {
	Phone string `json:"phone" validate:"required,min=11,max=15"`
}

type InitAuthResponse struct {
	SessionID string `json:"sessionId"`
}

type VerifyCodeRequest struct {
	SessionID string `json:"sessionId" validate:"required"`
	Code      string `json:"code" validate:"required,len=4"`
}

type VerifyCodeResponse struct {
	Token string `json:"token"`
}
