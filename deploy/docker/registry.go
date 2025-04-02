package docker

import (
	"fmt"
	"log"
	"strings"

	"github.com/alcharra/docker-deploy-action-go/config"
	"github.com/alcharra/docker-deploy-action-go/ssh"
)

func DockerRegistryLogin(client *ssh.Client, cfg config.DeployConfig) {
	if cfg.RegistryHost == "" || cfg.RegistryUser == "" || cfg.RegistryPass == "" {
		fmt.Println("⏭️ Skipping registry login - credentials not provided")
		return
	}

	cmd := fmt.Sprintf(`
		echo "🔑 Attempting to log in to container registry: %s"
		echo "%s" | docker login "%s" -u "%s" --password-stdin
	`, cfg.RegistryHost, cfg.RegistryPass, cfg.RegistryHost, cfg.RegistryUser)

	stdout, stderr, err := client.RunCommandBuffered(cmd)
	if err != nil {
		log.Fatalf("❌ Registry login failed: %v\nDetails: %s", err, stderr)
	}

	fmt.Println(strings.TrimSpace(stdout))
}
