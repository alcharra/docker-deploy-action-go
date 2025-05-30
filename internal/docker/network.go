package docker

import (
	"fmt"
	"strings"

	"github.com/alcharra/docker-deploy-action-go/config"
	"github.com/alcharra/docker-deploy-action-go/internal/logs"
	"github.com/alcharra/docker-deploy-action-go/internal/ssh/client"
)

func EnsureDockerNetwork(cli *client.Client, cfg config.DeployConfig) {
	logs.IsVerbose = cfg.Verbose
	network := cfg.DockerNetwork
	if network == "" {
		return
	}

	mode := cfg.Mode
	driver := cfg.DockerNetworkDriver
	if driver == "" {
		if mode == "stack" {
			driver = "overlay"
		} else {
			driver = "bridge"
		}
	}

	attachable := cfg.DockerNetworkAttach

	logs.Step("\U0001F310 Docker network checks...")

	existsCmd := fmt.Sprintf(`docker network inspect %s >/dev/null 2>&1 && echo EXISTS || echo MISSING`, network)
	logs.Verbosef("Checking if Docker network '%s' exists", network)
	logs.VerboseCommandf("docker network inspect %s >/dev/null", network)

	existsOut, _, err := cli.RunCommandBuffered(existsCmd)
	if err != nil {
		logs.Fatalf("Failed to check Docker network existence: %v", err)
	}

	switch strings.TrimSpace(existsOut) {
	case "EXISTS":
		logs.Successf("Network '%s' already exists", network)

		driverCmd := fmt.Sprintf("docker network inspect --format '{{ .Driver }}' %s", network)
		logs.Verbosef("Checking driver of network '%s'", network)
		logs.VerboseCommandf("%s", driverCmd)

		driverOut, _, err := cli.RunCommandBuffered(driverCmd)
		if err != nil {
			logs.Fatalf("Could not verify driver for network '%s': %v", network, err)
		}

		actual := strings.TrimSpace(driverOut)
		if actual != driver {
			logs.Warnf("Driver mismatch: found '%s', expected '%s'", actual, driver)
			logs.Info("Consider removing and recreating the network")
		} else {
			logs.Successf("Driver matches expected: '%s'", driver)
		}

	case "MISSING":
		logs.Infof("Network '%s' does not exist", network)
		logs.Substepf("\U0001F527 Creating network '%s' (driver: '%s')", network, driver)

		createCmd := fmt.Sprintf("docker network create --driver %s", driver)
		if driver == "overlay" && mode == "stack" {
			createCmd += " --scope swarm"
			if attachable {
				createCmd += " --attachable"
			}
		}
		createCmd += " " + network

		logs.VerboseCommandf("%s", createCmd)

		stdout, stderr, err := cli.RunCommandBuffered(createCmd)
		if err != nil {
			logs.Fatalf("Failed to create Docker network '%s': %v\nDetails: %s", network, err, stderr)
		}

		networkID := strings.TrimSpace(stdout)
		if networkID != "" {
			logs.Successf("Network '%s' created successfully (ID: %s)", network, networkID)
		} else {
			logs.Successf("Network '%s' created successfully", network)
		}

	default:
		logs.Fatalf("Unexpected output from network inspect: %s", existsOut)
	}
}
