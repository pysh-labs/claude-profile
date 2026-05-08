package spec

type Profile struct {
	APIVersion        string                 `yaml:"apiVersion"`
	Kind              string                 `yaml:"kind"`
	Metadata          Metadata               `yaml:"metadata"`
	Statusline        Statusline             `yaml:"statusline"`
	Marketplaces      map[string]Marketplace `yaml:"marketplaces,omitempty"`
	Plugins           []string               `yaml:"plugins,omitempty"`
	MCPServers        map[string]MCPServer   `yaml:"mcp_servers,omitempty"`
	SettingsOverrides map[string]any         `yaml:"settings_overrides,omitempty"`
	ClaudeMD          string                 `yaml:"claude_md,omitempty"`
}

type Metadata struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description,omitempty"`
}

type Statusline struct {
	Label string `yaml:"label"`
	Color string `yaml:"color"`
}

type Marketplace struct {
	Type string `yaml:"type"`
	Repo string `yaml:"repo"`
}

type MCPServer struct {
	Command string            `yaml:"command"`
	Args    []string          `yaml:"args,omitempty"`
	Env     map[string]string `yaml:"env,omitempty"`
}
