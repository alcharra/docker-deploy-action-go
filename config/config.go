package config

func LoadConfig() DeployConfig {
	return DeployConfig{
		SSHHost:               getEnv("SSH_HOST", ""),
		SSHPort:               getEnv("SSH_PORT", "22"),
		SSHUser:               getEnv("SSH_USER", ""),
		SSHKey:                getEnv("SSH_KEY", ""),
		SSHKeyPassphrase:      getEnv("SSH_KEY_PASSPHRASE", ""),
		SSHKnownHosts:         getEnv("SSH_KNOWN_HOSTS", ""),
		SSHFingerprint:        getEnv("SSH_FINGERPRINT", ""),
		SSHTimeout:            getEnv("SSH_TIMEOUT", "10s"),
		ProjectPath:           getEnv("PROJECT_PATH", ""),
		DeployFile:            getEnv("DEPLOY_FILE", "docker-compose.yml"),
		ExtraFiles:            splitEnv("EXTRA_FILES"),
		Mode:                  getEnv("MODE", "compose"),
		StackName:             getEnv("STACK_NAME", ""),
		ComposePull:           getBool("COMPOSE_PULL", true),
		ComposeBuild:          getBool("COMPOSE_BUILD", false),
		ComposeNoDeps:         getBool("COMPOSE_NO_DEPS", false),
		ComposeTargetServices: splitEnv("COMPOSE_TARGET_SERVICES"),
		DockerNetwork:         getEnv("DOCKER_NETWORK", ""),
		DockerNetworkDriver:   getEnv("DOCKER_NETWORK_DRIVER", "bridge"),
		DockerNetworkAttach:   getBool("DOCKER_NETWORK_ATTACHABLE", false),
		DockerPrune:           getEnv("DOCKER_PRUNE", "none"),
		RegistryHost:          getEnv("REGISTRY_HOST", ""),
		RegistryUser:          getEnv("REGISTRY_USER", ""),
		RegistryPass:          getEnv("REGISTRY_PASS", ""),
		EnableRollback:        getBool("ENABLE_ROLLBACK", false),
		EnvVars:               getEnv("ENV_VARS", ""),
		Verbose:               getBool("VERBOSE", false),
	}
}
