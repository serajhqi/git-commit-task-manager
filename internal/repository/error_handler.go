package repository

import (
	"errors"
	"fmt"

	"github.com/danielgtaylor/huma/v2"
)

// HandleError maps application errors to HTTP responses.
func HandleError(err error) error {
	if err == nil {
		return nil
	}

	switch {
	// Example for "not found" error
	case err != nil && err.Error() == "pg: no rows in result set":
		return huma.Error404NotFound("Resource not found")

	case errors.Is(err, fmt.Errorf("bad request")):
		return huma.Error400BadRequest("invalid input provided")

	default:
		return huma.Error500InternalServerError("internal error", err)
	}
}
