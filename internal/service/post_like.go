package service

import (
	"context"
	"errors"

	"MicroCraft/internal/dao/mysql"
	perr "MicroCraft/pkg/errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ToggleLikeResp struct {
	PostID    uint  `json:"post_id"`
	Liked     bool  `json:"liked"`      
	LikeCount int64 `json:"like_count"` 
}

func ToggleLikePost(ctx context.Context, userID uint, postID uint) (*ToggleLikeResp, error) {
	if userID == 0 || postID == 0 {
		return nil, perr.New(perr.INVALID_PARAM)
	}

	var out *ToggleLikeResp

	err := mysql.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		p, err := mysql.GetPostByIDForUpdate(tx, postID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return perr.New(perr.POST_NOT_FOUND)
			}
			return perr.Wrap(perr.DB_QUERY_FAIL, err)
		}

		if p.Status != "published" {
			return perr.New(perr.POST_FORBIDDEN)
		}

		liked, err := mysql.IsPostLikedTx(tx, postID, userID)
		if err != nil {
			return perr.Wrap(perr.DB_QUERY_FAIL, err)
		}

		if !liked {
			if err := mysql.CreatePostLikeTx(tx, postID, userID); err != nil {
				return perr.Wrap(perr.DB_CREATE_FAIL, err)
			}
			if err := mysql.IncPostLikeCountTx(tx, postID, 1); err != nil {
				return perr.Wrap(perr.DB_UPDATE_FAIL, err)
			}
			p.LikeCount++
			out = &ToggleLikeResp{PostID: postID, Liked: true, LikeCount: p.LikeCount}
			return nil
		}

		if err := mysql.DeletePostLikeTx(tx, postID, userID); err != nil {
			return perr.Wrap(perr.DB_UPDATE_FAIL, err)
		}
		if p.LikeCount > 0 {
			if err := mysql.IncPostLikeCountTx(tx, postID, -1); err != nil {
				return perr.Wrap(perr.DB_UPDATE_FAIL, err)
			}
			p.LikeCount--
		}
		out = &ToggleLikeResp{PostID: postID, Liked: false, LikeCount: p.LikeCount}
		return nil
	})

	if err != nil {
		return nil, err
	}
	return out, nil
}

var _ = clause.Locking{}