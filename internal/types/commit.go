package types

type CommitRequest struct {
	Token string `json:"token"`
	Body  struct {
		ProjectID     int64  `json:"project_id"`
		CommitMessage string `json:"commit_message"`
		CommitHash    string `json:"commit_hash"`
		Branch        string `json:"branch"`
	}
}

type CommitResponse struct {
}
