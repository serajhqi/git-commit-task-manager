package repository

import (
	"fmt"
	"git-project-management/internal/database"
	"git-project-management/internal/types"
)

func GetUserTask(taskId int64, userId int64) (*types.TaskEntity, error) {
	var task types.TaskEntity
	err := database.GetDB().Model(&task).Where("id = ? AND created_by = ?", taskId, userId).First()
	if err != nil {
		return nil, fmt.Errorf("failed to read record by ID: %w", err)
	}
	return &task, nil
}

func GetUserTasks(userId int64, offset, limit int) ([]types.TaskEntity, int, error) {
	var task []types.TaskEntity
	count, err := database.GetDB().Model(&task).
		Where("created_by = ?", userId).
		Offset(offset).
		Limit(limit).
		Order("id DESC").
		SelectAndCount()
	if err != nil {
		return nil, count, fmt.Errorf("failed to read record by ID: %w", err)
	}
	return task, count, nil
}

func GetProjectTasks(projectId int64, userId int64, offset, limit int) ([]types.TaskEntity, int, error) {
	var task []types.TaskEntity
	count, err := database.GetDB().Model(&task).
		Where("project_id = ? AND created_by = ?", projectId, userId).
		Offset(int(offset)).
		Limit(int(limit)).
		Order("id DESC").
		SelectAndCount()
	if err != nil {
		return nil, count, fmt.Errorf("failed to read record by ID: %w", err)
	}
	return task, count, nil
}
