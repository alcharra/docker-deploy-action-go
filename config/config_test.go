//go:build unit
// +build unit

package config

import (
	"os"
	"reflect"
	"testing"
)

func TestLoadConfig_Defaults(t *testing.T) {
	os.Clearenv()
	cfg := LoadConfig()

	if cfg.DeployFile != "docker-compose.yml" {
		t.Errorf("expected DeployFile to be 'docker-compose.yml', got %s", cfg.DeployFile)
	}
	if cfg.Mode != "compose" {
		t.Errorf("expected Mode to be 'compose', got %s", cfg.Mode)
	}
	if cfg.DockerNetworkAttach {
		t.Errorf("expected DockerNetworkAttach to be false, got true")
	}
	if cfg.EnableRollback {
		t.Errorf("expected EnableRollback to be false, got true")
	}
	if cfg.SSHTimeout != "10s" {
		t.Errorf("expected SSHTimeout to default to '10s', got %s", cfg.SSHTimeout)
	}
	if len(cfg.ExtraFiles) != 0 {
		t.Errorf("expected ExtraFiles to be empty, got %v", cfg.ExtraFiles)
	}
	if len(cfg.ComposeTargetServices) != 0 {
		t.Errorf("expected ComposeTargetServices to be empty, got %v", cfg.ComposeTargetServices)
	}
}

func TestLoadConfig_WithEnvOverrides(t *testing.T) {
	t.Setenv("DEPLOY_FILE", "override.yml")
	t.Setenv("MODE", "stack")
	t.Setenv("DOCKER_NETWORK_ATTACHABLE", "true")
	t.Setenv("ENABLE_ROLLBACK", "true")
	t.Setenv("SSH_TIMEOUT", "20s")

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
	if cfg.SSHTimeout != "20s" {
		t.Errorf("expected SSHTimeout to be '20s', got %s", cfg.SSHTimeout)
	}
}

func TestLoadConfig_SliceParsing_Newline(t *testing.T) {
	t.Setenv("EXTRA_FILES", `
		flatten ./tests/testdata/stack/redis.conf
		flatten ./tests/testdata/stack/nginx.conf:dir/
	`)
	t.Setenv("COMPOSE_TARGET_SERVICES", `
		web
		db
	`)

	cfg := LoadConfig()

	expectedExtra := []ExtraFile{
		{Src: "./tests/testdata/stack/redis.conf", Dst: "", Flatten: true},
		{Src: "./tests/testdata/stack/nginx.conf", Dst: "dir/", Flatten: true},
	}
	expectedServices := []string{"web", "db"}

	if !reflect.DeepEqual(cfg.ExtraFiles, expectedExtra) {
		t.Errorf("expected ExtraFiles to be %v, got %v", expectedExtra, cfg.ExtraFiles)
	}

	if !reflect.DeepEqual(cfg.ComposeTargetServices, expectedServices) {
		t.Errorf("expected ComposeTargetServices to be %v, got %v", expectedServices, cfg.ComposeTargetServices)
	}
}

func TestLoadConfig_EmptySlices(t *testing.T) {
	t.Setenv("EXTRA_FILES", "")
	t.Setenv("COMPOSE_TARGET_SERVICES", "")

	cfg := LoadConfig()

	if len(cfg.ExtraFiles) != 0 {
		t.Errorf("expected ExtraFiles to be empty, got %v", cfg.ExtraFiles)
	}
	if len(cfg.ComposeTargetServices) != 0 {
		t.Errorf("expected ComposeTargetServices to be empty, got %v", cfg.ComposeTargetServices)
	}
}

func TestLoadConfig_BoolParsing(t *testing.T) {
	t.Setenv("COMPOSE_BUILD", "true")
	if !LoadConfig().ComposeBuild {
		t.Error("expected ComposeBuild to be true for 'true'")
	}

	t.Setenv("COMPOSE_BUILD", "false")
	if LoadConfig().ComposeBuild {
		t.Error("expected ComposeBuild to be false for 'false'")
	}

	t.Setenv("COMPOSE_BUILD", "yes")
	if LoadConfig().ComposeBuild {
		t.Error("expected ComposeBuild to be false for non-'true' input")
	}
}

func TestLoadConfig_AllFieldsSet(t *testing.T) {
	t.Setenv("SSH_HOST", "example.com")
	t.Setenv("SSH_PORT", "2222")
	t.Setenv("SSH_USER", "deployer")
	t.Setenv("SSH_KEY", "id_rsa")
	t.Setenv("SSH_KEY_PASSPHRASE", "secret")
	t.Setenv("SSH_KNOWN_HOSTS", "known_hosts")
	t.Setenv("SSH_FINGERPRINT", "fp123")
	t.Setenv("SSH_TIMEOUT", "30s")
	t.Setenv("PROJECT_PATH", "/app")
	t.Setenv("DEPLOY_FILE", "docker-stack.yml")
	t.Setenv("MODE", "stack")
	t.Setenv("STACK_NAME", "my-stack")
	t.Setenv("COMPOSE_PULL", "false")
	t.Setenv("COMPOSE_BUILD", "true")
	t.Setenv("COMPOSE_NO_DEPS", "true")
	t.Setenv("DOCKER_NETWORK", "custom-net")
	t.Setenv("DOCKER_NETWORK_DRIVER", "overlay")
	t.Setenv("DOCKER_NETWORK_ATTACHABLE", "true")
	t.Setenv("DOCKER_PRUNE", "all")
	t.Setenv("REGISTRY_HOST", "docker.io")
	t.Setenv("REGISTRY_USER", "admin")
	t.Setenv("REGISTRY_PASS", "hunter2")
	t.Setenv("ENABLE_ROLLBACK", "true")
	t.Setenv("ENV_VARS", "FOO=bar")
	t.Setenv("EXTRA_FILES", "file1.env\nfile2.env")
	t.Setenv("COMPOSE_TARGET_SERVICES", "web\nworker")

	cfg := LoadConfig()

	if cfg.SSHHost != "example.com" || cfg.SSHUser != "deployer" || cfg.ProjectPath != "/app" {
		t.Errorf("unexpected SSH or project config values: %+v", cfg)
	}
	if !cfg.ComposeBuild || !cfg.ComposeNoDeps || cfg.ComposePull {
		t.Errorf("unexpected Compose config values: %+v", cfg)
	}
	if cfg.DockerNetwork != "custom-net" || cfg.DockerNetworkDriver != "overlay" || !cfg.DockerNetworkAttach {
		t.Errorf("unexpected Docker network config: %+v", cfg)
	}
	if cfg.RegistryPass != "hunter2" || cfg.EnvVars != "FOO=bar" {
		t.Errorf("unexpected registry or env config: %+v", cfg)
	}
}
