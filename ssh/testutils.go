package ssh

import (
	"os"
	"testing"

	"github.com/alcharra/docker-deploy-action-go/config"
)

func getTestConfig(t *testing.T) config.DeployConfig {
	host := os.Getenv("SSH_HOST")
	port := os.Getenv("SSH_PORT")
	user := os.Getenv("SSH_USER")
	key := os.Getenv("SSH_KEY")
	pass := os.Getenv("SSH_KEY_PASSPHRASE")

	if host == "" || user == "" || key == "" {
		t.Skip("Skipping SSH test: SSH_HOST, SSH_USER, or SSH_KEY not set")
	}
	if port == "" {
		port = "22"
	}

	return config.DeployConfig{
		SSHHost:          host,
		SSHPort:          port,
		SSHUser:          user,
		SSHKey:           key,
		SSHKeyPassphrase: pass,
	}
}
