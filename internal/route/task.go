package route

import (
	"git-project-management/internal/controller"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

func SetupTask(api huma.API) {

	ctrl := controller.NewTaskController()

	huma.Register(api, huma.Operation{
		OperationID: "get-one-task",
		Method:      http.MethodGet,
		Path:        "/tasks/{id}",
		Summary:     "get task by id",
		Description: "",
		Tags:        []string{"Task"},
	}, ctrl.GetTask)

	// huma.Register(*api, huma.Operation{
	// 	OperationID: "set-task-status",
	// 	Method:      http.MethodPut,
	// 	Path:        "/tasks/{id}/set-status",
	// 	Summary:     "set task status",
	// 	Description: "",
	// 	Tags:        []string{"Task"},
	// }, controller.setStatus)

}
