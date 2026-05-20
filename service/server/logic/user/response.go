package user

type NonceResponse struct {
	Nonce string `json:"nonce"`
}

type LoginResponse struct {
	JWT     string `json:"jwt"`
	UserID  int64  `json:"user_id"`
	Address string `json:"address"`
}

type ProfileResponse struct {
	ID      int64  `json:"id"`
	Address string `json:"address"`
}
