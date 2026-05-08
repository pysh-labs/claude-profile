package spec

import (
	"fmt"
	"os"
	"regexp"
	"sort"
	"strings"
)

var envVarRe = regexp.MustCompile(`\$\{([A-Za-z_][A-Za-z0-9_]*)\}`)

func Interpolate(p *Profile) error {
	var errs []string

	serverNames := make([]string, 0, len(p.MCPServers))
	for name := range p.MCPServers {
		serverNames = append(serverNames, name)
	}
	sort.Strings(serverNames)

	for _, name := range serverNames {
		srv := p.MCPServers[name]
		envKeys := make([]string, 0, len(srv.Env))
		for k := range srv.Env {
			envKeys = append(envKeys, k)
		}
		sort.Strings(envKeys)

		for _, k := range envKeys {
			expanded, err := expand(srv.Env[k])
			if err != nil {
				errs = append(errs,
					fmt.Sprintf("mcp_servers.%s.env.%s: %s", name, k, err.Error()))
				continue
			}
			srv.Env[k] = expanded
		}
		p.MCPServers[name] = srv
	}

	if len(errs) > 0 {
		return fmt.Errorf("env interpolation failed:\n  - %s", strings.Join(errs, "\n  - "))
	}
	return nil
}

func expand(s string) (string, error) {
	var missing []string
	out := envVarRe.ReplaceAllStringFunc(s, func(match string) string {
		name := match[2 : len(match)-1]
		val, ok := os.LookupEnv(name)
		if !ok {
			missing = append(missing, name)
			return match
		}
		return val
	})
	if len(missing) > 0 {
		return "", fmt.Errorf("environment variable %s is not set", missing[0])
	}
	return out, nil
}
