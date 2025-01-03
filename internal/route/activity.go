package route

import (
	"git-project-management/internal/project"
	"git-project-management/internal/task"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/go-pg/pg/v10"
)

func Setup(api *huma.API, db *pg.DB) {

	controller := NewController(NewRepo(db), task.NewRepo(db), project.NewRepo(db))

	huma.Register(*api, huma.Operation{
		OperationID: "get-one-activity",
		Method:      http.MethodGet,
		Path:        "/activities/{id}",
		Summary:     "one activity",
		Description: "",
		Tags:        []string{"Activity"},
	}, controller.getOne)

	huma.Register(*api, huma.Operation{
		OperationID: "add-activity",
		Method:      http.MethodPost,
		Path:        "/activities",
		Summary:     "add activity",
		Description: "",
		Tags:        []string{"Activity"},
	}, controller.create)

	huma.Register(*api, huma.Operation{
		OperationID: "get-all-task-activities",
		Method:      http.MethodGet,
		Path:        "/activities",
		Summary:     "get all task activities",
		Description: "",
		Tags:        []string{"Task"},
	}, controller.getAll)
}
