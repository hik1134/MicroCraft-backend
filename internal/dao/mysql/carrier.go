package mysql

import (
	"MicroCraft/internal/model"
)

func ListCarriers(status string) ([]model.Carrier, error) {
	var list []model.Carrier
	q := DB.Model(&model.Carrier{})
	if status != "" {
		q = q.Where("status = ?", status)
	}
	if err := q.Order("created_at DESC").Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}