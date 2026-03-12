package middleware

import (
	"errors"
	"strings"

	"MicroCraft/internal/config"
	perr "MicroCraft/pkg/errors"
	"MicroCraft/pkg/response"
	"MicroCraft/pkg/utils"

	"github.com/gin-gonic/gin"
)

const CtxUserIDKey = "user_id"

func AuthJWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := strings.TrimSpace(c.GetHeader("Authorization"))
		if auth == "" {
			response.FailByErr(c, perr.New(perr.AUTH_TOKEN_MISSING))
			c.Abort()
			return
		}
		tokenStr := auth
		if strings.HasPrefix(strings.ToLower(auth), "bearer ") {
			tokenStr = strings.TrimSpace(auth[7:])
		}
		if tokenStr == "" {
			response.FailByErr(c, perr.New(perr.AUTH_TOKEN_MISSING))
			c.Abort()
			return
		}
		if config.Conf == nil || config.Conf.Jwt.Secret == "" {
			response.FailByErr(c, perr.New(perr.JWT_SECRET_EMPTY))
			c.Abort()
			return
		}
		userID, err := utils.ParseToken(config.Conf.Jwt.Secret, tokenStr)
		if err != nil {
			if errors.Is(err, utils.ErrTokenExpired) {
				response.FailByErr(c, perr.New(perr.AUTH_TOKEN_EXPIRED))
			} else {
				response.FailByErr(c, perr.New(perr.AUTH_TOKEN_INVALID))
			}
			c.Abort()
			return
		}
		c.Set(CtxUserIDKey, userID)
		c.Next()
	}
}