package deploy

import (
	"github.com/alcharra/docker-deploy-action-go/config"
	"github.com/alcharra/docker-deploy-action-go/internal/logs"
	"github.com/alcharra/docker-deploy-action-go/internal/ssh/client"
)

func ConnectToSSH(cfg config.DeployConfig) *client.Client {
	logs.Step("\U0001F50C Connecting to remote server...")
	logs.Substepf("\u2022 Host: %s", cfg.SSHHost)
	logs.Substepf("\u2022 User: %s", cfg.SSHUser)

	cli, err := client.NewClient(cfg)
	if err != nil {
		logs.Fatalf("Unable to establish SSH connection: %v", err)
	}

	logs.Success("SSH connection established")
	return cli
}
