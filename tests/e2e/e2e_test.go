//go:build e2e
// +build e2e

package tests

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"testing"

	"github.com/alcharra/docker-deploy-action-go/config"
)

func TestDeployBinary_MainDeploy(t *testing.T) {
	cfg := config.LoadConfig()

	binary := "deploy-action"
	if runtime.GOOS == "windows" {
		binary += ".exe"
	}

	binaryPath, err := filepath.Abs(filepath.Join("../../", binary))
	if err != nil {
		t.Fatalf("failed to resolve binary path: %v", err)
	}

	if _, err := os.Stat(binaryPath); os.IsNotExist(err) {
		t.Fatalf("deploy-action binary not found at %s", binaryPath)
	}

	cmd := exec.Command(binaryPath)
	cmd.Dir = "../../"

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cmd.Env = append(os.Environ(),
		"SSH_HOST="+cfg.SSHHost,
		"SSH_PORT="+cfg.SSHPort,
		"SSH_USER="+cfg.SSHUser,
		"SSH_KEY="+cfg.SSHKey,
		"SSH_KEY_PASSPHRASE="+cfg.SSHKeyPassphrase,
		"SSH_KNOWN_HOSTS="+cfg.SSHKnownHosts,
		"SSH_FINGERPRINT="+cfg.SSHFingerprint,
		"SSH_TIMEOUT="+cfg.SSHTimeout,
		"PROJECT_PATH="+cfg.ProjectPath,
		"DEPLOY_FILE="+cfg.DeployFile,
		"EXTRA_FILES="+strings.Join(extraFilesToEnv(cfg.ExtraFiles), "\n"),
		"MODE="+cfg.Mode,
		"STACK_NAME="+cfg.StackName,
		"COMPOSE_PULL="+strconv.FormatBool(cfg.ComposePull),
		"COMPOSE_BUILD="+strconv.FormatBool(cfg.ComposeBuild),
		"COMPOSE_NO_DEPS="+strconv.FormatBool(cfg.ComposeNoDeps),
		"COMPOSE_TARGET_SERVICES="+strings.Join(cfg.ComposeTargetServices, "\n"),
		"DOCKER_NETWORK="+cfg.DockerNetwork,
		"DOCKER_NETWORK_DRIVER="+cfg.DockerNetworkDriver,
		"DOCKER_NETWORK_ATTACHABLE="+strconv.FormatBool(cfg.DockerNetworkAttach),
		"DOCKER_PRUNE="+cfg.DockerPrune,
		"REGISTRY_HOST="+cfg.RegistryHost,
		"REGISTRY_USER="+cfg.RegistryUser,
		"REGISTRY_PASS="+cfg.RegistryPass,
		"ENABLE_ROLLBACK="+strconv.FormatBool(cfg.EnableRollback),
		"ENV_VARS="+cfg.EnvVars,
		"VERBOSE="+strconv.FormatBool(cfg.Verbose),
	)

	t.Logf("\U0001F680 Running E2E deploy with: %s", binaryPath)

	if err := cmd.Run(); err != nil {
		t.Fatalf("\U0001F6A8 deploy-action binary failed: %v", err)
	}

	t.Log("\u2705 E2E deployment test passed")
}
