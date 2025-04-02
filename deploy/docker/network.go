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

		echo "🌐 Checking if Docker network '$NETWORK' exists"

		if [ -z "$DRIVER" ]; then
			echo "❌ Network driver is not set"
			exit 1
		fi

		if docker network inspect "$NETWORK" > /dev/null 2>&1; then
			echo "✅ Network '$NETWORK' already exists - verifying driver..."

			EXISTING_DRIVER=$(docker network inspect --format '{{ .Driver }}' "$NETWORK")

			if [ "$EXISTING_DRIVER" != "$DRIVER" ]; then
				echo "⚠️ Network '$NETWORK' is using driver '$EXISTING_DRIVER', expected '$DRIVER'"
				echo "ℹ️ Consider removing and recreating it to avoid unexpected behaviour"
			else
				echo "✅ Driver matches expected: '$DRIVER'"
			fi
		else
			echo "🔧 Creating network '$NETWORK' with driver '$DRIVER'"

			if [ "$DRIVER" = "overlay" ] && [ "$MODE" = "stack" ] && ! docker info | grep -q 'Swarm: active'; then
				echo "⚠️ Swarm mode is not active - overlay networks require Swarm for multi-node setups"
				echo "ℹ️ The network will still work in single-node mode as a bridge"
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
				echo "✅ Network '$NETWORK' created successfully"
			else
				echo "❌ Failed to create network '$NETWORK'"
				exit 1
			fi
		fi
	`, cfg.DockerNetwork, cfg.DockerNetworkDriver, cfg.Mode, cfg.DockerNetworkAttach)

	stdout, stderr, err := client.RunCommandBuffered(cmd)
	if err != nil {
		log.Fatalf("❌ Could not ensure Docker network: %v\nStderr: %s", err, stderr)
	}

	fmt.Println(strings.TrimSpace(stdout))
}
