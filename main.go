package main

import (
	"fmt"

	"github.com/alcharra/docker-deploy-action-go/config"
	"github.com/alcharra/docker-deploy-action-go/deploy"
	"github.com/alcharra/docker-deploy-action-go/deploy/docker"
	"github.com/alcharra/docker-deploy-action-go/deploy/files"
	"github.com/alcharra/docker-deploy-action-go/utils"
)

func main() {
	cfg := config.LoadConfig()

	client := deploy.ConnectToSSH(cfg)
	defer client.Close()

	files.CheckOrCreateRemotePath(client, cfg)
	localFiles := utils.CollectFiles(cfg)

	files.BackupDeploymentFiles(client, cfg)
	files.UploadFiles(client, cfg.ProjectPath, localFiles)
	files.CheckFilesExistRemote(client, cfg.ProjectPath, localFiles)

	docker.CheckDockerInstalled(client)
	docker.EnsureDockerNetwork(client, cfg)
	docker.DockerRegistryLogin(client, cfg)
	docker.DeployDockerStack(client, cfg)
	docker.DeployDockerCompose(client, cfg)
	docker.RunDockerPrune(client, cfg)

	fmt.Println("✅ Deployment complete")
}
