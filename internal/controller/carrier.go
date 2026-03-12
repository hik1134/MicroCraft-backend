package controller

import (
	"MicroCraft/internal/service"
	perr "MicroCraft/pkg/errors"
	"MicroCraft/pkg/response"
	"github.com/gin-gonic/gin"
)

func GetCarriersHandler(c *gin.Context) {
	var req service.CarrierListReq
	if err := c.ShouldBindQuery(&req); err != nil {
		response.FailByErr(c, perr.New(perr.INVALID_PARAM))
		return
	}
	resp, err := service.GetCarriers(c.Request.Context(), req)
	if err != nil {
		response.FailByErr(c, err)
		return
	}
	response.OK(c, resp)
}