package types

type GetAllRequest struct {
	Offset int `query:"offset"`
	Limit  int `query:"limit"`
}

type TokenizedRequest[T any] struct {
	Token string `header:"token"`
	Body  T
}
