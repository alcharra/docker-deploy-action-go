package utils

import (
	"log"
	"os"

	"github.com/alcharra/docker-deploy-action-go/config"
)

func CollectFiles(cfg config.DeployConfig) []string {
	files := []string{cfg.DeployFile}

	for _, file := range cfg.ExtraFiles {
		info, err := os.Stat(file)
		if err != nil {
			if os.IsNotExist(err) {
				log.Fatalf("❌ Extra file %s not found", file)
			}
			log.Fatalf("❌ Failed to stat file %s: %v", file, err)
		}

		if info.IsDir() {
			log.Fatalf("❌ %s is a directory, not a file", file)
		}

		files = append(files, file)
	}

	return files
}
