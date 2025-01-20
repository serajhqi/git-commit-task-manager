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

	var createdBy int64 = 1
	project, err := repository.Create(ctx, types.ProjectEntity{
		Name:        req.Body.Name,
		Description: req.Body.Description,
		CreatedBy:   createdBy,
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

func GetProjects(ctx context.Context, req *types.GetAllRequest) (*lc.RespWithBodyPaginated[types.ProjectDTO], error) {

	if req.Limit == 0 {
		req.Limit = 10
	}

	projectEntities, total, err := repository.ReadAll[types.ProjectEntity](req.Offset, req.Limit)

	if err != nil {
		return nil, repository.HandleError(err)
	}

	var projects []types.ProjectDTO
	for _, projectEntity := range projectEntities {
		projects = append(projects, entityToDTO(projectEntity))
	}

	return &lc.RespWithBodyPaginated[types.ProjectDTO]{
		Body: lc.PaginatedBody[types.ProjectDTO]{
			Total: total,
			Items: projects,
		},
	}, nil
}

// -----------

func entityToDTO(input types.ProjectEntity) types.ProjectDTO {
	return types.ProjectDTO{
		ID:          input.ID,
		Name:        input.Name,
		Description: input.Description,
		CreatedBy:   input.CreatedBy,
		CreatedAt:   input.CreatedAt,
	}
}
