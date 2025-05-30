package docker

import (
	"fmt"
	"path"
	"strings"

	"github.com/alcharra/docker-deploy-action-go/config"
	"github.com/alcharra/docker-deploy-action-go/internal/logs"
	"github.com/alcharra/docker-deploy-action-go/internal/ssh/client"
	"github.com/alcharra/docker-deploy-action-go/internal/utils"
	"github.com/alcharra/docker-deploy-action-go/internal/validator"
)

func DeployDockerStack(cli *client.Client, cfg config.DeployConfig) {
	logs.IsVerbose = cfg.Verbose

	if cfg.Mode != "stack" {
		return
	}

	if err := validateStackFile(cfg); err != nil {
		logs.Errorf("%s", err)
		logs.Fatalf("Aborting deployment")
	}

	logs.Step("\u2693 Deploying Docker stack...")
	logs.Verbosef("Stack name: %s", cfg.StackName)

	if err := runStackDeployment(cli, cfg); err != nil {
		logs.Errorf("%s", err)
		handleDeploymentFailures(cli, cfg)
		return
	}

	logs.Substepf("\U0001F6A2 All services in Docker stack '%s' have converged successfully", cfg.StackName)

	if err := validateStackStatus(cli, cfg, false); err != nil {
		handleDeploymentFailures(cli, cfg)
	}
}

func validateStackFile(cfg config.DeployConfig) error {
	deployFilePath := cfg.DeployFile

	logs.Step("\U0001F9EA Validating Docker Stack file...")
	logs.Verbosef("Stack file: %s", deployFilePath)

	stackCfg, err := validator.LoadComposeFile(deployFilePath)
	if err != nil {
		return err
	}

	if err := stackCfg.Validate(); err != nil {
		return err
	}

	logs.Success("Stack file validation passed")
	return nil
}

func runStackDeployment(cli *client.Client, cfg config.DeployConfig) error {
	stackName := cfg.StackName
	deployFilePath := path.Join(cfg.ProjectPath, path.Base(cfg.DeployFile))

	if cfg.EnvVars != "" {
		logs.Substep("\U0001F4C4 Loading environment variables")
		logs.VerboseCommand("set -a")
		logs.VerboseCommandf(`source "%s/.env"`, cfg.ProjectPath)
		logs.VerboseCommand("set +a")
	}

	withAuth := ""
	if cfg.RegistryHost != "" && cfg.RegistryUser != "" && cfg.RegistryPass != "" {
		withAuth = "--with-registry-auth"
	}

	logs.Substepf("\U0001F4E6 Deploying stack '%s'", stackName)
	logs.VerboseCommandf(`docker stack deploy -c "%s" "%s" %s --detach=false`, deployFilePath, stackName, withAuth)

	cmd := fmt.Sprintf(`
		STACK="%s"
		PROJECT_PATH="%s"
		DEPLOY_FILE="%s"
		ENV_VARS='%s'
		WITH_AUTH="%s"

		if [ -f "$PROJECT_PATH/.env" ] && [ -n "$ENV_VARS" ]; then
			set -a
			source "$PROJECT_PATH/.env"
			set +a
		fi

		docker stack deploy -c "$DEPLOY_FILE" "$STACK" $WITH_AUTH --detach=false
	`, stackName, cfg.ProjectPath, deployFilePath, cfg.EnvVars, withAuth)

	return cli.RunCommandStreamed(cmd)
}

func validateStackStatus(cli *client.Client, cfg config.DeployConfig, afterDeployFailure bool) error {
	logs.Step("\U0001F50E Validating stack status...")
	logs.Verbosef("Validating status of stack '%s'...", cfg.StackName)

	cmd := fmt.Sprintf(`docker service ls --filter "label=com.docker.stack.namespace=%s"`, cfg.StackName)
	logs.VerboseCommand(cmd)

	output, _, err := cli.RunCommandBuffered(cmd)
	if err != nil {
		return fmt.Errorf("error verifying services: %w", err)
	}

	var failedServices []string

	for line := range strings.SplitSeq(strings.TrimSpace(output), "\n") {
		fields := strings.Fields(line)
		if len(fields) < 5 {
			continue
		}

		replicas := fields[3]
		name := fields[1]
		image := fields[4]

		if parts := strings.Split(replicas, "/"); len(parts) == 2 && parts[0] != parts[1] {
			failedServices = append(failedServices,
				fmt.Sprintf("      \u2192 %s — REPLICAS: %s, IMAGE: %s", name, replicas, image),
			)
		}
	}

	if len(failedServices) > 0 {
		logs.Substepf("\u2022 Health check failed for %d service%s", len(failedServices), utils.Plural(len(failedServices)))
		for _, msg := range failedServices {
			fmt.Println(msg)
		}
		logs.Errorf("Stack validation failed for '%s'", cfg.StackName)
		return fmt.Errorf("one or more services failed to start")
	}

	if afterDeployFailure {
		fmt.Printf("   \u2705 All services in stack '%s' are healthy %s(despite deployment error)%s\n",
			cfg.StackName,
			logs.GrayColor,
			logs.ResetColor,
		)
	} else {
		logs.Successf("All services in stack '%s' are healthy", cfg.StackName)
	}

	return nil
}

func handleDeploymentFailures(cli *client.Client, cfg config.DeployConfig) {
	if err := validateStackStatus(cli, cfg, true); err == nil {
		logs.Fatalf("Deployment failed")
	}

	if cfg.EnableRollback && !cfg.RollbackTriggered {
		cfg.RollbackTriggered = true
		logs.Step("\U0001F504 Starting rollback...")

		services := getServiceStatus(cli, cfg.StackName)
		if rollbackStack(cli, services) {
			logs.Fatalf("Deployment failed — rollback succeeded")
		}
	}

	logs.Fatalf("Deployment failed")
}

func getServiceStatus(cli *client.Client, stack string) []string {
	logs.Verbosef("Fetching service list for rollback in stack '%s'...", stack)

	cmd := fmt.Sprintf(`docker service ls --filter "label=com.docker.stack.namespace=%s" --format "{{.Name}} {{.Replicas}}"`, stack)
	logs.VerboseCommand(cmd)

	output, _, err := cli.RunCommandBuffered(cmd)
	if err != nil {
		logs.Warnf("Could not retrieve service list: %v", err)
		return nil
	}

	return strings.Split(strings.TrimSpace(output), "\n")
}

func rollbackStack(cli *client.Client, lines []string) bool {
	var rolledBack bool

	for _, line := range lines {
		parts := strings.Fields(line)
		if len(parts) != 2 {
			continue
		}

		name := parts[0]
		replicas := parts[1]

		replicaParts := strings.Split(replicas, "/")
		if len(replicaParts) == 2 && replicaParts[0] != replicaParts[1] {
			logs.Substepf("\U0001F501 Rolling back %s", name)
			cmd := fmt.Sprintf(`docker service update --rollback "%s"`, name)
			logs.VerboseCommand(cmd)

			if err := cli.RunCommandStreamed(cmd); err != nil {
				logs.Warnf("Rollback failed for %s", name)
			} else {
				logs.Successf("Rolled back: %s", name)
				rolledBack = true
			}
		} else {
			logs.Verbosef("No rollback needed for %s (replicas: %s)", name, replicas)
		}
	}

	return rolledBack
}
