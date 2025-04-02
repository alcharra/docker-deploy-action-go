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

	// Start building the backup script
	cmd := fmt.Sprintf(`
		echo "📦 Creating backup of deployment files"
		mkdir -p "%[1]s"

		if [ -f "%[2]s" ]; then
			cp "%[2]s" "%[1]s/"
			echo "✅ Backed up main file: '%[3]s'"
		else
			echo "⚠️ Main deployment file '%[3]s' not found - skipping"
		fi
	`, backupDir, mainFilePath, mainFile)

	// Add any extra files
	for _, file := range cfg.ExtraFiles {
		fileName := path.Base(file)
		filePath := path.Join(cfg.ProjectPath, fileName)

		cmd += fmt.Sprintf(`
		if [ -f "%[1]s" ]; then
			cp "%[1]s" "%[2]s/"
			echo "✅ Backed up extra file: '%[3]s'"
		else
			echo "⚠️ Extra file '%[3]s' not found - skipping"
		fi
		`, filePath, backupDir, fileName)
	}

	cmd += fmt.Sprintf(`
		echo "📂 Backup directory created at: %s"
	`, backupDir)

	if err := client.RunCommandStreamed(cmd); err != nil {
		log.Printf("⚠️ Backup step encountered an issue: %v", err)
	}
}
