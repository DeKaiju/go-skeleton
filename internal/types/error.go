package types

import "errors"

var (
	ErrInternalError         = errors.New("internal error")
	ErrParamsLost            = errors.New("parameter lost")
	ErrInvalidNonce          = errors.New("invalid nonce")
	ErrInvalidAddress        = errors.New("invalid address")
	ErrInvalidRequest        = errors.New("invalid request")
	ErrInvalidParameter      = errors.New("invalid parameter")
	ErrInsufficientVoteQuota = errors.New("insufficient vote quota")
	ErrInvalidRequestToken   = errors.New("invalid request token")
)

var ErrorCodeMap = map[error]int{
	ErrInternalError:         5000,
	ErrParamsLost:            5001,
	ErrInvalidNonce:          5002,
	ErrInvalidAddress:        5003,
	ErrInvalidRequest:        5004,
	ErrInvalidParameter:      5005,
	ErrInsufficientVoteQuota: 5006,
	ErrInvalidRequestToken:   5017,
}
