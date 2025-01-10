package route

import (
	"git-project-management/internal/controller"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

func SetupProject(api huma.API) {

	_ = controller.NewProjectController()

	huma.Register(api, huma.Operation{
		OperationID: "add-task",
		Method:      http.MethodPost,
		Path:        "/projects",
		Summary:     "add project",
		Description: "",
		Tags:        []string{"Project"},
	}, controller.AddProject)

	// huma.Register(*api, huma.Operation{
	// 	OperationID: "get-one-project",
	// 	Method:      http.MethodGet,
	// 	Path:        "/projects/{id}",
	// 	Summary:     "one project",
	// 	Description: "",
	// 	Tags:        []string{"Project"},
	// }, controller.getOne)

	huma.Register(api, huma.Operation{
		OperationID: "get-all-projects",
		Method:      http.MethodGet,
		Path:        "/projects",
		Summary:     "all projects",
		Description: "",
		Tags:        []string{"Project"},
	}, controller.GetProjects)

	// huma.Register(*api, huma.Operation{
	// 	OperationID: "get-all-tasks",
	// 	Method:      http.MethodGet,
	// 	Path:        "/projects/{project_id}/tasks",
	// 	Summary:     "all tasks",
	// 	Description: "",
	// 	Tags:        []string{"Project"},
	// }, controller.getAllTasks)

}
