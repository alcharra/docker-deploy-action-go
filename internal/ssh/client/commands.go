package client

import (
	"bytes"
	"fmt"
	"strings"
)

func (cli *Client) RunCommandBuffered(cmd string) (string, string, error) {
	session, err := cli.sshClient.NewSession()
	if err != nil {
		return "", "", fmt.Errorf("failed to create SSH session: %w", err)
	}
	defer session.Close()

	var stdout, stderr bytes.Buffer
	session.Stdout = &stdout
	session.Stderr = &stderr

	err = session.Run(cmd)
	return stdout.String(), stderr.String(), err
}

func (cli *Client) RunCommandStreamed(cmd string) error {
	session, err := cli.sshClient.NewSession()
	if err != nil {
		return fmt.Errorf("failed to create SSH session: %w", err)
	}
	defer session.Close()

	stdout, err := session.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to get stdout pipe: %w", err)
	}
	stderr, err := session.StderrPipe()
	if err != nil {
		return fmt.Errorf("failed to get stderr pipe: %w", err)
	}

	if err := session.Start(cmd); err != nil {
		return fmt.Errorf("failed to start remote command: %w", err)
	}

	if strings.Contains(cmd, "docker compose") {
		go streamComposeOutput(stdout)
		go streamComposeOutput(stderr)
	} else {
		go streamStackOutput(stdout)
		go streamStackOutput(stderr)
	}

	if err := session.Wait(); err != nil {
		return fmt.Errorf("remote command failed: %w", err)
	}

	return nil
}
