package docker

import (
	"fmt"
	"log"
	"strings"

	"github.com/alcharra/docker-deploy-action-go/config"
	"github.com/alcharra/docker-deploy-action-go/ssh"
)

func DeployDockerCompose(client *ssh.Client, cfg config.DeployConfig) {
	if cfg.Mode != "compose" {
		return
	}

	projectPath := cfg.ProjectPath

	// Prepare flags for `docker compose up`
	var upFlags []string
	upFlags = append(upFlags, "-d")

	if cfg.ComposeBuild {
		upFlags = append(upFlags, "--build")
	}

	if cfg.ComposeNoDeps {
		upFlags = append(upFlags, "--no-deps")
	}

	upFlagsStr := strings.Join(upFlags, " ")

	// Build the up command
	var upCmd string
	if len(cfg.ComposeTargetServices) > 0 {
		var cmds []string
		for _, service := range cfg.ComposeTargetServices {
			cmds = append(cmds, fmt.Sprintf(`$COMPOSE up %s %s`, upFlagsStr, service))
		}
		upCmd = strings.Join(cmds, "\n")
	} else {
		upCmd = fmt.Sprintf(`$COMPOSE down && $COMPOSE up %s`, upFlagsStr)
	}

	cmd := fmt.Sprintf(`
		PROJECT_PATH="%s"
		ENABLE_ROLLBACK="%t"
		COMPOSE_PULL="%t"

		echo "🐳 Deploying with Docker Compose"

		if docker compose version >/dev/null 2>&1; then
			COMPOSE="docker compose"
		elif docker-compose version >/dev/null 2>&1; then
			COMPOSE="docker-compose"
		else
			echo "❌ Docker Compose is not installed"
			exit 1
		fi

		cd "$PROJECT_PATH" || { echo "❌ Failed to change directory to $PROJECT_PATH"; exit 1; }

		if [ "$COMPOSE_PULL" = "true" ]; then
			echo "📥 Pulling latest images"
			$COMPOSE pull || { echo "❌ Pull failed"; exit 1; }
		else
			echo "⏩ Skipping image pull"
		fi

		%s

		echo "🔍 Checking service status"
		if $COMPOSE ps | grep -E "Exit|Restarting|Dead"; then
			echo "❌ One or more services did not start properly"
			$COMPOSE ps

			if [ "$ENABLE_ROLLBACK" = "true" ]; then
				echo "🔄 Attempting to roll back"

				LATEST_BACKUP=$(ls -td .backup_* 2>/dev/null | head -n 1)

				if [ -n "$LATEST_BACKUP" ]; then
					echo "📦 Restoring backup from $LATEST_BACKUP"
					cp "$LATEST_BACKUP"/* . || echo "❌ Failed to restore backup"

					echo "♻️ Re-deploying previous version"
					$COMPOSE down
					$COMPOSE up -d

					if $COMPOSE ps | grep -E "Exit|Restarting|Dead"; then
						echo "❌ Rollback failed"
						$COMPOSE ps
					else
						echo "✅ Rollback successful"
					fi
				else
					echo "⚠️ No backup found"
				fi
			else
				echo "⚠️ Rollback is disabled"
			fi

			exit 1
		else
			echo "✅ All services are running"
		fi

		rm -rf .backup_* 2>/dev/null || true
	`, projectPath, cfg.EnableRollback, cfg.ComposePull, upCmd)

	err := client.RunCommandStreamed(cmd)
	if err != nil {
		log.Fatalf("❌ Docker Compose deployment failed: %v", err)
	}

	fmt.Println("✅ Docker Compose deployment completed")
}
