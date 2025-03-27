package files

import (
	"fmt"
	"log"
	"path"
	"path/filepath"

	"github.com/alcharra/docker-deploy-action-go/ssh"
)

func UploadFiles(client *ssh.Client, remoteDir string, files []string) {
	fmt.Printf("📂 Uploading files to %s:\n", remoteDir)

	for _, file := range files {
		remotePath := path.Join(remoteDir, filepath.Base(file))

		fmt.Printf("📦 Uploading: %s → %s\n", file, remotePath)
		err := client.UploadFileSCP(file, remotePath)
		if err != nil {
			log.Fatalf("❌ Failed to upload %s: %v", file, err)
		}
		fmt.Printf("✅ Uploaded: %s → %s\n", file, remotePath)
	}
}
