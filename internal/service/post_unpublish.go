package service

import (
	"context"
	"errors"

	"MicroCraft/internal/dao/mysql"
	perr "MicroCraft/pkg/errors"

	"gorm.io/gorm"
)

func UnpublishPost(ctx context.Context, userID uint, postID uint) error {
	if userID == 0 || postID == 0 {
		return perr.New(perr.INVALID_PARAM)
	}

	p, err := mysql.GetPostByID(postID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return perr.New(perr.POST_NOT_FOUND)
		}
		return perr.Wrap(perr.DB_QUERY_FAIL, err)
	}

	w, err := mysql.GetWorkByID(p.WorkID)
	if err != nil {
		if mysql.IsNotFound(err) {
			return perr.New(perr.POST_NOT_FOUND)
		}
		return perr.Wrap(perr.DB_QUERY_FAIL, err)
	}

	if w.UserID != userID {
		return perr.New(perr.POST_FORBIDDEN)
	}

	if p.Status == "draft" {
		return nil
	}

	if err := mysql.UpdatePostStatus(p.ID, "draft"); err != nil {
		return perr.Wrap(perr.DB_UPDATE_FAIL, err)
	}

	return nil
}