package ssh

import (
	"strings"
	"testing"

	"github.com/alcharra/docker-deploy-action-go/config"
)

func TestNewClient_Valid(t *testing.T) {
	cfg := getTestConfig(t)

	client, err := NewClient(cfg)
	if err != nil {
		t.Fatalf("expected SSH client to connect, got error: %v", err)
	}
	defer client.Close()
}

func TestNewClient_InvalidKey(t *testing.T) {
	cfg := config.DeployConfig{
		SSHHost: "example.com",
		SSHPort: "22",
		SSHUser: "invalid",
		SSHKey:  "not_a_real_key",
	}

	_, err := NewClient(cfg)
	if err == nil {
		t.Fatal("expected error with invalid SSH key, got nil")
	}
}

func TestRunCommandBuffered(t *testing.T) {
	cfg := getTestConfig(t)

	client, err := NewClient(cfg)
	if err != nil {
		t.Fatalf("SSH client creation failed: %v", err)
	}
	defer client.Close()

	stdout, stderr, err := client.RunCommandBuffered("echo Hello")
	if err != nil {
		t.Fatalf("expected command to run successfully, got error: %v", err)
	}
	if !strings.Contains(stdout, "Hello") {
		t.Errorf("expected stdout to contain 'Hello', got: %s", stdout)
	}
	if stderr != "" {
		t.Errorf("expected no stderr, got: %s", stderr)
	}
}

func TestRunCommandBuffered_InvalidCommand(t *testing.T) {
	cfg := getTestConfig(t)

	client, err := NewClient(cfg)
	if err != nil {
		t.Fatalf("SSH client creation failed: %v", err)
	}
	defer client.Close()

	_, stderr, err := client.RunCommandBuffered("non_existing_command_1234")
	if err == nil {
		t.Fatal("expected error from invalid command, got nil")
	}
	if stderr == "" {
		t.Error("expected stderr message for invalid command, got empty string")
	}
}

func TestRunCommandStreamed(t *testing.T) {
	cfg := getTestConfig(t)

	client, err := NewClient(cfg)
	if err != nil {
		t.Fatalf("SSH client creation failed: %v", err)
	}
	defer client.Close()

	err = client.RunCommandStreamed("echo 'Streamed Output'")
	if err != nil {
		t.Fatalf("expected streamed command to run successfully, got: %v", err)
	}
}
