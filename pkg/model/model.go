package model

type LockData struct {
	LockID    string `json:"ID"`
	Operation string `json:"Operation"`
	Info      string `json:"Info"`
	Who       string `json:"Who"`
	Version   string `json:"Version"`
	Created   string `json:"Created"`
	Path      string `json:"Path"`
}

type State struct {
	Version          float64        `json:"version"`
	TerraformVersion string         `json:"terraform_version"`
	Serial           float64        `json:"serial"`
	Lineage          string         `json:"lineage"`
	Outputs          map[string]any `json:"outputs"`
	Resources        []Resource     `json:"resources"`
	CheckResults     any            `json:"check_results"`
}

type Resource struct {
	Type      string     `json:"type"`
	Name      string     `json:"name"`
	Mode      string     `json:"mode"`
	Provider  string     `json:"provider"`
	Instances []Instance `json:"instances"`
}

type Instance struct {
	SchemaVersion       float64        `json:"schema_version"`
	Attributes          map[string]any `json:"attributes"`
	SensitiveAttributes []any          `json:"sensitive_attributes"`
}
