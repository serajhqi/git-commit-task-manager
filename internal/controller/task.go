package controller

import (
	"context"
	"git-project-management/internal/repository"
	"git-project-management/internal/types"

	"gitea.com/logicamp/lc"
)

type TaskController struct{}

func NewTaskController() TaskController {
	return TaskController{}
}

func AddTask(ctx context.Context, req *types.CreateTaskRequest) (*lc.RespWithBody[types.TaskDTO], error) {

	var createdBy int64 = 1

	// check if it has access to this project or parent task id

	taskEntity, err := repository.Create(ctx, types.TaskEntity{
		ParentID:    req.Body.ParentID,
		Title:       req.Body.Title,
		Description: req.Body.Description,
		Status:      req.Body.Status,
		Priority:    req.Body.Priority,
		AssigneeID:  req.Body.AssigneeID,
		ProjectID:   req.ProjectId,
		DueDate:     req.Body.DueDate,
		CreatedBy:   createdBy,
	})

	if err != nil {
		return nil, repository.HandleError(err)
	}

	dto := taskEntityToDTO(*taskEntity)
	return &lc.RespWithBody[types.TaskDTO]{
		Body: &dto,
	}, nil
}

// ---------------------

func taskEntityToDTO(entity types.TaskEntity) types.TaskDTO {
	return types.TaskDTO{
		ID:          entity.ID,
		ParentID:    entity.ParentID,
		Title:       entity.Title,
		Description: entity.Description,
		Status:      entity.Status,
		Priority:    entity.Priority,
		AssigneeID:  entity.AssigneeID,
		ProjectID:   entity.ProjectID,
		DueDate:     entity.DueDate,
		CreatedBy:   entity.CreatedBy,
		CreatedAt:   entity.CreatedAt,
		UpdatedAt:   entity.UpdatedAt,
	}
}
