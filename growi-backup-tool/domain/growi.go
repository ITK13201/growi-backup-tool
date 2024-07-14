package domain

type GrowiRevision struct {
	ID     string `json:"id"`
	Body   string `json:"body"`
	PageID string `json:"page-id"`
}

type GrowiPage struct {
	ID               string  `json:"id"`
	Path             string  `json:"path"`
	LatestRevisionID *string `json:"latest-revision-id"`
	IsDumped         bool    `json:"is-dumped"`
}
