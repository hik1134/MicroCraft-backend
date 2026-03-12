package mysql

import "MicroCraft/internal/model"

func CreatePost(p *model.Post) error {
	return DB.Create(p).Error
}