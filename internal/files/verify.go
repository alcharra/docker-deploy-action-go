package files

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/alcharra/docker-deploy-action-go/config"
	"github.com/alcharra/docker-deploy-action-go/internal/logs"
	"github.com/alcharra/docker-deploy-action-go/internal/ssh/client"
)

func CheckFilesExistRemote(cli *client.Client, cfg config.DeployConfig, files []UploadedFile) {
	logs.IsVerbose = cfg.Verbose
	logs.Step("🧪 Verifying uploaded files...")

	for _, file := range files {
		remotePath := filepath.ToSlash(file.RemotePath)

		logs.Verbosef("Checking if remote file exists: %s", remotePath)
		logs.VerboseCommandf("stat %s", remotePath)

		cmd := fmt.Sprintf(`
			if stat "%s" >/dev/null 2>&1; then
				echo "OK"
			else
				echo "MISSING"
			fi
		`, remotePath)

		stdout, stderr, err := cli.RunCommandBuffered(cmd)
		if err != nil {
			logs.Fatalf("Unable to verify remote file '%s': %v\nDetails: %s", remotePath, err, stderr)
		}

		switch strings.TrimSpace(stdout) {
		case "OK":
			logs.Success(remotePath)
		case "MISSING":
			logs.Fatalf("File missing after upload: %s", remotePath)
		default:
			logs.Fatalf("Unexpected verification response for: %s", remotePath)
		}
	}
}
