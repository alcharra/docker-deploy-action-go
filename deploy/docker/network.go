package docker

import (
	"fmt"
	"log"
	"strings"

	"github.com/alcharra/docker-deploy-action-go/config"
	"github.com/alcharra/docker-deploy-action-go/ssh"
)

func EnsureDockerNetwork(client *ssh.Client, cfg config.DeployConfig) {
	if cfg.DockerNetwork == "" {
		return
	}

	cmd := fmt.Sprintf(`
		NETWORK="%s"
		DRIVER="%s"
		MODE="%s"
		ATTACHABLE="%t"

		echo "🌐 Ensuring network $NETWORK exists"

		if [ -z "$DRIVER" ]; then
			echo "❌ DOCKER_NETWORK_DRIVER is not set!"
			exit 1
		fi

		if docker network inspect "$NETWORK" > /dev/null 2>&1; then
			echo "✅ Network $NETWORK exists. Checking driver..."

			EXISTING_DRIVER=$(docker network inspect --format '{{ .Driver }}' "$NETWORK")

			if [ "$EXISTING_DRIVER" != "$DRIVER" ]; then
				echo "⚠️ Network $NETWORK exists but uses driver '$EXISTING_DRIVER' instead of '$DRIVER'"
				echo "🚨 Consider deleting and recreating the network manually."
			else
				echo "✅ Network driver matches expected: $DRIVER"
			fi
		else
			echo "🔧 Creating $NETWORK network with driver $DRIVER"

			# Swarm warning if overlay + stack
			if [ "$DRIVER" = "overlay" ] && [ "$MODE" = "stack" ] && ! docker info | grep -q 'Swarm: active'; then
				echo "⚠️ Swarm mode is not active. Overlay networks need Swarm for multi-node communication."
				echo "ℹ️ It will still work in single-node mode as a bridge."
			fi

			CREATE_CMD="docker network create --driver $DRIVER"

			if [ "$DRIVER" = "overlay" ] && [ "$MODE" = "stack" ]; then
				CREATE_CMD="$CREATE_CMD --scope swarm"
				if [ "$ATTACHABLE" = "true" ]; then
					CREATE_CMD="$CREATE_CMD --attachable"
				fi
			fi

			$CREATE_CMD "$NETWORK"

			if docker network inspect "$NETWORK" > /dev/null 2>&1; then
				echo "✅ Network $NETWORK successfully created"
			else
				echo "❌ Network creation failed for $NETWORK!"
				exit 1
			fi
		fi
	`, cfg.DockerNetwork, cfg.DockerNetworkDriver, cfg.Mode, cfg.DockerNetworkAttach)

	stdout, stderr, err := client.RunCommandBuffered(cmd)
	if err != nil {
		log.Fatalf("❌ Failed to ensure Docker network: %v\nStderr: %s", err, stderr)
	}

	fmt.Println(strings.TrimSpace(stdout))
}
