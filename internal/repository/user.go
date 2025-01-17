package repository

import (
	"git-project-management/internal/database"
	"git-project-management/internal/types"
)

func GetUser(email, password string) (*types.UserEntity, error) {
	var user types.UserEntity
	err := database.GetDB().Model(&user).Where("email = ? AND password = ?", email, password).First()
	return &user, err
}

func GetUserByEmail(email string) (*types.UserEntity, error) {
	var user types.UserEntity
	err := database.GetDB().Model(&user).Where("email = ?", email).First()
	return &user, err
}
