package repository

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"git-project-management/internal/database"
	"git-project-management/internal/types"
)

func RemoveApiKey(key string) error {
	var entity types.ApiKeyEntity
	_, err := database.GetDB().Model(&entity).Where("key = ?", key).Delete()
	if err != nil {
		return fmt.Errorf("failed to delete record: %w", err)
	}
	return nil
}

func GetApiKey(key string) (*types.ApiKeyEntity, error) {
	hash := md5.Sum([]byte(key))
	hashedKey := hex.EncodeToString(hash[:])

	var entity types.ApiKeyEntity
	err := database.GetDB().Model(&entity).Where("key = ?", hashedKey).First()
	if err != nil {
		return nil, fmt.Errorf("failed to delete record: %w", err)
	}
	return &entity, nil
}
