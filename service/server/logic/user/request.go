package user

type NonceReq struct {
	Address string `form:"address" binding:"required"`
}

type LoginReq struct {
	Message   string `json:"message" binding:"required"`
	Signature string `json:"signature" binding:"required"`
}
