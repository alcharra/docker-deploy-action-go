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

	deployFile := path.Base(cfg.DeployFile)

	cmd := fmt.Sprintf(`
		STACK="%s"
		PROJECT_PATH="%s"
		ENABLE_ROLLBACK="%t"
		ENV_VARS='%s'

		if ! cd "$PROJECT_PATH"; then
			echo "❌ Failed to change directory to $PROJECT_PATH"
			exit 1
		fi

		if ! docker info | grep -q "Swarm: active"; then
			echo "❌ Docker Swarm mode is not enabled"
			echo "👉 Run 'docker swarm init' on the server to activate it"
			exit 1
		fi

		if [ -f ".env" ] && [ -n "${ENV_VARS}" ]; then
			echo "📄 Loading environment variables from .env"
			set -a
			source .env
			set +a
		fi

		echo "⚓ Deploying stack '$STACK' using Docker Swarm"

		DEPLOY_OUTPUT=$(mktemp)

		docker stack deploy -c "%s" "$STACK" --with-registry-auth --detach=false 2>&1 | tee "$DEPLOY_OUTPUT"

		# Check for known critical issues
		echo "🧪 Validating Stack file"
		
		if grep -Eqi "undefined volume|unsupported option|is not supported|no such file|error:" "$DEPLOY_OUTPUT"; then
			echo "❌ Stack deployment failed: validation error detected"
			echo "🔍 Reason:"
			grep -Ei "undefined volume|unsupported option|is not supported|no such file|error:" "$DEPLOY_OUTPUT"
			rm "$DEPLOY_OUTPUT"
			exit 1
		else
			echo "✅ Stack file is valid"
		fi

		rm "$DEPLOY_OUTPUT"

		echo "🔍 Verifying services in stack '$STACK'"

		if ! docker service ls --filter "label=com.docker.stack.namespace=$STACK" | grep -v REPLICAS | grep -q " 0/"; then
			echo "✅ All services in stack '$STACK' are running as expected"
		else
			echo "❌ One or more services failed to start in stack '$STACK'"
			docker service ls --filter "label=com.docker.stack.namespace=$STACK"

			if [ "$ENABLE_ROLLBACK" = "true" ]; then
				echo "🔄 Attempting to roll back failed services..."
				for service in $(docker service ls --filter "label=com.docker.stack.namespace=$STACK" --format "{{.Name}}"); do
					echo "↩️ Rolling back service: $service"
					if ! docker service update --rollback "$service"; then
						echo "⚠️ Rollback failed for: $service"
					fi
				done
			fi

			exit 1
		fi
	`, cfg.StackName, cfg.ProjectPath, cfg.EnableRollback, cfg.EnvVars, deployFile)

	err := client.RunCommandStreamed(cmd)
	if err != nil {
		log.Fatalf("❌ Stack deployment failed: %v", err)
	}

	fmt.Println("✅ Stack deployment completed successfully")
}
