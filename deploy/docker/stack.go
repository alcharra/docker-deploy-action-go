package docker

import (
	"fmt"
	"log"
	"path"

	"github.com/alcharra/docker-deploy-action-go/config"
	"github.com/alcharra/docker-deploy-action-go/ssh"
)

func DeployDockerStack(client *ssh.Client, cfg config.DeployConfig) {
	if cfg.Mode != "stack" {
		return
	}

	stack := cfg.StackName
	projectPath := cfg.ProjectPath
	rollback := cfg.EnableRollback
	deployFile := path.Base(cfg.DeployFile)

	cmd := fmt.Sprintf(`
		STACK="%s"
		PROJECT_PATH="%s"
		ENABLE_ROLLBACK="%t"

		cd "$PROJECT_PATH" || {
			echo "❌ Unable to access project directory: $PROJECT_PATH"
			exit 1
		}

		if ! docker info | grep -q "Swarm: active"; then
			echo "❌ Docker Swarm mode is not enabled"
			echo "👉 Run 'docker swarm init' on the server to activate it"
			exit 1
		fi

		echo "⚓ Deploying stack '$STACK' using Docker Swarm"
		docker stack deploy -c "%s" "$STACK" --with-registry-auth --detach=false

		echo "🔍 Checking status of services in stack '$STACK'"

		if ! docker service ls --filter "label=com.docker.stack.namespace=$STACK" | grep -v REPLICAS | grep -q " 0/"; then
			echo "✅ All services in stack '$STACK' are running as expected"
		else
			echo "❌ One or more services failed to start in stack '$STACK'"
			docker service ls --filter "label=com.docker.stack.namespace=$STACK"

			if [ "$ENABLE_ROLLBACK" = "true" ]; then
				echo "🔄 Attempting to roll back failed services..."
				for service in $(docker service ls --filter "label=com.docker.stack.namespace=$STACK" --format "{{.Name}}"); do
					echo "↩️ Rolling back service: $service"
					docker service update --rollback "$service" || echo "⚠️ Rollback failed for: $service"
				done
			else
				echo "⚠️ Rollback is disabled"
			fi

			exit 1
		fi
	`, stack, projectPath, rollback, deployFile)

	err := client.RunCommandStreamed(cmd)
	if err != nil {
		log.Fatalf("❌ Stack deployment failed: %v", err)
	}

	fmt.Println("✅ Stack deployment completed successfully")
}
