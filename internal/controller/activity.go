package controller

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"git-project-management/internal/project"
	"git-project-management/internal/task"
	"strconv"
	"time"

	"gitea.com/logicamp/lc"
	"github.com/danielgtaylor/huma/v2"
	"github.com/go-pg/pg/v10"
)

type Controller struct {
	repo        *Repo
	taskRepo    *task.Repo
	projectRepo *project.Repo
}

func NewController(repo *Repo, taskRepo *task.Repo, projectRepo *project.Repo) *Controller {
	return &Controller{
		repo:        repo,
		taskRepo:    taskRepo,
		projectRepo: projectRepo,
	}
}

func (c *Controller) getOne(_ context.Context, req *GetOneRequest) (*GetOneResponse, error) {

	activity, err := c.repo.getByID(req.Id)
	if err != nil {
		return nil, lc.SendInternalErrorResponse(err, "[activity] get all")
	}

	return &GetOneResponse{
		Body: ToActivityDTO(*activity),
	}, nil
}

// ----------------------------------------------------

func (c *Controller) create(_ context.Context, req *CreateActivityRequest) (*CreateActivityResponse, error) {

	data, _ := base64.StdEncoding.DecodeString(req.Body.CommitMessage)

	parsedValues := ParseFullCommitMessage(string(data))
	commitMessage := parsedValues["Commit Message"]
	author := parsedValues["Author"]
	hash := parsedValues["Hash"]
	branch := parsedValues["Branch"]
	date := parsedValues["Date"]
	isoDate := parsedValues["ISO Date"]

	fmt.Println("\nExtracted Fields:")
	fmt.Println("Author:", author)
	fmt.Println("Hash:", hash)
	fmt.Println("Branch:", branch)
	fmt.Println("Date:", date)
	fmt.Println("ISO Date:", isoDate)

	parsedData, err := ParseCommitMessage(commitMessage)
	if err != nil {
		fmt.Println("Error:", err)
		return nil, huma.Error400BadRequest(err.Error())
	}

	fmt.Println("Parsed Data:")
	fmt.Println("Activity Title:", parsedData.ActivityTitle)
	fmt.Println("Activity Description:", parsedData.ActivityDescription)
	fmt.Println("Task ID:", parsedData.TaskID)
	fmt.Println("Task Status:", parsedData.TaskStatus)
	fmt.Println("Project ID:", parsedData.ProjectID)
	fmt.Println("New Task Title:", parsedData.NewTaskTitle)
	fmt.Println("New Task Description:", parsedData.NewTaskDescription)
	fmt.Println("New Task Status:", parsedData.NewTaskStatus)
	fmt.Println("Time Spent:", parsedData.TimeSpentMinutes)

	// check whether task exists

	activity := &ActivityEntity{
		CommitHash:  hash,
		Branch:      branch,
		Title:       parsedData.ActivityTitle,
		Description: parsedData.ActivityDescription,
		CreatedBy:   1,
	}

	if len(parsedData.TaskID) > 0 {
		// error already handled in ParseCommitMessage()
		taskId, _ := strconv.Atoi(parsedData.TaskID)
		activity.TaskID = int64(taskId)

		// TODO later we may need to check if the task is owned by current user
		task_, err := c.taskRepo.GetByID(int64(taskId))
		if err != nil {
			if errors.Is(err, pg.ErrNoRows) {
				return nil, huma.Error404NotFound(fmt.Sprintf("task id %d not found", taskId))
			}
			return nil, lc.SendInternalErrorResponse(err, "[activity] get task by id")
		}

		if len(parsedData.TaskStatus) > 0 {
			task_.Status = task.TaskStatus(parsedData.TaskStatus)
			c.taskRepo.Update(task_.ID, task_)
		}

		if parsedData.TimeSpentMinutes > 0 {
			activity.Duration = &parsedData.TimeSpentMinutes
		} else {
			// get the last activity with the same task id and user id
			// get the diff from the last createAt and the current one and take it as task duration
			penultimateActivity, err := c.repo.findByUserID(1)
			if err == nil {
				duration := int(time.Since(penultimateActivity.CreatedAt).Minutes())
				activity.Duration = &duration
			}
		}

	} else if len(parsedData.ProjectID) > 0 {
		projectId, _ := strconv.Atoi(parsedData.ProjectID)
		_, err := c.projectRepo.GetByID(int64(projectId))
		if err != nil {
			if errors.Is(err, pg.ErrNoRows) {
				return nil, huma.Error404NotFound(fmt.Sprintf("project id %d not found", projectId))
			}
			return nil, lc.SendInternalErrorResponse(err, "[activity] get project by id")
		}

		//create the task
		taskEntity := &task.TaskEntity{
			Title:       parsedData.NewTaskTitle,
			Description: parsedData.ActivityDescription,
			Status:      task.TaskStatus(parsedData.NewTaskStatus),
			Priority:    task.MEDIUM,
			ProjectID:   int64(projectId),
			AssigneeID:  1,
			CreatedBy:   1,
		}

		taskEntity_, err := c.taskRepo.Create(taskEntity)
		if err != nil {
			return nil, lc.SendInternalErrorResponse(err, "[activity] create task error")
		}
		activity.TaskID = taskEntity_.ID
	}

	activityEntity, err := c.repo.create(activity)
	if err != nil {
		return nil, lc.SendInternalErrorResponse(err, "[activity] create activity error")
	}

	return &CreateActivityResponse{
		Body: IdBody{
			Id: activityEntity.ID,
		},
	}, nil
}

func (c *Controller) getAll(_ context.Context, req *GetAllActivitiesRequest) (*GetAllActivitiesResponse, error) {

	limit := 100
	offset := 0
	if req.Limit > 0 {
		limit = req.Limit
	}

	if req.Offset > 0 {
		offset = req.Offset
	}

	activities, err := c.repo.getAll(req.TaskID, limit, offset)
	if err != nil {
		return nil, lc.SendInternalErrorResponse(err, "[activity] get all")
	}

	var activitiesDTO []ActivityDTO
	for _, v := range activities {
		activitiesDTO = append(activitiesDTO, ToActivityDTO(v))
	}

	return &GetAllActivitiesResponse{
		Body: activitiesDTO,
	}, nil
}
