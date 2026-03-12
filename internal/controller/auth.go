package controller

import (
	"strings"
	"MicroCraft/internal/middleware"
	"MicroCraft/internal/service"
	perr "MicroCraft/pkg/errors"
	"MicroCraft/pkg/response"

	"github.com/gin-gonic/gin"
)

type SendCodeReq struct {
	Email string `json:"email"`
}

func SendEmailCodeHandler(c *gin.Context) {
	var req SendCodeReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailByErr(c, perr.New(perr.INVALID_PARAM))
		return
	}
	email := strings.TrimSpace(req.Email)
	if email == "" || !strings.Contains(email, "@") {
		response.FailByErr(c, perr.New(perr.EMAIL_INVALID))
		return
	}
	expire, err := service.SendEmailCode(c.Request.Context(), email)
	if err != nil {
		response.FailByErr(c, err)
		return
	}
	response.OK(c, gin.H{
		"expire_seconds": expire,
	})
}

func RegisterHandler(c *gin.Context) {
	var req service.RegisterReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailByErr(c, perr.New(perr.INVALID_PARAM))
		return
	}
	resp, err := service.Register(c.Request.Context(), req)
	if err != nil {
		response.FailByErr(c, err)
		return
	}
	response.OK(c, resp)
}

func LoginHandler(c *gin.Context) {
	var req service.LoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailByErr(c, perr.New(perr.INVALID_PARAM))
		return
	}
	resp, err := service.Login(c.Request.Context(), req)
	if err != nil {
		response.FailByErr(c, err)
		return
	}
	response.OK(c, resp)
}

func MeHandler(c *gin.Context) {
	v, ok := c.Get(middleware.CtxUserIDKey)
	if !ok {
		response.FailByErr(c, perr.New(perr.AUTH_TOKEN_MISSING))
		return
	}
	userID, ok := v.(uint)
	if !ok || userID == 0 {
		response.FailByErr(c, perr.New(perr.AUTH_TOKEN_INVALID))
		return
	}
	resp, err := service.GetMe(userID)
	if err != nil {
		response.FailByErr(c, err)
		return
	}
	response.OK(c, resp)
}