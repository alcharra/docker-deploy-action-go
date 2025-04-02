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
			FILE_PATH="%s"
			if [ ! -f "$FILE_PATH" ]; then
				echo "❌ File missing after upload: $FILE_PATH"
				exit 1
			else
				echo "✅ Verified file exists: $FILE_PATH"
			fi
		`, remotePath)

		stdout, stderr, err := client.RunCommandBuffered(cmd)
		if err != nil {
			log.Fatalf("❌ Remote file check failed for '%s': %v\nDetails: %s", filename, err, stderr)
		}

		fmt.Println(strings.TrimSpace(stdout))
	}
}
