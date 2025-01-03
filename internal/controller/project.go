package controller

import (
	"context"
	"git-project-management/internal/repository"
	"git-project-management/internal/types"
	"time"

	"gitea.com/logicamp/lc"
)

type ProjectController struct{}

func NewProjectController() ProjectController {
	return ProjectController{}
}

func AddProject(ctx context.Context, req *types.CreateProjectRequest) (*lc.RespWithBody[types.ProjectDTO], error) {

	project, err := repository.Create(ctx, types.ProjectEntity{
		Name:        req.Name,
		Description: req.Description,
		CreatedBy:   0,
		CreatedAt:   time.Time{},
	})

	if err != nil {
		return nil, repository.HandleError(err)
	}

	return &lc.RespWithBody[types.ProjectDTO]{
		Body: &types.ProjectDTO{
			ID:          project.ID,
			Name:        project.Name,
			Description: project.Description,
			CreatedBy:   0,
			CreatedAt:   time.Time{},
		},
	}, nil
}
