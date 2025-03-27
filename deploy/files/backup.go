package files

import (
	"fmt"
	"log"
	"path"
	"time"

	"github.com/alcharra/docker-deploy-action-go/config"
	"github.com/alcharra/docker-deploy-action-go/ssh"
)

func BackupDeploymentFiles(client *ssh.Client, cfg config.DeployConfig) {
	if cfg.Mode != "compose" || !cfg.EnableRollback {
		return
	}

	timestamp := time.Now().Format("20060102_150405")
	backupDirName := ".backup_" + timestamp
	backupDir := path.Join(cfg.ProjectPath, backupDirName)

	mainFile := path.Base(cfg.DeployFile)
	mainFilePath := path.Join(cfg.ProjectPath, mainFile)

	cmd := fmt.Sprintf(`
		echo "📦 Creating backup of deployment files"
		mkdir -p "%[1]s"

		if [ -f "%[2]s" ]; then
			cp "%[2]s" "%[1]s/"
		else
			echo "⚠️ Main deploy file '%[3]s' not found, skipping backup."
		fi
	`, backupDir, mainFilePath, mainFile)

	for _, file := range cfg.ExtraFiles {
		fileName := path.Base(file)
		filePath := path.Join(cfg.ProjectPath, fileName)

		cmd += fmt.Sprintf(`
		if [ -f "%[1]s" ]; then
			cp "%[1]s" "%[2]s/"
		else
			echo "⚠️ Extra file '%[3]s' not found, skipping."
		fi
		`, filePath, backupDir, fileName)
	}

	cmd += fmt.Sprintf(`
		echo "✅ Backup created at %s"
	`, backupDir)

	if err := client.RunCommandStreamed(cmd); err != nil {
		log.Printf("⚠️ Backup step failed: %v", err)
	}
}
