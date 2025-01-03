package controller

import (
	"context"
	"errors"
	"git-project-management/internal/task"

	"gitea.com/logicamp/lc"
	"github.com/go-pg/pg/v10"
)

type ProjectController struct {
}

func NewProjectController(repo *Repo) *ProjectController {
	return &ProjectController{}
}

func (c *Controller) getOne(_ context.Context, req *GetOneRequest) (*GetOneResponse, error) {

	project, err := c.repo.GetByID(req.Id)
	if err != nil {
		return nil, lc.SendInternalErrorResponse(err, "[activity] get all")
	}

	return &GetOneResponse{
		Body: ToProjectDTO(*project),
	}, nil
}

func (c *Controller) getAll(_ context.Context, req *GetAllRequest) (*GetAllResponse, error) {
	limit := 100
	offset := 0
	if req.Limit > 0 {
		limit = req.Limit
	}

	if req.Offset > 0 {
		offset = req.Offset
	}

	projects, err := c.repo.getAll(limit, offset)
	if err != nil {
		return nil, lc.SendInternalErrorResponse(err, "[activity] get all")
	}
	var projectsDTO []ProjectDTO
	for _, v := range projects {
		projectsDTO = append(projectsDTO, ToProjectDTO(v))
	}

	return &GetAllResponse{
		Body: projectsDTO,
	}, nil
}

func (c *Controller) getAllTasks(_ context.Context, req *GetAllTasksRequest) (*GetAllTasksResponse, error) {
	limit := 100
	offset := 0
	if req.Limit > 0 {
		limit = req.Limit
	}

	if req.Offset > 0 {
		offset = req.Offset
	}

	tasks, err := c.repo.getAllTasks(req.ProjectId, limit, offset)
	if err != nil {
		return nil, lc.SendInternalErrorResponse(err, "[activity] get all")
	}

	var tasksDTO []task.TaskDTO
	for _, v := range tasks {
		tasksDTO = append(tasksDTO, task.ToTaskDTO(v))
	}

	return &GetAllTasksResponse{
		Body: tasksDTO,
	}, nil
}

func (c *Controller) addTask(_ context.Context, req *CreateTaskRequest) (*CreateTaskResponse, error) {

	// TODO project_id validation
	// TODO parent_id validation
	// TODO assignee_id validation

	status := task.IN_PROGRESS
	if req.Body.Status != "" {
		status = req.Body.Status
	}
	priority := task.MEDIUM
	if req.Body.Priority != "" {
		priority = req.Body.Priority
	}

	taskEntity := &task.TaskEntity{
		Title:       req.Body.Title,
		Description: req.Body.Description,
		Status:      status,
		Priority:    priority,
		AssigneeID:  1,
		ProjectID:   req.ProjectId,
		DueDate:     req.Body.DueDate,
		CreatedBy:   1,
	}
	err := c.repo.addTask(taskEntity)
	if err != nil {
		return nil, lc.SendInternalErrorResponse(err, "[project] create task")
	}

	if errors.Is(err, pg.ErrNoRows) {
		return nil, errors.New("project not found")
	}
	return &CreateTaskResponse{Body: task.ToTaskDTO(*taskEntity)}, nil
}

func ToProjectDTO(model ProjectEntity) ProjectDTO {
	return ProjectDTO{
		ID:          model.ID,
		Name:        model.Name,
		Description: model.Description,
		StartDate:   model.StartDate,
		EndDate:     model.EndDate,
		CreatedBy:   model.CreatedBy,
		CreatedAt:   model.CreatedAt,
	}
}
