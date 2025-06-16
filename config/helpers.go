package config

import (
	"os"
	"strings"
)

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}

func getBool(key string, fallback bool) bool {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}
	return val == "true"
}

func splitEnv(key string) []string {
	val := os.Getenv(key)
	if val == "" {
		return []string{}
	}

	lines := strings.Split(val, "\n")
	var cleaned []string
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed != "" {
			cleaned = append(cleaned, trimmed)
		}
	}
	return cleaned
}

func ParseExtraFilesFromEnv(key string) []ExtraFile {
	val := os.Getenv(key)
	if val == "" {
		return nil
	}

	lines := strings.Split(val, "\n")
	var files []ExtraFile
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		flatten := false
		if strings.HasPrefix(line, "flatten ") {
			flatten = true
			line = strings.TrimPrefix(line, "flatten ")
		}

		parts := strings.SplitN(line, ":", 2)
		src := ""
		dst := ""

		if len(parts) == 2 {
			src = strings.TrimSpace(parts[0])
			dst = strings.TrimSpace(parts[1])
		} else {
			src = strings.TrimSpace(line)
		}

		files = append(files, ExtraFile{
			Src:     src,
			Dst:     dst,
			Flatten: flatten,
		})
	}
	return files
}
