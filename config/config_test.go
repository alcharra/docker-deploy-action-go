//go:build unit
// +build unit

package config

import (
	"os"
	"testing"
)

func TestLoadConfigDefaults(t *testing.T) {
	os.Clearenv()

	cfg := LoadConfig()

	if cfg.DeployFile != "docker-compose.yml" {
		t.Errorf("expected default DeployFile to be 'docker-compose.yml', got %s", cfg.DeployFile)
	}
	if cfg.Mode != "compose" {
		t.Errorf("expected default Mode to be 'compose', got %s", cfg.Mode)
	}
	if cfg.DockerNetworkAttach {
		t.Errorf("expected default DockerNetworkAttach to be false, got true")
	}
	if cfg.EnableRollback {
		t.Errorf("expected default EnableRollback to be false, got true")
	}
}

func TestLoadConfigWithEnvOverrides(t *testing.T) {
	t.Setenv("DEPLOY_FILE", "override.yml")
	t.Setenv("MODE", "stack")
	t.Setenv("DOCKER_NETWORK_ATTACHABLE", "true")
	t.Setenv("ENABLE_ROLLBACK", "true")

	cfg := LoadConfig()

	if cfg.DeployFile != "override.yml" {
		t.Errorf("expected DeployFile to be 'override.yml', got %s", cfg.DeployFile)
	}
	if cfg.Mode != "stack" {
		t.Errorf("expected Mode to be 'stack', got %s", cfg.Mode)
	}
	if !cfg.DockerNetworkAttach {
		t.Errorf("expected DockerNetworkAttach to be true, got false")
	}
	if !cfg.EnableRollback {
		t.Errorf("expected EnableRollback to be true, got false")
	}
}

func TestLoadConfig_SliceParsing(t *testing.T) {
	t.Setenv("EXTRA_FILES", "a.env,b.env")
	t.Setenv("COMPOSE_TARGET_SERVICES", "web,db")

	cfg := LoadConfig()

	if len(cfg.ExtraFiles) != 2 || cfg.ExtraFiles[0] != "a.env" || cfg.ExtraFiles[1] != "b.env" {
		t.Errorf("expected ExtraFiles to be [a.env b.env], got %#v", cfg.ExtraFiles)
	}

	if len(cfg.ComposeTargetServices) != 2 || cfg.ComposeTargetServices[1] != "db" {
		t.Errorf("expected ComposeTargetServices to be [web db], got %#v", cfg.ComposeTargetServices)
	}
}

func TestLoadConfig_DefaultTimeout(t *testing.T) {
	os.Clearenv()
	cfg := LoadConfig()
	if cfg.Timeout != "10s" {
		t.Errorf("expected Timeout to default to '10s', got %s", cfg.Timeout)
	}
}

func TestLoadConfig_AllFields(t *testing.T) {
	t.Setenv("SSH_HOST", "host")
	t.Setenv("SSH_PORT", "2022")
	t.Setenv("SSH_USER", "user")
	t.Setenv("SSH_KEY", "key")
	t.Setenv("FINGERPRINT", "fp")
	t.Setenv("TIMEOUT", "30s")
	t.Setenv("EXTRA_FILES", "a,b")
	t.Setenv("REGISTRY_PASS", "pass")

	cfg := LoadConfig()

	if cfg.SSHHost != "host" || cfg.SSHPort != "2022" || cfg.RegistryPass != "pass" {
		t.Errorf("unexpected config values: %+v", cfg)
	}
}
