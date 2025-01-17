package controller

import (
	"context"
	"encoding/base64"
	"fmt"
	"git-project-management/internal/controller/utils"
	"git-project-management/internal/repository"
	"git-project-management/internal/types"
	"time"

	"gitea.com/logicamp/lc"
	"github.com/danielgtaylor/huma/v2"
	"go.uber.org/zap"
)

type CommitController struct{}

func NewCommitController() CommitController {
	return CommitController{}
}

func (cc CommitController) Commit(ctx context.Context, req *types.CommitRequest) (*struct {
	Message string `json:"message"`
}, error) {
	response := "OK"

	decodedData, err := base64.StdEncoding.DecodeString(req.Body.CommitMessage)
	if err != nil {
		fmt.Println("Error decoding string:", err)
		return nil, huma.Error400BadRequest("")
	}

	// Convert the decoded byte slice to a string
	decodedString := string(decodedData)
	commit, customPattern := utils.ParseCommitMessage(decodedString)
	commit.Hash = req.Body.CommitHash
	commit.Branch = req.Body.Branch
	lc.Logger.Debug("parsed message", zap.Any("customPattern", customPattern), zap.Any("commit", commit))

	// Validate the last line of the commit message
	validationResult := utils.IsValidCommitLine(customPattern)
	if validationResult.ProjectID == 0 {
		validationResult.ProjectID = req.Body.ProjectID
	}
	if validationResult.IsValid {
		fmt.Println("✅ Valid")
	} else {
		fmt.Println(validationResult)
		fmt.Printf("Invalid commit message: %s %d\n", validationResult.Message, validationResult.ProjectID)
		fmt.Println("❌ Invalid")
	}

	userID := utils.GetCtxUserID(ctx)
	// run the action ------------------------------------------------------
	if !validationResult.IsValid {
		return nil, huma.Error400BadRequest(validationResult.Message)
	}

	// check project
	_, err = repository.GetUserProject(validationResult.ProjectID, userID)
	if err != nil {
		lc.Logger.Error("[commit] get project error", zap.Error(err))
		return nil, repository.HandleError(err)
	}

	// check existing task
	var task *types.TaskEntity
	if !validationResult.IsNewTask {
		task, err = repository.GetUserTask(validationResult.TaskID, userID)
		if err != nil {
			lc.Logger.Error("[commit] get task error", zap.Error(err))
			return nil, repository.HandleError(err)
		}

		_, err = repository.Update(context.Background(), task.ID, task)
		if err != nil {
			lc.Logger.Error("[commit] update task error", zap.Error(err))
			return nil, repository.HandleError(err)
		}
	} else {

		task = &types.TaskEntity{
			Title:      validationResult.TaskTitle,
			Status:     types.DEFAULT_TASK_STATUS,
			Priority:   types.DEFAULT_TASK_PRIORITY,
			AssigneeID: userID,
			Weight:     types.DEFAULT_TASK_WEIGHT,
			ProjectID:  validationResult.ProjectID,
			CreatedBy:  userID,
			UpdatedAt:  time.Now(),
		}
		task, err = repository.Create(context.Background(), *task)
		if err != nil {
			lc.Logger.Error("[commit] failed to create new task", zap.Error(err))
			return nil, repository.HandleError(err)
		}
	}

	var timelog int
	if validationResult.HasTimelog {
		timelog, _ = utils.ConvertTimelogToMinutes(validationResult.Timelog)
	}

	taskActivity := types.ActivityEntity{
		TaskID:      task.ID,
		CommitHash:  commit.Hash,
		Branch:      commit.Branch,
		Title:       commit.Title,
		Description: commit.Description,
		Duration:    timelog,
		CreatedBy:   userID,
	}

	_, err = repository.Create(context.Background(), taskActivity)
	if err != nil {
		return nil, repository.HandleError(err)
	}
	// ---------------------------------------------------------------------

	return &struct {
		Message string `json:"message"`
	}{Message: response}, nil
}

func (cc CommitController) ImportHistory(ctx context.Context, req *types.ImportHistory) (*struct {
	Message string `json:"message"`
}, error) {

	return nil, nil
}
