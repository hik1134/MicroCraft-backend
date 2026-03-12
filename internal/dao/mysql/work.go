package mysql

import (
	"gorm.io/gorm"
	"MicroCraft/internal/model"
)

func CreateWork(w *model.Work) error {
	return DB.Create(w).Error
}

func UpdateWorkFile(wid uint, fileKey, fileURL, thumbKey, thumbURL string) error {
	return DB.Model(&model.Work{}).Where("id = ?", wid).Updates(map[string]interface{}{
		"file_key":       fileKey,
		"file_url":       fileURL,
		"thumbnail_key":  thumbKey,
		"thumbnail_url":  thumbURL,
	}).Error
}

type WorkListOpt struct {
	Type   string
	Source string
	Status string
	Offset int
	Limit  int
}

func ListWorksByUser(userID uint, opt WorkListOpt) (list []model.Work, total int64, err error) {
	if DB == nil {
		return nil, 0, gorm.ErrInvalidDB
	}
	q := DB.Model(&model.Work{}).Where("user_id = ?", userID)
	if opt.Type != "" {
		q = q.Where("type = ?", opt.Type)
	}
	if opt.Source != "" {
		q = q.Where("source = ?", opt.Source)
	}
	if opt.Status != "" {
		q = q.Where("status = ?", opt.Status)
	}
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := q.Order("created_at DESC").
		Offset(opt.Offset).
		Limit(opt.Limit).
		Find(&list).Error; err != nil {
		return nil, 0, err
	}

	return list, total, nil
}

func GetWorkByID(id uint) (*model.Work, error) {
	var w model.Work
	if err := DB.First(&w, id).Error; err != nil {
		return nil, err
	}
	return &w, nil
}

func SoftDeleteWork(id uint) error {
	return DB.Model(&model.Work{}).
		Where("id = ?", id).
		Update("status", "deleted").Error
}