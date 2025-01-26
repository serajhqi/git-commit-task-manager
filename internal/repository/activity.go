package repository

import (
	"git-project-management/internal/database"
	"git-project-management/internal/types"
)

func GetActivities(taskID int64, userID int64) (*[]types.ActivityEntity, *int, error) {
	var activityEntities []types.ActivityEntity

	// query := `
	// 	SELECT
	// 	*
	// 	FROM task AS t
	// 	LEFT JOIN activities AS a ON t.id = a.task_id
	// 	LEFT JOIN user AS u ON a.created_by = u.id
	// 	WHERE t.id = ?

	// `
	count, err := database.GetDB().Model(&activityEntities).Where("task_id = ? AND user_id = ", taskID, userID).Order("id DESC").SelectAndCount()
	if err != nil {
		return nil, nil, err
	}
	return &activityEntities, &count, err
}
