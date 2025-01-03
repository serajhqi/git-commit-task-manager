package route

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/go-pg/pg/v10"
)

func Setup(api *huma.API, db *pg.DB) {

	controller := NewController(NewRepo(db))

	huma.Register(*api, huma.Operation{
		OperationID: "get-one-task",
		Method:      http.MethodGet,
		Path:        "/tasks/{id}",
		Summary:     "get one task",
		Description: "",
		Tags:        []string{"Task"},
	}, controller.getOne)

	huma.Register(*api, huma.Operation{
		OperationID: "set-task-status",
		Method:      http.MethodPut,
		Path:        "/tasks/{id}/set-status",
		Summary:     "set task status",
		Description: "",
		Tags:        []string{"Task"},
	}, controller.setStatus)

}
