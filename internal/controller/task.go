package controller

import (
	"context"
	"errors"

	"gitea.com/logicamp/lc"
	"github.com/go-pg/pg/v10"
)

type TaskController struct {
}

func NewTaskController() *TaskController {
	return &TaskController{}
}

func (c *Controller) getOne(_ context.Context, req *GetOneRequest) (*GetOneResponse, error) {

	task, err := c.repo.GetByID(req.Id)
	if err != nil {
		return nil, lc.SendInternalErrorResponse(err, "[activity] get one")
	}

	if errors.Is(err, pg.ErrNoRows) {
		return nil, errors.New("project not found")
	}

	return &GetOneResponse{
		Body: ToTaskDTO(*task),
	}, nil
}

func (c *Controller) setStatus(_ context.Context, req *SetTaskStatusRequest) (*SetTaskStatusResponse, error) {

	task, err := c.repo.GetByID(req.TaskID)
	if err != nil {
		if errors.Is(err, pg.ErrNoRows) {
			return nil, errors.New("project not found")
		}
		return nil, lc.SendInternalErrorResponse(err, "[task] update task")
	}

	task.Status = req.Body.Status
	updatedTask, err := c.repo.Update(req.TaskID, task)
	if err != nil {
		return nil, lc.SendInternalErrorResponse(err, "[task] update task")
	}
	return &SetTaskStatusResponse{
		Body: ToTaskDTO(*updatedTask),
	}, nil
}

func ToTaskDTO(model TaskEntity) TaskDTO {
	return TaskDTO{
		ID:          model.ID,
		ParentID:    model.ParentID,
		Title:       model.Title,
		Description: model.Description,
		Status:      model.Status,
		Priority:    model.Priority,
		AssigneeID:  model.AssigneeID,
		ProjectID:   model.ProjectID,
		DueDate:     model.DueDate,
		CreatedBy:   model.CreatedBy,
		CreatedAt:   model.CreatedAt,
	}
}
