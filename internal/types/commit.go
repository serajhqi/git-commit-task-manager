package types

type CommitRequest struct {
	Token string `json:"token"`
	Body  struct {
		ProjectID     string `json:"project_id"`
		CommitMessage string `json:"commit_message"`
	}
}

type CommitResponse struct {
}
