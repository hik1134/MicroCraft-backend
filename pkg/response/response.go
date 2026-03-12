package response

import (
	"log"
	"net/http"

	perr "MicroCraft/pkg/errors"
	"github.com/gin-gonic/gin"
)

type Resp struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func OK(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Resp{
		Code:    string(perr.OK),
		Message: "success",
		Data:    data,
	})
}

func Fail(c *gin.Context, httpStatus int, code string, msg string) {
	c.JSON(httpStatus, Resp{
		Code:    code,
		Message: msg,
		Data:    gin.H{},
	})
}

func FailByErr(c *gin.Context, err error) {
	code := perr.GetCode(err)

	if code == perr.INTERNAL_ERROR && err != nil {
		log.Printf("internal error: %+v", err)
	}

	meta := perr.GetMeta(code)
	Fail(c, meta.HTTPStatus, string(code), meta.Message)
}