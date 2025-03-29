package ssh

import (
	"bytes"
	"crypto/subtle"
	"fmt"
	"net"
	"os"
	"strings"
	"time"

	"github.com/alcharra/docker-deploy-action-go/config"
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
				return nil, fmt.Errorf("failed to decrypt SSH key with passphrase: %w", err)
			}
		} else {
			return nil, fmt.Errorf("failed to parse SSH key: %w", err)
		}
	}

	var hostKeyCallback ssh.HostKeyCallback

	switch {
	case cfg.SSHKnownHosts != "":
		tmpFile, err := os.CreateTemp("", "known_hosts")
		if err != nil {
			return nil, fmt.Errorf("failed to create temp known_hosts file: %w", err)
		}
		defer tmpFile.Close()

		if _, err := tmpFile.WriteString(cfg.SSHKnownHosts); err != nil {
			return nil, fmt.Errorf("failed to write to temp known_hosts file: %w", err)
		}

		hostKeyCallback, err = knownhosts.New(tmpFile.Name())
		if err != nil {
			return nil, fmt.Errorf("failed to parse known_hosts content: %w", err)
		}

	case cfg.Fingerprint != "":
		expected := strings.TrimSpace(cfg.Fingerprint)
		hostKeyCallback = func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			actual := ssh.FingerprintSHA256(key)
			if subtle.ConstantTimeCompare([]byte(actual), []byte(expected)) != 1 {
				return fmt.Errorf("host key mismatch: got %s, want %s", actual, expected)
			}
			return nil
		}

	default:
		fmt.Fprintln(os.Stderr, "WARNING: using insecure host key verification (InsecureIgnoreHostKey)")
		hostKeyCallback = ssh.InsecureIgnoreHostKey()
	}

	timeout := 10 * time.Second
	if cfg.Timeout != "" {
		if parsed, err := time.ParseDuration(cfg.Timeout); err == nil {
			timeout = parsed
		} else {
			return nil, fmt.Errorf("invalid timeout duration: %w", err)
		}
	}

	clientConfig := &ssh.ClientConfig{
		User: cfg.SSHUser,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
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
		client:     conn,
	}, nil
}

func (c *Client) RunCommandBuffered(cmd string) (string, string, error) {
	session, err := c.client.NewSession()
	if err != nil {
		return "", "", fmt.Errorf("failed to create session: %w", err)
	}
	defer session.Close()

	var stdout, stderr bytes.Buffer
	session.Stdout = &stdout
	session.Stderr = &stderr

	err = session.Run(cmd)
	return stdout.String(), stderr.String(), err
}

func (c *Client) RunCommandStreamed(cmd string) error {
	session, err := c.client.NewSession()
	if err != nil {
		return fmt.Errorf("failed to create session: %w", err)
	}
	defer session.Close()

	session.Stdout = os.Stdout
	session.Stderr = os.Stderr

	if err := session.Run(cmd); err != nil {
		return fmt.Errorf("command failed: %w", err)
	}

	return nil
}

func (c *Client) Close() error {
	return c.client.Close()
}
