package files

import (
	"fmt"
	"log"
	"strings"

	"github.com/alcharra/docker-deploy-action-go/config"
	"github.com/alcharra/docker-deploy-action-go/ssh"
)

func CheckOrCreateRemotePath(client *ssh.Client, cfg config.DeployConfig) {
	cmd := fmt.Sprintf(`
		PROJECT_PATH="%s"
		SSH_USER="%s"

		if [ ! -d "$PROJECT_PATH" ]; then
			echo '📁 Project path not found - creating it...'
			mkdir -p "$PROJECT_PATH"
			chown "$SSH_USER":"$SSH_USER" "$PROJECT_PATH"
			chmod 750 "$PROJECT_PATH"

			if [ ! -d "$PROJECT_PATH" ]; then
				echo '❌ Failed to create project path!'
				exit 1
			fi

			echo '✅ Project path created and verified.'
		else
			echo '✅ Project path already exists.'
		fi
	`, cfg.ProjectPath, cfg.SSHUser)

	stdout, stderr, err := client.RunCommandBuffered(cmd)
	if err != nil {
		log.Fatalf("❌ Failed to check/create project path: %v\nStderr: %s", err, stderr)
	}

	fmt.Println(strings.TrimSpace(stdout))
}
