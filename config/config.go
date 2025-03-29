package config

import (
	"os"
	"strings"
)

func LoadConfig() DeployConfig {
	return DeployConfig{
		SSHHost:             getEnv("SSH_HOST", ""),
		SSHPort:             getEnv("SSH_PORT", "22"),
		SSHUser:             getEnv("SSH_USER", ""),
		SSHKey:              getEnv("SSH_KEY", ""),
		SSHKeyPassphrase:    getEnv("SSH_KEY_PASSPHRASE", ""),
		SSHKnownHosts:       getEnv("SSH_KNOWN_HOSTS", ""),
		Fingerprint:         getEnv("FINGERPRINT", ""),
		Timeout:             getEnv("TIMEOUT", "10s"),
		ProjectPath:         getEnv("PROJECT_PATH", ""),
		DeployFile:          getEnv("DEPLOY_FILE", "docker-compose.yml"),
		ExtraFiles:          splitEnv("EXTRA_FILES"),
		Mode:                getEnv("MODE", "compose"),
		StackName:           getEnv("STACK_NAME", ""),
		DockerNetwork:       getEnv("DOCKER_NETWORK", ""),
		DockerNetworkDriver: getEnv("DOCKER_NETWORK_DRIVER", "bridge"),
		DockerNetworkAttach: getEnv("DOCKER_NETWORK_ATTACHABLE", "false") == "true",
		DockerPrune:         getEnv("DOCKER_PRUNE", "none"),
		RegistryHost:        getEnv("REGISTRY_HOST", ""),
		RegistryUser:        getEnv("REGISTRY_USER", ""),
		RegistryPass:        getEnv("REGISTRY_PASS", ""),
		EnableRollback:      getEnv("ENABLE_ROLLBACK", "false") == "true",
	}
}

func splitEnv(key string) []string {
	val := os.Getenv(key)

	if val == "" {
		return []string{}
	}

	return strings.Split(val, ",")
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
