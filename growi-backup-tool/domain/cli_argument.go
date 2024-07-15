package domain

type CLIArgumentExpand struct {
	PagesFilePath     string `json:"pages-file-path"`
	RevisionsFilePath string `json:"revisions-file-path"`
	OutputDirPath     string `json:"output-dir-path"`
}
