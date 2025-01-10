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

	// Unwrap error if it's wrapped
	unwrappedErr := errors.Unwrap(err)

	switch {
	// Example for "not found" error
	case unwrappedErr != nil && unwrappedErr.Error() == "pg: no rows in result set":
		return huma.Error404NotFound("Resource not found")

	case errors.Is(unwrappedErr, fmt.Errorf("bad request")):
		return huma.Error400BadRequest("invalid input provided")

	default:
		return huma.Error500InternalServerError("internal error", err)
	}
}
