package validator

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

func LoadComposeFile(path string) (*ComposeFile, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	content := os.ExpandEnv(string(data))

	var root yaml.Node
	if err := yaml.Unmarshal([]byte(content), &root); err != nil {
		return nil, fmt.Errorf("failed to parse YAML:\n   \u2192 %s", err)
	}

	var cfg ComposeFile
	if err := root.Decode(&cfg); err != nil {
		return nil, fmt.Errorf("failed to decode Compose structure:\n   \u2192 %s", err)
	}

	return &cfg, nil
}

func (c *ComposeFile) Validate() error {
	if len(c.Services) == 0 {
		return fmt.Errorf("no services defined")
	}

	var errs []string

	for name, svc := range c.Services {
		if svc.Image == "" {
			errs = append(errs, fmt.Sprintf("service '%s' is missing 'image'", name))
		}
		if svc.Build != nil {
			errs = append(errs, fmt.Sprintf("service '%s' uses 'build', which is not supported in docker stack deploy", name))
		}
		if svc.Deploy != nil && svc.Deploy.Replicas != nil && *svc.Deploy.Replicas < 1 {
			errs = append(errs, fmt.Sprintf("service '%s' has invalid 'deploy.replicas': must be >= 1", name))
		}
		if svc.Deploy != nil && svc.Deploy.Placement != nil {
			for _, constraint := range svc.Deploy.Placement.Constraints {
				if !strings.Contains(constraint, "==") && !strings.Contains(constraint, "!=") {
					errs = append(errs, fmt.Sprintf("service '%s' has invalid constraint '%s': must contain '==' or '!='", name, constraint))
				}
			}
		}
		for i, port := range svc.Ports {
			switch p := port.(type) {
			case string:
				if !strings.Contains(p, ":") {
					errs = append(errs, fmt.Sprintf("service '%s' port '%s' must contain at least one ':'", name, p))
				}
			default:
				errs = append(errs, fmt.Sprintf("service '%s' ports[%d] must be a string like 'HOST:CONTAINER'", name, i))
			}
		}
		switch cmd := svc.Command.(type) {
		case string:
		case []interface{}:
			for i, part := range cmd {
				if _, ok := part.(string); !ok {
					errs = append(errs, fmt.Sprintf("service '%s' command[%d] must be a string", name, i))
				}
			}
		case nil:
		default:
			errs = append(errs, fmt.Sprintf("service '%s' command must be a string or list of strings", name))
		}
		for _, ref := range svc.Configs {
			if ref.Source == "" {
				errs = append(errs, fmt.Sprintf("service '%s' references a config with no 'source'", name))
			} else if c.Configs == nil || c.Configs[ref.Source] == nil {
				errs = append(errs, fmt.Sprintf("service '%s' references undefined config '%s'", name, ref.Source))
			}
		}
		for _, ref := range svc.Secrets {
			if ref.Source == "" {
				errs = append(errs, fmt.Sprintf("service '%s' references a secret with no 'source'", name))
			} else if c.Secrets == nil || c.Secrets[ref.Source] == nil {
				errs = append(errs, fmt.Sprintf("service '%s' references undefined secret '%s'", name, ref.Source))
			}
		}
	}

	checkExternal := func(kind, name string, val interface{}) {
		switch val.(type) {
		case bool, map[string]interface{}:
		default:
			errs = append(errs, fmt.Sprintf("%s '%s' has invalid 'external' value", kind, name))
		}
	}

	for name, def := range c.Volumes {
		if def != nil && def.External != nil {
			checkExternal("volume", name, def.External)
		}
	}
	for name, def := range c.Networks {
		if def != nil && def.External != nil {
			checkExternal("network", name, def.External)
		}
	}
	for name, def := range c.Configs {
		if def != nil && def.External != nil {
			checkExternal("config", name, def.External)
		}
	}
	for name, def := range c.Secrets {
		if def != nil && def.External != nil {
			checkExternal("secret", name, def.External)
		}
	}

	if len(errs) > 0 {
		return errors.New("validation failed:\n      \u2192 " + strings.Join(errs, "\n      \u2192 "))
	}

	return nil
}
