package config

type Firestore struct {
	ProjectID          string `yaml:"project_id"`
	DatabaseID         string `yaml:"database_id"`
	JsonCredentialFile string `yaml:"json_credential_file"`
}
