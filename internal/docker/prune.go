package docker

import (
	"strings"

	"github.com/alcharra/docker-deploy-action-go/config"
	"github.com/alcharra/docker-deploy-action-go/internal/logs"
	"github.com/alcharra/docker-deploy-action-go/internal/ssh/client"
)

func RunDockerPrune(cli *client.Client, cfg config.DeployConfig) {
	logs.IsVerbose = cfg.Verbose

	pruneType := strings.ToLower(cfg.DockerPrune)

	if pruneType == "" || pruneType == "none" {
		return
	}

	var cmd string

	switch pruneType {
	case "system":
		cmd = "docker system prune -f"
	case "volumes":
		cmd = "docker volume prune -f"
	case "networks":
		cmd = "docker network prune -f"
	case "images":
		cmd = "docker image prune -f"
	case "containers":
		cmd = "docker container prune -f"
	default:
		logs.Fatalf("Invalid prune type: '%s'. Accepted values are: system, volumes, networks, images, containers, or none.", pruneType)
	}

	logs.Step("\U0001F9F9 Docker prune...")
	logs.Verbose("Running Docker prune command...")
	logs.VerboseCommandf(cmd)
	logs.Substepf("\u2022 Prune type: %s", pruneType)

	stdout, stderr, err := cli.RunCommandBuffered(cmd)
	if err != nil {
		logs.Fatalf("Docker prune command failed: %v\nDetails: %s", err, stderr)
	}

	lines := strings.Split(strings.TrimSpace(stdout), "\n")
	var lastHeader string

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		switch {
		case strings.HasPrefix(line, "Deleted ") || strings.HasPrefix(line, "Unused "):
			lastHeader = "\u2022 " + line
			logs.Substep(lastHeader)
		case strings.HasPrefix(line, "Total reclaimed space:"):
			space := strings.TrimPrefix(line, "Total reclaimed space: ")
			logs.Substepf("\u2022 Reclaimed space: %s", space)
			lastHeader = ""
		case strings.HasPrefix(line, "No "):
			logs.Substepf("\u2022 %s", line)
			lastHeader = ""
		default:
			if lastHeader != "" {
				logs.Substepf("   \u2192 %s", line)
			} else {
				logs.Substepf("\u2022 %s", line)
			}
		}
	}

	logs.Success("Docker prune completed successfully")
}
