package docker

import (
	"fmt"
	"log"
	"strings"

	"github.com/alcharra/docker-deploy-action-go/config"
	"github.com/alcharra/docker-deploy-action-go/ssh"
)

func RunDockerPrune(client *ssh.Client, cfg config.DeployConfig) {
	pruneType := strings.ToLower(cfg.DockerPrune)

	if pruneType == "" || pruneType == "none" {
		fmt.Println("⏭️ Skipping Docker prune - no type specified")
		return
	}

	var cmd string
	var label string

	switch pruneType {
	case "system":
		label = "🧹 Running full system prune"
		cmd = "docker system prune -f"
	case "volumes":
		label = "📦 Removing unused volumes"
		cmd = "docker volume prune -f"
	case "networks":
		label = "🌐 Cleaning up unused networks"
		cmd = "docker network prune -f"
	case "images":
		label = "🖼️ Removing unused images"
		cmd = "docker image prune -f"
	case "containers":
		label = "📦 Removing stopped containers"
		cmd = "docker container prune -f"
	default:
		log.Fatalf("❌ Invalid prune type: '%s'. Accepted values are: system, volumes, networks, images, containers, or none.", pruneType)
	}

	fmt.Println(label)

	err := client.RunCommandStreamed(cmd)
	if err != nil {
		log.Fatalf("❌ Docker prune command failed: %v", err)
	}
}
