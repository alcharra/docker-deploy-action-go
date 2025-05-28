package deploy

import (
	"fmt"

	"github.com/alcharra/docker-deploy-action-go/config"
	"github.com/alcharra/docker-deploy-action-go/internal/logs"
	"github.com/alcharra/docker-deploy-action-go/internal/ssh/client"
)

func Cleanup(client *client.Client, cfg config.DeployConfig) {
	logs.IsVerbose = cfg.Verbose

	if cfg.Mode != "compose" || !cfg.EnableRollback {
		return
	}

	logs.Step("\U0001F9FC Post-deployment cleanup started...")

	cleanupCmd := fmt.Sprintf(`find "%s" -maxdepth 1 -type d -name ".backup_*" -exec rm -rf {} +`, cfg.ProjectPath)

	logs.Verbose("Removing all backup directories")
	logs.VerboseCommandf("%s", cleanupCmd)

	if _, stderr, err := client.RunCommandBuffered(cleanupCmd); err != nil {
		logs.Warnf("Failed to clean up backup directories: %v\nDetails: %s", err, stderr)
		return
	}

	logs.Success("Backup directories cleaned up successfully")
}
