package ssh

import (
	"os"

	"golang.org/x/crypto/ssh"
)

func parsePrivateKey(key string) (ssh.Signer, error) {
	if _, err := os.Stat(key); err == nil {
		keyBytes, err := os.ReadFile(key)
		if err != nil {
			return nil, err
		}
		return ssh.ParsePrivateKey(keyBytes)
	}

	return ssh.ParsePrivateKey([]byte(key))
}
