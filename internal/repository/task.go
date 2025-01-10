package repository

import (
	"fmt"
	"git-project-management/internal/database"
	"git-project-management/internal/types"
)

func GetUserTask(taskId int64, userId int64) (*types.TaskEntity, error) {
	var task types.TaskEntity
	// todo check if the task belongs to this user
	err := database.GetDB().Model(&task).Where("id = ? AND created_by = ?", taskId, userId).Select()
	if err != nil {
		return nil, fmt.Errorf("failed to read record by ID: %w", err)
	}
	return &task, nil
}
