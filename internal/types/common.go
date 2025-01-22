package types

type GetAllRequest struct {
	Authorization string `header:"Authorization"`
	Offset        int    `query:"offset,minimum=0"`
	Limit         int    `query:"limit,minimum=1,maximum=1000"`
}
