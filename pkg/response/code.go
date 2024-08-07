package response

import "github.com/dekaiju/go-skeleton/internal/types"

func getCodeByError(err error) int {
	if code, ok := types.ErrorCodeMap[err]; ok {
		return code
	}
	return 5000
}
