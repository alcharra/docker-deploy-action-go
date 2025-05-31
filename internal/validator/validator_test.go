//go:build unit
// +build unit

package validator

import (
	"os"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestValidateComposeWithBindMountsAndEnv(t *testing.T) {
	yamlContent := `
services:
  web:
    image: nginx:latest
    volumes:
      - ./nginx.conf:/etc/nginx/nginx.conf:ro
    ports:
      - "8080:80"
    networks:
      - test_network_stack

  redis:
    image: redis:latest
    volumes:
      - ./redis.conf:/usr/local/etc/redis/redis.conf:ro
    command: ["redis-server", "/usr/local/etc/redis/redis.conf"]
    deploy:
      replicas: 1
      restart_policy:
        condition: on-failure
    networks:
      - test_network_stack

  db:
    image: postgres:latest
    environment:
      - POSTGRES_DB=mydb
      - POSTGRES_USER=admin
      - POSTGRES_PASSWORD=secret
    volumes:
      - db_data:/var/lib/postgresql/data
    networks:
      - test_network_stack

volumes:
  db_data: {}

networks:
  test_network_stack:
    driver: overlay
`

	os.Setenv("POSTGRES_DB", "mydb")
	os.Setenv("POSTGRES_USER", "admin")
	os.Setenv("POSTGRES_PASSWORD", "secret")

	content := os.ExpandEnv(yamlContent)

	var cfg ComposeFile
	if err := yaml.Unmarshal([]byte(content), &cfg); err != nil {
		t.Fatalf("Failed to parse YAML: %v", err)
	}

	err := cfg.Validate()
	if err != nil {
		t.Errorf("Validation failed unexpectedly:\n%s", err)
	}
}

func loadAndValidateYAML(t *testing.T, yamlContent string) error {
	content := os.ExpandEnv(yamlContent)

	var cfg ComposeFile
	if err := yaml.Unmarshal([]byte(content), &cfg); err != nil {
		t.Fatalf("Failed to parse YAML: %v", err)
	}

	return cfg.Validate()
}

func TestValidComposeFile(t *testing.T) {
	yaml := `
services:
  web:
    image: nginx
    volumes:
      - ./site:/usr/share/nginx/html
    ports:
      - "8080:80"
    networks:
      - net1
  db:
    image: postgres
    environment:
      - POSTGRES_PASSWORD=secret
    volumes:
      - pgdata:/var/lib/postgresql/data
    networks:
      - net1
volumes:
  pgdata: {}
networks:
  net1: {}
`
	if err := loadAndValidateYAML(t, yaml); err != nil {
		t.Errorf("Expected valid Compose file, got error: %v", err)
	}
}

func TestMissingImage(t *testing.T) {
	yaml := `
services:
  bad:
    ports:
      - "1234:80"
`
	err := loadAndValidateYAML(t, yaml)
	if err == nil || !strings.Contains(err.Error(), "missing 'image'") {
		t.Errorf("Expected error about missing image, got: %v", err)
	}
}

func TestUnsupportedBuildDirective(t *testing.T) {
	yaml := `
services:
  app:
    build: .
    image: myapp
`
	err := loadAndValidateYAML(t, yaml)
	if err == nil || !strings.Contains(err.Error(), "uses 'build'") {
		t.Errorf("Expected error for 'build', got: %v", err)
	}
}

func TestInvalidPortFormat(t *testing.T) {
	yaml := `
services:
  web:
    image: nginx
    ports:
      - "8080"
`
	err := loadAndValidateYAML(t, yaml)
	if err == nil || !strings.Contains(err.Error(), "must contain at least one ':'") {
		t.Errorf("Expected port format error, got: %v", err)
	}
}

func TestInvalidCommandListType(t *testing.T) {
	yaml := `
services:
  app:
    image: busybox
    command: ["echo", 123]
`
	err := loadAndValidateYAML(t, yaml)
	if err == nil || !strings.Contains(err.Error(), "command[1] must be a string") {
		t.Errorf("Expected command list string error, got: %v", err)
	}
}

func TestInvalidEntrypointListType(t *testing.T) {
	yaml := `
services:
  app:
    image: busybox
    entrypoint: ["run.sh", 42]
`
	err := loadAndValidateYAML(t, yaml)
	if err == nil || !strings.Contains(err.Error(), "entrypoint[1] must be a string") {
		t.Errorf("Expected entrypoint list string error, got: %v", err)
	}
}

func TestUndefinedVolumeReference(t *testing.T) {
	yaml := `
services:
  app:
    image: nginx
    volumes:
      - myvolume:/data
`
	err := loadAndValidateYAML(t, yaml)
	if err == nil || !strings.Contains(err.Error(), "uses undefined volume") {
		t.Errorf("Expected undefined volume error, got: %v", err)
	}
}

func TestUndefinedNetworkReference(t *testing.T) {
	yaml := `
services:
  app:
    image: nginx
    networks:
      - ghostnet
`
	err := loadAndValidateYAML(t, yaml)
	if err == nil || !strings.Contains(err.Error(), "uses undefined network") {
		t.Errorf("Expected undefined network error, got: %v", err)
	}
}

func TestInvalidDeployReplicas(t *testing.T) {
	yaml := `
services:
  app:
    image: nginx
    deploy:
      replicas: 0
`
	err := loadAndValidateYAML(t, yaml)
	if err == nil || !strings.Contains(err.Error(), "must be >= 1") {
		t.Errorf("Expected replicas >= 1 error, got: %v", err)
	}
}

func TestInvalidConstraintSyntax(t *testing.T) {
	yaml := `
services:
  app:
    image: nginx
    deploy:
      placement:
        constraints:
          - node.role=manager
`
	err := loadAndValidateYAML(t, yaml)
	if err == nil || !strings.Contains(err.Error(), "must contain '==' or '!='") {
		t.Errorf("Expected constraint syntax error, got: %v", err)
	}
}

func TestInvalidPortType(t *testing.T) {
	yaml := `
services:
  app:
    image: nginx
    ports:
      - 8080
`
	err := loadAndValidateYAML(t, yaml)
	if err == nil || !strings.Contains(err.Error(), "must be a string") {
		t.Errorf("Expected port type error, got: %v", err)
	}
}

func TestVolumeMapMissingSource(t *testing.T) {
	yaml := `
services:
  app:
    image: nginx
    volumes:
      - { target: /data }
volumes:
  vol: {}
`
	err := loadAndValidateYAML(t, yaml)
	if err == nil || !strings.Contains(err.Error(), "missing or invalid 'source'") {
		t.Errorf("Expected volume map missing source error, got: %v", err)
	}
}

func TestNetworkMapMissingName(t *testing.T) {
	yaml := `
services:
  app:
    image: nginx
    networks:
      - { driver: bridge }
networks:
  net1: {}
`
	err := loadAndValidateYAML(t, yaml)
	if err == nil || !strings.Contains(err.Error(), "missing or invalid 'name'") {
		t.Errorf("Expected network map missing name error, got: %v", err)
	}
}

func TestConfigReferenceMissingSource(t *testing.T) {
	yaml := `
services:
  app:
    image: nginx
    configs:
      - target: /app/config
configs:
  myconfig:
    file: ./config.yaml
`
	err := loadAndValidateYAML(t, yaml)
	if err == nil || !strings.Contains(err.Error(), "config with no 'source'") {
		t.Errorf("Expected config missing source error, got: %v", err)
	}
}

func TestUndefinedSecretReference(t *testing.T) {
	yaml := `
services:
  app:
    image: nginx
    secrets:
      - source: missing_secret
secrets:
  real_secret:
    file: ./secret.txt
`
	err := loadAndValidateYAML(t, yaml)
	if err == nil || !strings.Contains(err.Error(), "references undefined secret") {
		t.Errorf("Expected undefined secret error, got: %v", err)
	}
}

func TestInvalidExternalFieldType(t *testing.T) {
	yaml := `
services:
  app:
    image: nginx

volumes:
  badvol:
    external: "yes"
`
	err := loadAndValidateYAML(t, yaml)
	if err == nil || !strings.Contains(err.Error(), "invalid 'external' value") {
		t.Errorf("Expected invalid external value error, got: %v", err)
	}
}

func TestValidVolumeWithExternalMap(t *testing.T) {
	yaml := `
services:
  app:
    image: nginx
    volumes:
      - myvol:/data
volumes:
  myvol:
    driver: local
    external:
      name: external_vol
`
	err := loadAndValidateYAML(t, yaml)
	if err != nil {
		t.Errorf("Expected valid external volume map, got: %v", err)
	}
}

func TestEmptyServicesBlock(t *testing.T) {
	yaml := `
services: {}
`
	err := loadAndValidateYAML(t, yaml)
	if err == nil || !strings.Contains(err.Error(), "no services defined") {
		t.Errorf("Expected error about no services, got: %v", err)
	}
}

func TestSecretReferenceMissingSource(t *testing.T) {
	yaml := `
services:
  app:
    image: nginx
    secrets:
      - target: /run/secrets/password
secrets:
  mysecret:
    file: ./secret.txt
`
	err := loadAndValidateYAML(t, yaml)
	if err == nil || !strings.Contains(err.Error(), "secret with no 'source'") {
		t.Errorf("Expected secret missing source error, got: %v", err)
	}
}

func TestEnvironmentAsMap(t *testing.T) {
	yaml := `
services:
  app:
    image: postgres
    environment:
      POSTGRES_USER: admin
      POSTGRES_PASSWORD: secret
`
	err := loadAndValidateYAML(t, yaml)
	if err != nil {
		t.Errorf("Expected valid environment map, got: %v", err)
	}
}

func TestLabelsAsList(t *testing.T) {
	yaml := `
services:
  app:
    image: nginx
    labels:
      - "com.example.description=web app"
      - "com.example.version=1.0"
`
	err := loadAndValidateYAML(t, yaml)
	if err != nil {
		t.Errorf("Expected valid labels list, got: %v", err)
	}
}

func TestLabelsAsMap(t *testing.T) {
	yaml := `
services:
  app:
    image: nginx
    labels:
      com.example.description: "web app"
      com.example.version: "1.0"
`
	err := loadAndValidateYAML(t, yaml)
	if err != nil {
		t.Errorf("Expected valid labels map, got: %v", err)
	}
}

func TestIsBindMount(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"relative path", "./data:/path", true},
		{"absolute path", "/data:/path", true},
		{"parent dir", "../data:/path", true},
		{"home dir", "~/data:/path", true},
		{"named volume", "myvolume:/path", false},
		{"named volume with options", "data:/path:ro", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			volume := strings.SplitN(tt.input, ":", 2)[0]
			got := isBindMount(volume)
			if got != tt.expected {
				t.Errorf("isBindMount(%q) = %v, expected %v", volume, got, tt.expected)
			}
		})
	}
}
