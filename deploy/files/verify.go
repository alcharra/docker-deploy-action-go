package files

import (
	"fmt"
	"log"
	"path"
	"path/filepath"
	"strings"

	"github.com/alcharra/docker-deploy-action-go/ssh"
)

func CheckFilesExistRemote(client *ssh.Client, projectPath string, files []string) {
	for _, localFile := range files {
		filename := filepath.Base(localFile)
		remotePath := path.Join(projectPath, filename)

		cmd := fmt.Sprintf(`
		PATH="%s"
		if [ ! -f "$PATH" ]; then
			echo '❌ Missing file after upload: %s'
			exit 1
		else
			echo '✅ File verified: %s'
		fi
		`, remotePath, remotePath, remotePath)

		stdout, stderr, err := client.RunCommandBuffered(cmd)
		if err != nil {
			log.Fatalf("❌ Remote file check failed for %s: %v\nStderr: %s", filename, err, stderr)
		}

		fmt.Println(strings.TrimSpace(stdout))
	}
}
