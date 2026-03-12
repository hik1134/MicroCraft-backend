package mysql

import (
	"MicroCraft/internal/model"
	"errors"

	"gorm.io/gorm"
)

func GetUserByEmail(email string) (*model.User, error) {
	var u model.User
	err := DB.Where("email = ?", email).First(&u).Error
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func GetUserByID(id uint) (*model.User, error) {
	var u model.User
	if err := DB.First(&u, id).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

func CreateUser(u *model.User) error {
	return DB.Create(u).Error
}

func IsNotFound(err error) bool {
	return errors.Is(err, gorm.ErrRecordNotFound)
}