package files

import (
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"

	"github.com/alcharra/docker-deploy-action-go/config"
	"github.com/alcharra/docker-deploy-action-go/ssh"
)

func UploadFiles(client *ssh.Client, cfg config.DeployConfig) []string {
	fmt.Printf("📂 Uploading files to remote directory: %s\n", cfg.ProjectPath)

	var filesToUpload []string
	filesToUpload = append(filesToUpload, cfg.DeployFile)

	for _, file := range cfg.ExtraFiles {
		info, err := os.Stat(file)
		if err != nil {
			if os.IsNotExist(err) {
				log.Fatalf("❌ Extra file '%s' not found", file)
			}
			log.Fatalf("❌ Failed to stat file '%s': %v", file, err)
		}

		if info.IsDir() {
			log.Fatalf("❌ '%s' is a directory, not a file", file)
		}

		filesToUpload = append(filesToUpload, file)
	}

	if cfg.EnvVars != "" {
		envFile := ".env"
		err := os.WriteFile(envFile, []byte(cfg.EnvVars), 0644)
		if err != nil {
			log.Fatalf("❌ Failed to create .env file: %v", err)
		}
		filesToUpload = append(filesToUpload, envFile)
		defer func() {
			err := os.Remove(envFile)
			if err != nil {
				log.Printf("⚠️  Failed to remove temporary .env file: %v", err)
			} else {
				fmt.Println("🧹 Removed temporary .env file after upload")
			}
		}()
		fmt.Println("🌿 Generated .env file from provided environment variables")
	}

	for _, file := range filesToUpload {
		remotePath := path.Join(cfg.ProjectPath, filepath.Base(file))

		fmt.Printf("📦 Uploading: %s → %s\n", file, remotePath)
		err := client.UploadFileSCP(file, remotePath)
		if err != nil {
			log.Fatalf("❌ Failed to upload '%s': %v", file, err)
		}
		fmt.Printf("✅ Successfully uploaded: %s → %s\n", file, remotePath)
	}

	return filesToUpload
}
