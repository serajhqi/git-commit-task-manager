package types

type GetAllRequest struct {
	Offset int `query:"offset"`
	Limit  int `query:"limit"`
}
