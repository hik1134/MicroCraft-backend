package mysql

import (
	"MicroCraft/internal/model"
	"errors"
	"gorm.io/gorm"
)

func IsPostLiked(postID, userID uint) (bool, error) {
	var pl model.PostLike
	err := DB.Where("post_id = ? AND user_id = ?", postID, userID).First(&pl).Error
	if err == nil {
		return true, nil
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	}
	return false, err
}
func IsPostLikedTx(tx *gorm.DB, postID, userID uint) (bool, error) {
	var pl model.PostLike
	err := tx.Where("post_id = ? AND user_id = ?", postID, userID).First(&pl).Error
	if err == nil {
		return true, nil
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	}
	return false, err
}
func CreatePostLikeTx(tx *gorm.DB, postID, userID uint) error {
	pl := &model.PostLike{
		PostID: postID,
		UserID: userID,
	}
	return tx.Create(pl).Error
}

func DeletePostLikeTx(tx *gorm.DB, postID, userID uint) error {
	return tx.Where("post_id = ? AND user_id = ?", postID, userID).
		Delete(&model.PostLike{}).Error
}