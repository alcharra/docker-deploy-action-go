//go:build unit
// +build unit

package ssh

import (
	"os"
	"path"
	"path/filepath"
	"testing"
)

func TestUploadFileSCP_Positive(t *testing.T) {
	cfg := getTestConfig(t)

	client, err := NewClient(cfg)
	if err != nil {
		t.Fatalf("failed to create SSH client: %v", err)
	}
	defer client.Close()

	tmpFile, err := os.CreateTemp("", "test_upload_*.txt")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	tmpFile.WriteString("Hello, SCP test!")
	tmpFile.Close()

	remoteDir := "/tmp/test_upload_scp"
	remotePath := path.Join(remoteDir, filepath.Base(tmpFile.Name()))

	_, _, err = client.RunCommandBuffered("mkdir -p " + remoteDir)
	if err != nil {
		t.Fatalf("failed to create remote directory: %v", err)
	}

	err = client.UploadFileSCP(tmpFile.Name(), remotePath)
	if err != nil {
		t.Fatalf("UploadFileSCP failed: %v", err)
	}
}

func TestUploadFileSCP_InvalidFile(t *testing.T) {
	cfg := getTestConfig(t)

	client, err := NewClient(cfg)
	if err != nil {
		t.Fatalf("failed to create SSH client: %v", err)
	}
	defer client.Close()

	err = client.UploadFileSCP("this_file_does_not_exist.txt", "/tmp/should_fail.txt")
	if err == nil {
		t.Error("expected error when uploading nonexistent file, got nil")
	}
}
