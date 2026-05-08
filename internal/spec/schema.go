package spec

type Profile struct {
	APIVersion        string                 `yaml:"apiVersion"        validate:"required,eq=claude-profile.io/v1"`
	Kind              string                 `yaml:"kind"              validate:"required,eq=Profile"`
	Metadata          Metadata               `yaml:"metadata"          validate:"required"`
	Statusline        Statusline             `yaml:"statusline"        validate:"required"`
	Marketplaces      map[string]Marketplace `yaml:"marketplaces,omitempty"`
	Plugins           []string               `yaml:"plugins,omitempty"`
	MCPServers        map[string]MCPServer   `yaml:"mcp_servers,omitempty"`
	SettingsOverrides map[string]any         `yaml:"settings_overrides,omitempty"`
	ClaudeMD          string                 `yaml:"claude_md,omitempty"`
}

type Metadata struct {
	Name        string `yaml:"name"        validate:"required"`
	Description string `yaml:"description,omitempty"`
}

type Statusline struct {
	Label string `yaml:"label" validate:"required"`
	Color string `yaml:"color" validate:"required,oneof=red green yellow blue magenta cyan white"`
}

type Marketplace struct {
	Type string `yaml:"type" validate:"required,oneof=github"`
	Repo string `yaml:"repo" validate:"required"`
}

type MCPServer struct {
	Command string            `yaml:"command" validate:"required"`
	Args    []string          `yaml:"args,omitempty"`
	Env     map[string]string `yaml:"env,omitempty"`
}
