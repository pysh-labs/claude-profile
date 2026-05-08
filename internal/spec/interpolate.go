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
		env := p.MCPServers[name].Env
		envKeys := make([]string, 0, len(env))
		for k := range env {
			envKeys = append(envKeys, k)
		}
		sort.Strings(envKeys)

		for _, k := range envKeys {
			expanded, err := expand(env[k])
			if err != nil {
				errs = append(errs,
					fmt.Sprintf("mcp_servers.%s.env.%s: %s", name, k, err.Error()))
				continue
			}
			env[k] = expanded
		}
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
		seen := map[string]bool{}
		uniq := missing[:0]
		for _, m := range missing {
			if !seen[m] {
				seen[m] = true
				uniq = append(uniq, m)
			}
		}
		return "", fmt.Errorf("environment variables not set: %s", strings.Join(uniq, ", "))
	}
	return out, nil
}
