package docker

import (
	"fmt"
	"log"

	"github.com/alcharra/docker-deploy-action-go/config"
	"github.com/alcharra/docker-deploy-action-go/ssh"
)

func DeployDockerCompose(client *ssh.Client, cfg config.DeployConfig) {
	if cfg.Mode != "compose" {
		return
	}

	projectPath := cfg.ProjectPath

	cmd := fmt.Sprintf(`
		PROJECT_PATH="%s"
		ENABLE_ROLLBACK="%t"
		COMPOSE_PULL="%t"

		echo "🐳 Deploying using Docker Compose"

		if docker compose version >/dev/null 2>&1; then
			COMPOSE="docker compose"
		elif docker-compose version >/dev/null 2>&1; then
			COMPOSE="docker-compose"
		else
			echo "❌ Docker Compose not found! Please install it."
			exit 1
		fi

		cd "$PROJECT_PATH" || { echo "❌ Failed to cd into $PROJECT_PATH"; exit 1; }

		if [ "$COMPOSE_PULL" = "true" ]; then
			$COMPOSE pull || { echo "❌ Failed to pull images"; exit 1; }
		else
			echo "⏩ Skipping image pull"
		fi

		$COMPOSE down &&
		$COMPOSE up -d

		echo "✅ Verifying Compose services"
		if $COMPOSE ps | grep -E "Exit|Restarting|Dead"; then
			echo "❌ One or more services failed to start!"
			$COMPOSE ps

			if [ "$ENABLE_ROLLBACK" = "true" ]; then
				echo "🔄 Attempting rollback..."

				LATEST_BACKUP=$(ls -td .backup_* 2>/dev/null | head -n 1)

				if [ -n "$LATEST_BACKUP" ]; then
					echo "📦 Restoring backup from $LATEST_BACKUP"
					cp "$LATEST_BACKUP"/* . || echo "❌ Failed to restore backup files"

					echo "♻️ Re-running deployment after rollback"
					$COMPOSE down
					$COMPOSE up -d

					if $COMPOSE ps | grep -E "Exit|Restarting|Dead"; then
						echo "❌ Rollback deployment failed"
						$COMPOSE ps
					else
						echo "✅ Rollback deployment successful"
					fi
				else
					echo "⚠️ No backup found to restore"
				fi
				echo "🧼 Cleaning up all backups"
				rm -rf .backup_* 2>/dev/null || true
			else
				echo "⚠️ Rollback is disabled"
			fi
			exit 1
		else
			echo "✅ All services are running"
		fi
	`, projectPath, cfg.EnableRollback, cfg.ComposePull)

	err := client.RunCommandStreamed(cmd)
	if err != nil {
		log.Fatalf("❌ Compose deployment failed: %v", err)
	}

	fmt.Println("✅ Compose deployment completed")
}
