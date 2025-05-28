package client

import (
	"crypto/subtle"
	"fmt"
	"net"
	"os"
	"strings"
	"time"

	"github.com/alcharra/docker-deploy-action-go/config"
	"github.com/alcharra/docker-deploy-action-go/internal/logs"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/knownhosts"
)

func NewClient(cfg config.DeployConfig) (*Client, error) {
	keyBytes := []byte(cfg.SSHKey)

	var signer ssh.Signer
	var err error

	signer, err = ssh.ParsePrivateKey(keyBytes)
	if err != nil {
		if _, ok := err.(*ssh.PassphraseMissingError); ok && cfg.SSHKeyPassphrase != "" {
			signer, err = ssh.ParsePrivateKeyWithPassphrase(keyBytes, []byte(cfg.SSHKeyPassphrase))
			if err != nil {
				return nil, fmt.Errorf("failed to decrypt SSH key using passphrase: %w", err)
			}
		} else {
			return nil, fmt.Errorf("failed to parse SSH private key: %w", err)
		}
	}

	var hostKeyCallback ssh.HostKeyCallback

	switch {
	case cfg.SSHKnownHosts != "":
		tmpFile, err := os.CreateTemp("", "known_hosts")
		if err != nil {
			return nil, fmt.Errorf("failed to create temporary known_hosts file: %w", err)
		}
		defer tmpFile.Close()

		if _, err := tmpFile.WriteString(cfg.SSHKnownHosts); err != nil {
			return nil, fmt.Errorf("failed to write contents to known_hosts file: %w", err)
		}

		hostKeyCallback, err = knownhosts.New(tmpFile.Name())
		if err != nil {
			return nil, fmt.Errorf("failed to parse known_hosts data: %w", err)
		}

	case cfg.SSHFingerprint != "":
		expected := strings.TrimSpace(cfg.SSHFingerprint)
		hostKeyCallback = func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			actual := ssh.FingerprintSHA256(key)
			if subtle.ConstantTimeCompare([]byte(actual), []byte(expected)) != 1 {
				return fmt.Errorf("SSH host key mismatch – got %s, expected %s", actual, expected)
			}
			return nil
		}

	default:
		logs.Warn("Host key verification is disabled (not recommended for production)")
		hostKeyCallback = ssh.InsecureIgnoreHostKey()
	}

	timeout := 10 * time.Second
	if cfg.SSHTimeout != "" {
		if parsed, err := time.ParseDuration(cfg.SSHTimeout); err == nil {
			timeout = parsed
		} else {
			return nil, fmt.Errorf("invalid SSH timeout duration: %w", err)
		}
	}

	clientConfig := &ssh.ClientConfig{
		User:            cfg.SSHUser,
		Auth:            []ssh.AuthMethod{ssh.PublicKeys(signer)},
		HostKeyCallback: hostKeyCallback,
		Timeout:         timeout,
	}

	conn, err := ssh.Dial("tcp", fmt.Sprintf("%s:%s", cfg.SSHHost, cfg.SSHPort), clientConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to SSH host: %w", err)
	}

	return &Client{
		Host:       cfg.SSHHost,
		Port:       cfg.SSHPort,
		User:       cfg.SSHUser,
		PrivateKey: cfg.SSHKey,
		sshClient:  conn,
	}, nil
}

func (cli *Client) NewSession() (*ssh.Session, error) {
	if cli.sshClient == nil {
		return nil, fmt.Errorf("SSH client is not initialised")
	}
	return cli.sshClient.NewSession()
}

func (cli *Client) Close() error {
	if cli.sshClient == nil {
		return nil
	}
	return cli.sshClient.Close()
}
