package appmodels

type LoginReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResp struct {
	Token string `json:"token"`
}

type ValidateTokenReq struct {
	Token string `json:"token"`
}

type ValidateTokenResp struct {
	ID int64 `json:"id"`
}
