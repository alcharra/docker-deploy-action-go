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
		fmt.Println("⏭️ Skipping docker prune")
		return
	}

	var cmd string
	var label string

	switch pruneType {
	case "system":
		label = "🧹 Running full system prune"
		cmd = "docker system prune -f"
	case "volumes":
		label = "📦 Running volume prune"
		cmd = "docker volume prune -f"
	case "networks":
		label = "🌐 Running network prune"
		cmd = "docker network prune -f"
	case "images":
		label = "🖼️ Running image prune"
		cmd = "docker image prune -f"
	case "containers":
		label = "📦 Running container prune"
		cmd = "docker container prune -f"
	default:
		log.Fatalf("❌ Invalid prune type: %s", pruneType)
	}

	fmt.Println(label)

	err := client.RunCommandStreamed(cmd)
	if err != nil {
		log.Fatalf("❌ Docker prune failed: %v", err)
	}
}
