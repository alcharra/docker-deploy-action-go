package docker

import (
	"strings"

	"github.com/alcharra/docker-deploy-action-go/config"
	"github.com/alcharra/docker-deploy-action-go/internal/logs"
	"github.com/alcharra/docker-deploy-action-go/internal/ssh/client"
)

func CheckDockerRequirements(cli *client.Client, cfg *config.DeployConfig) {
	logs.IsVerbose = cfg.Verbose

	switch cfg.Mode {
	case "stack":
		logs.Step("\U0001F433 Docker Stack checks...")
	case "compose":
		logs.Step("\U0001F433 Docker Compose checks...")
	default:
		logs.Step("\U0001F433 Docker checks...")
	}

	CheckDockerInstalled(cli)

	if cfg.Mode == "stack" {
		CheckSwarmMode(cli)
	} else {
		CheckComposeAvailable(cli, cfg)
	}
}

func CheckDockerInstalled(cli *client.Client) {
	logs.Verbose("Checking: Docker binary availability")
	logs.VerboseCommand("command -v docker")

	cmd := `
		if ! command -v docker >/dev/null 2>&1; then
			echo "MISSING"
		else
			echo "OK"
		fi
	`

	stdout, stderr, err := cli.RunCommandBuffered(cmd)
	if err != nil {
		logs.Fatalf("Unable to verify Docker installation: %v\nDetails: %s", err, stderr)
	}

	switch strings.TrimSpace(stdout) {
	case "OK":
		logs.Success("Docker is installed and accessible")
	case "MISSING":
		logs.Fatalf("Docker is not installed or not available in the system PATH")
	default:
		logs.Fatalf("Unexpected response while checking for Docker: %s", stdout)
	}
}

func CheckSwarmMode(cli *client.Client) {
	logs.Verbose("Checking: Docker Swarm mode status")
	logs.VerboseCommand("docker info --format '{{ .Swarm.LocalNodeState }}'")

	cmd := `
		if docker info 2>/dev/null | grep -q 'Swarm: active'; then
			echo "OK"
		else
			echo "MISSING"
		fi
	`

	stdout, stderr, err := cli.RunCommandBuffered(cmd)
	if err != nil {
		logs.Fatalf("Unable to verify Swarm mode: %v\nDetails: %s", err, stderr)
	}

	switch strings.TrimSpace(stdout) {
	case "OK":
		logs.Success("Swarm mode is active")
	case "MISSING":
		logs.Fatalf("Swarm mode is not active (required for stack mode)")
	default:
		logs.Fatalf("Unexpected response when checking Swarm mode: %s", stdout)
	}
}

func CheckComposeAvailable(cli *client.Client, cfg *config.DeployConfig) {
	logs.Verbose("Checking: Docker Compose availability")
	logs.VerboseCommand("docker compose version || docker-compose version")

	cmd := `
		if command -v docker compose >/dev/null 2>&1; then
			echo "docker compose"
		elif command -v docker-compose >/dev/null 2>&1; then
			echo "docker-compose"
		else
			echo "MISSING"
		fi
	`

	stdout, stderr, err := cli.RunCommandBuffered(cmd)
	if err != nil {
		logs.Fatalf("Unable to verify Docker Compose availability: %v\nDetails: %s", err, stderr)
	}

	binary := strings.TrimSpace(stdout)
	switch binary {
	case "docker compose", "docker-compose":
		logs.Success("Docker Compose is available")
		cfg.ComposeBinary = binary
	case "MISSING":
		logs.Fatalf("Docker Compose is not installed or accessible")
	default:
		logs.Fatalf("Unexpected response when checking Compose: %s", stdout)
	}
}
