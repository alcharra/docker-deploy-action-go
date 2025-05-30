package docker

import (
	"fmt"
	"strings"

	"github.com/alcharra/docker-deploy-action-go/config"
	"github.com/alcharra/docker-deploy-action-go/internal/logs"
	"github.com/alcharra/docker-deploy-action-go/internal/ssh/client"
)

func DockerRegistryLogin(cli *client.Client, cfg config.DeployConfig) {
	if cfg.RegistryHost == "" || cfg.RegistryUser == "" || cfg.RegistryPass == "" {
		return
	}

	logs.IsVerbose = cfg.Verbose
	logs.Step("\U0001F510 Docker registry login...")

	masked := strings.Repeat("*", len(cfg.RegistryPass))

	logs.Verbosef("Attempting login to registry: %s", cfg.RegistryHost)
	logs.VerboseCommandf(`echo "%s" | docker login %s -u %s --password-stdin`, masked, cfg.RegistryHost, cfg.RegistryUser)

	cmd := fmt.Sprintf(`
		echo "%s" | docker login "%s" -u "%s" --password-stdin >/dev/null
	`, cfg.RegistryPass, cfg.RegistryHost, cfg.RegistryUser)

	_, stderr, err := cli.RunCommandBuffered(cmd)
	if err != nil {
		logs.Fatalf("Registry login failed: %v\nDetails: %s", err, stderr)
	}

	logs.Successf("Logged in to: %s", cfg.RegistryHost)
}
