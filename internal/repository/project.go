package repository

import (
	"fmt"
	"git-project-management/internal/database"
	"git-project-management/internal/types"
)

func GetUserProject(projectId int64, userId int64) (*types.ProjectEntity, error) {
	var project types.ProjectEntity
	// todo check if the project is owned or assinged to this user
	err := database.GetDB().Model(&project).Where("id = ? AND created_by = ?", projectId, userId).Select()
	if err != nil {
		return nil, fmt.Errorf("failed to read record by ID: %w", err)
	}
	return &project, nil
}
