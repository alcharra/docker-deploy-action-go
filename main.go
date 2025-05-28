package main

import (
	"github.com/alcharra/docker-deploy-action-go/config"
	"github.com/alcharra/docker-deploy-action-go/internal/deploy"
	"github.com/alcharra/docker-deploy-action-go/internal/docker"
	"github.com/alcharra/docker-deploy-action-go/internal/files"
	"github.com/alcharra/docker-deploy-action-go/internal/logs"
)

func main() {
	logs.Step("\U0001F680 Starting deployment...")

	cfg := config.LoadConfig()
	client := deploy.ConnectToSSH(cfg)
	defer client.Close()

	files.BackupDeploymentFiles(client, cfg)
	uploadedFiles := files.UploadFiles(client, cfg)
	files.CheckFilesExistRemote(client, cfg, uploadedFiles)

	docker.CheckDockerRequirements(client, &cfg)
	docker.EnsureDockerNetwork(client, cfg)
	docker.DockerRegistryLogin(client, cfg)

	docker.DeployDockerStack(client, cfg)
	docker.DeployDockerCompose(client, cfg)

	docker.RunDockerPrune(client, cfg)
	deploy.Cleanup(client, cfg)

	logs.Step("\U0001F389 All done — deployment completed successfully")
}
