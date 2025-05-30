package files

import (
	"fmt"
	"path"
	"strings"
	"time"

	"github.com/alcharra/docker-deploy-action-go/config"
	"github.com/alcharra/docker-deploy-action-go/internal/logs"
	"github.com/alcharra/docker-deploy-action-go/internal/ssh/client"
)

func BackupDeploymentFiles(cli *client.Client, cfg config.DeployConfig) {
	logs.IsVerbose = cfg.Verbose

	if cfg.Mode != "compose" || !cfg.EnableRollback {
		return
	}

	logs.Step("\U0001F4E4 Creating backup of deployment files...")

	deployFileName := path.Base(cfg.DeployFile)
	deployFilePath := path.Join(cfg.ProjectPath, deployFileName)
	checkDeployFileCmd := fmt.Sprintf(`test -f "%s"`, deployFilePath)

	logs.Verbosef("Checking for deploy file in project path: %s", deployFileName)
	logs.VerboseCommandf("%s", checkDeployFileCmd)

	if _, _, err := cli.RunCommandBuffered(checkDeployFileCmd); err != nil {
		logs.Warnf("Deploy file not found, skipping backup: %s", deployFilePath)
		return
	}
	logs.Success("Deploy file found - proceeding with backup...")

	timestamp := time.Now().Format("20060102_150405")
	backupDir := path.Join(cfg.ProjectPath, ".backup_"+timestamp)

	logs.Verbosef("Backup directory: %s", backupDir)

	mkdirCmd := fmt.Sprintf(`mkdir -p "%s"`, backupDir)
	if _, stderr, err := cli.RunCommandBuffered(mkdirCmd); err != nil {
		logs.Fatalf("Failed to create backup directory: %v\nDetails: %s", err, stderr)
	}

	backupCmd := fmt.Sprintf(`rsync -a --exclude "%s" "%s/" "%s/"`, path.Base(backupDir), cfg.ProjectPath, backupDir)
	logs.VerboseCommandf("%s", backupCmd)

	if _, stderr, err := cli.RunCommandBuffered(backupCmd); err != nil {
		logs.Fatalf("Failed to back up project directory: %v\nDetails: %s", err, stderr)
	} else {
		logs.Successf("Project directory backed up successfully at: %s", backupDir)
	}
}

func RestoreBackup(cli *client.Client, projectPath string) error {
	logs.Step("\U0001F4BE Restoring latest backup...")

	logs.Verbose("Locating latest backup directory...")
	findBackupCmd := `ls -td .backup_* 2>/dev/null | head -n 1`
	logs.VerboseCommandf("%s", findBackupCmd)

	backupDir, _, err := cli.RunCommandBuffered(fmt.Sprintf(`cd "%s" && %s`, projectPath, findBackupCmd))
	backupDir = strings.TrimSpace(backupDir)
	if err != nil || backupDir == "" {
		msg := fmt.Sprintf("no backup found in %s", projectPath)
		logs.Error(msg)
		return fmt.Errorf("%s", msg)
	}

	logs.Substepf("\U0001F4C2 Restoring from backup: %s", backupDir)
	restoreCmd := fmt.Sprintf(`cp -r "%s"/* "%s"/`, path.Join(projectPath, backupDir), projectPath)
	logs.VerboseCommandf("%s", restoreCmd)

	if _, stderr, err := cli.RunCommandBuffered(restoreCmd); err != nil {
		msg := fmt.Sprintf("failed to restore backup: %s", stderr)
		logs.Error(msg)
		return fmt.Errorf("%s", msg)
	}

	logs.Success("Backup restored successfully")
	return nil
}
