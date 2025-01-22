package repository

import (
	"context"
	"fmt"
	"git-project-management/internal/database"
)

func Create[ENTITY any](ctx context.Context, entity ENTITY) (*ENTITY, error) {
	_, err := database.GetDB().Model(&entity).Context(ctx).Insert()
	if err != nil {
		return nil, fmt.Errorf("failed to create record: %w", err)
	}
	return &entity, nil
}

func ReadByID[ENTITY any](ctx context.Context, id int64) (*ENTITY, error) {
	var entity ENTITY
	err := database.GetDB().Model(&entity).Context(ctx).Where("id = ?", id).Select()
	if err != nil {
		return nil, fmt.Errorf("failed to read record by ID: %w", err)
	}
	return &entity, nil
}

func Update[ENTITY any](ctx context.Context, id int64, entity ENTITY) (*ENTITY, error) {
	_, err := database.GetDB().Model(&entity).Context(ctx).Where("id = ?", id).Update()
	if err != nil {
		return nil, fmt.Errorf("failed to update record: %w", err)
	}
	return &entity, nil
}

func Delete[ENTITY any](ctx context.Context, id int64) error {
	var entity ENTITY
	_, err := database.GetDB().Model(&entity).Context(ctx).Where("id = ?", id).Delete()
	if err != nil {
		return fmt.Errorf("failed to delete record: %w", err)
	}
	return nil
}

func ReadAll[ENTITY any](offset, limit int) ([]ENTITY, int, error) {
	var entities []ENTITY
	count, err := database.GetDB().Model(&entities).Offset(offset).Limit(limit).SelectAndCount()
	if err != nil {
		return nil, count, fmt.Errorf("failed to read all records: %w", err)
	}
	return entities, count, nil
}
