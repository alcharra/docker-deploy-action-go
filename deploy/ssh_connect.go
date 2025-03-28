package deploy

import (
	"fmt"
	"log"

	"github.com/alcharra/docker-deploy-action-go/config"
	"github.com/alcharra/docker-deploy-action-go/ssh"
)

func ConnectToSSH(cfg config.DeployConfig) *ssh.Client {
	fmt.Println("🚀 Connecting to remote server...")
	client, err := ssh.NewClient(cfg)
	if err != nil {
		log.Fatalf("❌ SSH connection failed: %v\n", err)
	}
	fmt.Println("✅ SSH connection established.")
	return client
}
