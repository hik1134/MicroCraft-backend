package service

import (
	"context"
	"strings"

	"MicroCraft/internal/dao/mysql"
	perr "MicroCraft/pkg/errors"
)

type CarrierListReq struct {
	Status string `form:"status"` // 默认 active
}

type CarrierListResp struct {
	List interface{} `json:"list"`
}

func GetCarriers(ctx context.Context, req CarrierListReq) (*CarrierListResp, error) {
	status := strings.TrimSpace(req.Status)
	if status == "" {
		status = "active"
	}

	if status != "active" && status != "disabled" {
		return nil, perr.New(perr.INVALID_PARAM)
	}

	list, err := mysql.ListCarriers(status)
	if err != nil {
		return nil, perr.Wrap(perr.DB_QUERY_FAIL, err)
	}

	return &CarrierListResp{
		List: list,
	}, nil
}