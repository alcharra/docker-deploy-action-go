package docker

import (
	"fmt"
	"log"
	"strings"

	"github.com/alcharra/docker-deploy-action-go/ssh"
)

func CheckDockerInstalled(client *ssh.Client) {
	cmd := `
		if ! command -v docker >/dev/null 2>&1; then
			echo "❌ Docker is not installed or not in PATH."
			exit 1
		else
			echo "✅ Docker is installed and available."
		fi
	`

	stdout, stderr, err := client.RunCommandBuffered(cmd)
	if err != nil {
		log.Fatalf("❌ Docker check failed: %v\nStderr: %s", err, stderr)
	}

	fmt.Println(strings.TrimSpace(stdout))
}
