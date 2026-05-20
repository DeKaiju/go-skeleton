package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/dekaiju/go-skeleton/pkg/response"
	"github.com/dekaiju/go-skeleton/service/server/logic/user"
	"github.com/dekaiju/go-skeleton/types"
)

func GetNonce(c *gin.Context) {
	req := new(user.NonceReq)
	if err := c.ShouldBindQuery(req); err != nil {
		response.Fail(c, http.StatusBadRequest, types.ErrParamsLost)
		return
	}

	res, err := user.GetNonce(c, req)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, err)
		return
	}

	response.Success(c, res)
}

func Login(c *gin.Context) {
	req := new(user.LoginReq)
	if err := c.ShouldBindJSON(req); err != nil {
		response.Fail(c, http.StatusBadRequest, types.ErrParamsLost)
		return
	}

	res, err := user.Login(c, req)
	if err != nil {
		if errors.Is(err, types.ErrInvalidNonce) || errors.Is(err, types.ErrUnauthorized) {
			response.Fail(c, http.StatusUnauthorized, err)
			return
		}

		response.Fail(c, http.StatusInternalServerError, err)
		return
	}

	response.Success(c, res)
}

func GetProfile(c *gin.Context) {
	res, err := user.GetProfile(c)
	if err != nil {
		switch {
		case errors.Is(err, types.ErrUnauthorized):
			response.Fail(c, http.StatusUnauthorized, err)
		case errors.Is(err, types.ErrUserNotFound):
			response.Fail(c, http.StatusNotFound, err)
		default:
			response.Fail(c, http.StatusInternalServerError, err)
		}
		return
	}

	response.Success(c, res)
}
