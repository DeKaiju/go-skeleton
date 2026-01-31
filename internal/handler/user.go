package handler

import (
	"github.com/gin-gonic/gin"

	"github.com/dekaiju/go-skeleton/internal/logic/user"
	"github.com/dekaiju/go-skeleton/internal/types"
	"github.com/dekaiju/go-skeleton/pkg/response"
)

func Login(c *gin.Context) {
	req := new(user.LoginReq)

	if err := c.BindJSON(req); err != nil {
		response.Fail(c, types.ErrParamsLost)
		return
	}

	res, err := user.Login(c, req.Message, req.Signature)
	if err != nil {
		response.Fail(c, err)
		return
	}

	response.Success(c, res)
}
