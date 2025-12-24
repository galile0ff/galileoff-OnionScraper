package classifier

type ClassificationConfig struct {
	Categories []Category `yaml:"categories"`
}

type Category struct {
	ID             string          `yaml:"id"`
	Name           string          `yaml:"name"`
	Tag            string          `yaml:"tag"`
	Color          string          `yaml:"color"`
	Priority       int             `yaml:"priority"`
	Keywords       KeywordRules    `yaml:"keywords"`
	StructureRules []StructureRule `yaml:"structure_rules"`
	MaxLinks       int             `yaml:"max_links"`
}

type KeywordRules struct {
	High    []string `yaml:"high"`
	Medium  []string `yaml:"medium"`
	Exclude []string `yaml:"exclude"`
}

type StructureRule struct {
	Selector string `yaml:"selector"`
}

type Result struct {
	CategoryID string
	Tag        string
	Color      string
	Score      int
	IsUnknown  bool
}
