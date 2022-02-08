package xoidc

type (
	Config struct {
		OIDCAuth AuthConfig `json:"OIDCAuth" yaml:"OIDCAuth"`
	}
	AuthConfig struct {
		Issuer         string   `json:"issuer" yaml:"issuer"`
		AllowProjects  []string `json:"allow_projects" yaml:"allow_projects"`
		Scopes         []string `json:"scopes" yaml:"scopes"`
		AllowedIssuers []string `json:"allowed_issuers" yaml:"allowed_issuers"`
	}
)
