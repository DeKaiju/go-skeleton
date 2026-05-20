package types

import "errors"

var (
	ErrInternalError  = errors.New("internal error")
	ErrParamsLost     = errors.New("parameter lost")
	ErrInvalidNonce   = errors.New("invalid nonce")
	ErrInvalidAddress = errors.New("invalid address")
	ErrInvalidToken   = errors.New("invalid token")
	ErrUnauthorized   = errors.New("unauthorized")
	ErrUserNotFound   = errors.New("user not found")
)
