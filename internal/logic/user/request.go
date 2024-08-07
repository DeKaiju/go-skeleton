package user

type LoginReq struct {
	Message   string `json:"message"`
	Signature string `json:"signature"`
}
