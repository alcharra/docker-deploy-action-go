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
