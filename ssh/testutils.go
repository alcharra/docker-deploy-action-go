package ssh

import (
	"os"
	"testing"
)

func getSSHEnv(t *testing.T) (host, port, user, key string) {
	host = os.Getenv("SSH_HOST")
	port = os.Getenv("SSH_PORT")
	user = os.Getenv("SSH_USER")
	key = os.Getenv("SSH_KEY")

	if host == "" || user == "" || key == "" {
		t.Skip("Skipping SSH test: SSH_HOST, SSH_USER, or SSH_KEY not set")
	}
	if port == "" {
		port = "22"
	}
	return
}
