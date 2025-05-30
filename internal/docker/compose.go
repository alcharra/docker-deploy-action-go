package docker

import (
	"fmt"
	"path"
	"strings"
	"time"

	"github.com/alcharra/docker-deploy-action-go/config"
	"github.com/alcharra/docker-deploy-action-go/internal/files"
	"github.com/alcharra/docker-deploy-action-go/internal/logs"
	"github.com/alcharra/docker-deploy-action-go/internal/ssh/client"
	"github.com/alcharra/docker-deploy-action-go/internal/utils"
)

func DeployDockerCompose(cli *client.Client, cfg config.DeployConfig) {
	logs.IsVerbose = cfg.Verbose
	if cfg.Mode != "compose" {
		return
	}
	if cfg.ComposeBinary == "" {
		logs.Fatalf("Compose binary not set. Ensure CheckDockerRequirements is called before deployment.")
	}

	compose := cfg.ComposeBinary
	composeFilePath := path.Join(cfg.ProjectPath, path.Base(cfg.DeployFile))

	if !cfg.RollbackTriggered {
		validateComposeConfig(cli, compose, composeFilePath)
	}

	if cfg.RollbackTriggered {
		logs.Step("\U0001F501 Re-deploying after rollback...")
	} else {
		logs.Step("\U0001F433 Deploying with Docker Compose...")
	}

	if cfg.ComposePull {
		pullImages(cli, compose, composeFilePath)
	} else if logs.IsVerbose {
		logs.Verbose("Skipping image pull as ComposePull is disabled")
	}

	stopServices(cli, compose, composeFilePath)
	if err := startServices(cli, compose, composeFilePath, buildComposeFlags(cfg)); err != nil {
		handleComposeFailure(cli, cfg, err.Error())
		return
	}

	logs.Substep("\U0001F433 Docker Compose deployment completed successfully")

	if err := checkServiceStatus(cli, compose, composeFilePath); err != nil {
		handleComposeFailure(cli, cfg, err.Error())
		return
	}
}

func validateComposeConfig(cli *client.Client, compose, filePath string) {
	logs.Step("\U0001F9EA Validating Docker Compose file...")
	logs.Verbosef("Compose file: %s", filePath)

	cmd := fmt.Sprintf(`%s -f "%s" config`, compose, filePath)
	logs.VerboseCommandf("%s", cmd)

	if _, stderr, err := cli.RunCommandBuffered(cmd); err != nil {
		cleaned := strings.ReplaceAll(strings.TrimSpace(stderr), "\n", " ")
		logs.Error("Compose file validation failed")
		logs.Fatalf("%s", cleaned)
	}

	logs.Success("Compose file is valid")
}

func pullImages(cli *client.Client, compose, filePath string) {
	logs.Verbose("Pulling latest images...")
	cmd := fmt.Sprintf(`%s -f "%s" pull`, compose, filePath)
	logs.VerboseCommandf("%s", cmd)
	if err := cli.RunCommandStreamed(cmd); err != nil {
		logs.Fatalf("Pull failed: %v", err)
	}
}

func stopServices(cli *client.Client, compose, filePath string) {
	logs.Verbose("Stopping existing services...")
	cmd := fmt.Sprintf(`%s -f "%s" down`, compose, filePath)
	logs.VerboseCommandf("%s", cmd)
	if err := cli.RunCommandStreamed(cmd); err != nil {
		logs.Fatalf("Failed to stop services: %v", err)
	}
}

func startServices(cli *client.Client, compose, filePath, flags string) error {
	logs.Verbose("Starting all services...")
	cmd := fmt.Sprintf(`%s -f "%s" up %s`, compose, filePath, flags)
	logs.VerboseCommandf("%s", cmd)
	return cli.RunCommandStreamed(cmd)
}

func buildComposeFlags(cfg config.DeployConfig) string {
	var flags []string
	flags = append(flags, "-d")
	if cfg.ComposeBuild {
		flags = append(flags, "--build")
	}
	if cfg.ComposeNoDeps {
		flags = append(flags, "--no-deps")
	}
	return strings.Join(flags, " ")
}

func checkServiceStatus(cli *client.Client, compose, filePath string) error {
	logs.Step("\U0001F50E Validating Docker Compose status...")
	logs.Verbose("Checking container status after deployment...")

	cmd := fmt.Sprintf(`%s -f "%s" ps`, compose, filePath)
	logs.VerboseCommandf("%s", cmd)

	time.Sleep(1 * time.Second)

	out, _, err := cli.RunCommandBuffered(cmd)
	if err != nil {
		return fmt.Errorf("failed to inspect services: %v", err)
	}

	lines := strings.Split(strings.TrimSpace(out), "\n")
	if len(lines) <= 1 {
		logs.Warn("No container lines found in `docker compose ps` output")
		return fmt.Errorf("no containers found to verify")
	}

	var failedContainers []string

	for _, line := range lines[1:] {
		fields := strings.Fields(line)
		if len(fields) < 5 {
			continue
		}

		status := strings.ToLower(strings.Join(fields[4:], " "))
		if strings.Contains(status, "exit") || strings.Contains(status, "dead") || strings.Contains(status, "restarting") {
			failedContainers = append(failedContainers, fmt.Sprintf("      \u2192 %s", line))
		}
	}

	if len(failedContainers) > 0 {
		logs.Substepf("\u2022 Container check failed for %d container%s", len(failedContainers), utils.Plural(len(failedContainers)))
		for _, msg := range failedContainers {
			fmt.Println(msg)
		}
		return fmt.Errorf("One or more containers failed to start")
	}

	logs.Success("All containers are running as expected")
	return nil
}

func handleComposeFailure(cli *client.Client, cfg config.DeployConfig, reason string) {
	logs.Error(reason)

	if cfg.EnableRollback && !cfg.RollbackTriggered {
		cfg.RollbackTriggered = true

		if err := files.RestoreBackup(cli, cfg.ProjectPath); err != nil {
			logs.Fatalf("Rollback failed — could not restore backup: %v", err)
		}

		DeployDockerCompose(cli, cfg)
		logs.Fatalf("Deployment failed — rollback completed successfully")
	}

	if cfg.RollbackTriggered {
		logs.Fatalf("Deployment failed — rollback attempted but still unsuccessful")
	}

	logs.Fatalf("Deployment failed")
}
