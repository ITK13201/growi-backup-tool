package domain

type ExpandConfig struct {
	PagesFileName     string `json:"pages-file-name"`
	RevisionsFileName string `json:"revisions-file-name"`
}

func NewExpandConfig() *ExpandConfig {
	return &ExpandConfig{
		PagesFileName:     "pages.json",
		RevisionsFileName: "revisions.json",
	}
}

type Config struct {
	Expand *ExpandConfig `json:"expand"`
}

func NewConfig() *Config {
	return &Config{
		Expand: NewExpandConfig(),
	}
}
