package service

import (
	"time"
	"MicroCraft/internal/dao/mysql"
	perr "MicroCraft/pkg/errors"
)

type MeResp struct {
	UserID   uint   `json:"user_id"`
	Email    string `json:"email"`
	Nickname string `json:"nickname"`
	CreatedAt time.Time `json:"created_at"`
}

func GetMe(userID uint) (*MeResp, error) {
	u, err := mysql.GetUserByID(userID)
	if err != nil {
		if mysql.IsNotFound(err) {
			return nil, perr.New(perr.USER_NOT_FOUND)
		}
		return nil, perr.Wrap(perr.DB_QUERY_FAIL, err)
	}

	return &MeResp{
		UserID:   u.ID,
		Email:    u.Email,
		Nickname: u.Nickname,
		CreatedAt: u.CreatedAt,
	}, nil
}