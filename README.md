# 🐳 Docker Deploy Action (Go)

[![Run Tests](https://github.com/alcharra/docker-deploy-action-go/actions/workflows/go-test.yml/badge.svg)](https://github.com/alcharra/docker-deploy-action-go/actions/workflows/go-test.yml)
[![GitHub tag](https://img.shields.io/github/tag/alcharra/docker-deploy-action-go.svg)](https://github.com/alcharra/docker-deploy-action-go/releases)
[![Go Version](https://img.shields.io/github/go-mod/go-version/alcharra/docker-deploy-action-go)](https://golang.org/)
[![Go Report Card](https://goreportcard.com/badge/github.com/alcharra/docker-deploy-action-go)](https://goreportcard.com/report/github.com/alcharra/docker-deploy-action-go)
[![CodeQL](https://github.com/alcharra/docker-deploy-action-go/actions/workflows/codeql-analysis.yml/badge.svg)](https://github.com/alcharra/docker-deploy-action-go/actions/workflows/codeql-analysis.yml)
[![Deploy Test](https://github.com/alcharra/docker-deploy-action-go/actions/workflows/deploy-test.yml/badge.svg)](https://github.com/alcharra/docker-deploy-action-go/actions/workflows/deploy-test.yml)
[![GoDoc](https://godoc.org/github.com/alcharra/docker-deploy-action-go?status.svg)](https://godoc.org/github.com/alcharra/docker-deploy-action-go)

A **reliable and efficient GitHub Action** written in Go for deploying **Docker Compose** and **Docker Swarm** services over SSH.

This action securely **uploads deployment files**, prepares the **remote environment** and automatically **provisions Docker networks** if needed. It supports **health checks, rollback** and optional **resource cleanup**, ensuring smooth and stable deployments with minimal hassle.

## Performance Comparison

The Go-based deployment tool was built with speed in mind — here’s a real-world comparison against the original [PowerShell/Bash-based version](https://github.com/alcharra/docker-deploy-action).

### Test Details

Both tools were tested under identical conditions: a Docker Compose deployment using the same configuration file along with three additional files (~1KB each). The tests were run on the same server, using the same SSH key, network and project path, ensuring a fair comparison between the two implementations.

### Results

| Tool            | Average Time | Fastest Time | Slowest Time |
| --------------- | ------------ | ------------ | ------------ |
| PowerShell/Bash | ~8.64s       | 8.38s        | 8.84s        |
| Go              | ~4.85s       | 4.82s        | 4.90s        |

✅ **Result:** The Go version is consistently **~44% faster** on average.

This speed gain comes from running a single compiled binary without shell overhead, resulting in faster deployments, especially in CI environments.

<details>
<summary>📸 See test outputs</summary>

| Go Version |                   Output                    |
| ---------- | :-----------------------------------------: |
| Test 1     | ![Go Test 1](./screenshots/go-deploy-1.png) |
| Test 2     | ![Go Test 2](./screenshots/go-deploy-2.png) |
| Test 3     | ![Go Test 3](./screenshots/go-deploy-3.png) |

| PowerShell/Bash |                       Output                        |
| --------------- | :-------------------------------------------------: |
| Test 1          | ![Script Test 1](./screenshots/script-deploy-1.png) |
| Test 2          | ![Script Test 2](./screenshots/script-deploy-2.png) |
| Test 3          | ![Script Test 3](./screenshots/script-deploy-3.png) |

</details>

## Inputs

| Input Parameter             | Description                                                                                          | Required | Default Value        |
| --------------------------- | ---------------------------------------------------------------------------------------------------- | :------: | -------------------- |
| `ssh_host`                  | Hostname or IP address of the target server                                                          |    ✅    |                      |
| `ssh_port`                  | Port used for the SSH connection                                                                     |    ❌    | `22`                 |
| `ssh_user`                  | Username used for the SSH connection                                                                 |    ✅    |                      |
| `ssh_key`                   | Private SSH key for authentication                                                                   |    ✅    |                      |
| `ssh_key_passphrase`        | Passphrase for the encrypted SSH private key                                                         |    ❌    |                      |
| `ssh_known_hosts`           | Contents of the SSH `known_hosts` file used to verify the server's identity                          |    ❌    |                      |
| `fingerprint`               | SSH host fingerprint for verifying the server's identity (SHA256 format)                             |    ❌    |                      |
| `timeout`                   | SSH connection timeout (e.g. `10s`, `30s`, `1m`)                                                     |    ❌    | `10s`                |
| `project_path`              | Path on the server where files will be uploaded                                                      |    ✅    |                      |
| `deploy_file`               | Path to the file used for defining the deployment (e.g. Docker Compose)                              |    ✅    | `docker-compose.yml` |
| `extra_files`               | Comma-separated list of additional files to upload (e.g. .env, config.yml)                           |    ❌    |                      |
| `mode`                      | Deployment mode (`compose` or `stack`)                                                               |    ❌    | `compose`            |
| `stack_name`                | Stack name used during Swarm deployment (required if mode is `stack`)                                |    ❌    |                      |
| `compose_pull`              | Whether to pull the latest images before bringing up services with Docker Compose (`true` / `false`) |    ❌    | `true`               |
| `compose_build`             | Whether to build images before starting services with Docker Compose (`true` / `false`)              |    ❌    | `false`              |
| `compose_no_deps`           | Whether to skip starting linked services (dependencies) with Docker Compose (`true` / `false`)       |    ❌    | `false`              |
| `compose_target_services`   | Comma-separated list of services to restart (e.g. web,db) - Restarts all if unset                    |    ❌    |                      |
| `docker_network`            | Name of the Docker network to be used or created if missing                                          |    ❌    |                      |
| `docker_network_driver`     | Driver for the network (`bridge`, `overlay`, `macvlan`, etc.)                                        |    ❌    | `bridge`             |
| `docker_network_attachable` | Whether standalone containers can attach to the network (`true` / `false`)                           |    ❌    | `false`              |
| `docker_prune`              | Type of Docker resource prune to run after deployment                                                |    ❌    | `none`               |
| `registry_host`             | Host address for the registry or remote service requiring authentication                             |    ❌    |                      |
| `registry_user`             | Username for authenticating with the registry or remote service                                      |    ❌    |                      |
| `registry_pass`             | Password or token for authenticating with the registry or remote service                             |    ❌    |                      |
| `enable_rollback`           | Whether to enable automatic rollback if deployment fails (`true` / `false`)                          |    ❌    | `false`              |

## SSH Host Key Verification

This tool supports two secure options for verifying the SSH server's identity:

- Using a `known_hosts` file (OpenSSH-compatible)
- Providing the server's SHA256 fingerprint

You only need to provide one of these options — not both.

> [!WARNING]  
> If neither `ssh_known_hosts` nor `fingerprint` is specified, the tool will fall back to `ssh.InsecureIgnoreHostKey()`.  
> This disables host key verification and leaves your connection vulnerable to man-in-the-middle attacks.  
> Never use this configuration in production environments.

> [!IMPORTANT]  
> For secure deployments, always provide either a `known_hosts` entry or a `fingerprint` to verify the server’s identity and prevent impersonation.

> [!TIP]  
> Use `ssh_known_hosts` for compatibility with OpenSSH and support for multiple key types.  
> Use `fingerprint` for a simpler, one-line setup in single-host environments.  
> In either case, store the value securely using a GitHub environment variable or secret.

## Supported Prune Types

- `none` - No pruning (default)
- `system` - Remove unused images, containers, volumes and networks
- `volumes` - Remove unused volumes
- `networks` - Remove unused networks
- `images` - Remove unused images
- `containers` - Remove stopped containers

## Network Management

This action ensures the required Docker network exists before deploying. If it is missing, it will be created automatically using the specified driver.

### How it works

- If the network already exists, its driver is verified.
- If the network does not exist, it is created using the provided driver.
- If `docker_network_attachable` is set to `true`, the network is created with the `--attachable` flag.
- In `stack` mode with the `overlay` driver:
  - Swarm mode must be active on the target server.
  - A warning is displayed if Swarm is not active.
- If the existing network uses a different driver than specified, a warning is displayed.

### Network scenarios

A network will be created if:

- The specified network does not exist.
- A custom network is defined via `docker_network`.
- The provided driver is valid and supported.

Warnings will be displayed if:

- The existing network's driver does not match the one specified.
- Swarm mode is inactive but `overlay` is requested in `stack` mode.

### Example usage

```yaml
docker_network: my_network
docker_network_driver: overlay
docker_network_attachable: true
```

## Rollback Behaviour

This action supports automatic rollback if a deployment fails to start correctly.

### How it works

- In `stack` mode:

  - Docker Swarm’s built-in rollback is used.
  - The command `docker service update --rollback <service-name>` is run to revert services in the stack to the last working state.

- In `compose` mode:
  - A backup of the current deployment file is created before deployment.
  - If services fail to start, the backup is restored and Compose is re-deployed.
  - If rollback is successful, the backup file is removed to avoid stale data.

### Rollback triggers

Rollback will occur if:

- Services fail health checks.
- Containers immediately exit after starting.
- Docker returns an error during service startup.

Rollback will not occur if:

- The deployment succeeds but the application has internal errors.
- A service is manually stopped by the user.
- Rollback is disabled via `enable_rollback: false`.

## Example Workflow

```yaml
name: Deploy

on:
  push:
    branches:
      - main

jobs:
  deploy:
    runs-on: ubuntu-latest

    steps:
      - name: 📦 Checkout repository
        uses: actions/checkout@v4

      # 🐳 Example 1: Deploy using Docker Stack (Swarm Mode)
      - name: 🚀 Deploy using Docker Stack
        uses: alcharra/docker-deploy-action-go@v1
        with:
          # SSH Connection
          ssh_host: ${{ secrets.SSH_HOST }} # Remote server IP or hostname
          ssh_user: ${{ secrets.SSH_USER }} # SSH username
          ssh_key: ${{ secrets.SSH_KEY }} # Private SSH key
          ssh_key_passphrase: ${{ secrets.SSH_KEY_PASSPHRASE }} # (Optional) SSH key passphrase
          ssh_known_hosts: ${{ secrets.SSH_KNOWN_HOSTS }} # (Optional) known_hosts entry

          # Deployment Settings
          project_path: /opt/myapp # Remote directory for upload and deploy
          deploy_file: docker-stack.yml # Stack file to deploy
          mode: stack # Deployment mode: 'stack'
          stack_name: myapp # Stack name on the target host

          # Additional Files
          extra_files: traefik.yml # Upload additional files (e.g. configs)

          # Docker Network Settings
          docker_network: myapp_network # Network name to use or create
          docker_network_driver: overlay # Network driver (e.g. bridge, overlay)

          # Post-Deployment Cleanup
          docker_prune: system # Prune unused Docker resources

          # Registry Authentication
          registry_host: ghcr.io
          registry_user: ${{ github.actor }}
          registry_pass: ${{ secrets.GITHUB_TOKEN }}

      # 🐳 Example 2: Deploy using Docker Compose
      - name: 🚀 Deploy using Docker Compose
        uses: alcharra/docker-deploy-action-go@v1
        with:
          # SSH Connection
          ssh_host: ${{ secrets.SSH_HOST }}
          ssh_user: ${{ secrets.SSH_USER }}
          ssh_key: ${{ secrets.SSH_KEY }}
          fingerprint: ${{ secrets.SSH_FINGERPRINT }} # (Optional) SHA256 host fingerprint

          # Deployment Settings
          project_path: /opt/myapp
          deploy_file: docker-compose.yml
          mode: compose

          # Additional Files
          extra_files: .env,database.env,nginx.conf # Upload environment and config files

          # Compose Behaviour
          compose_pull: true # Pull latest images before up
          compose_build: true # Build images before starting services
          compose_no_deps: true # Don’t start linked services
          compose_target_services: web,db # Restart only selected services (optional)

          # Rollback Support
          enable_rollback: true # Automatically rollback on failure

          # Docker Network
          docker_network: myapp_network
          docker_network_driver: bridge

          # Post-Deployment Cleanup
          docker_prune: volumes
```

## Requirements on the Server

- Docker must be installed
- Docker Compose (if using `compose` mode)
- Docker Swarm must be initialised (if using `stack` mode)
- SSH access must be configured for the provided user and key

## Important Notes

- This action is designed for Linux servers (Debian, Ubuntu, Alpine, CentOS)
- The SSH user must have permissions to write files and run Docker commands
- If the `project_path` does not exist, it will be created with permissions `750` and owned by the provided SSH user
- If using Swarm mode, the target machine must be a Swarm manager

## References

- [Docker Compose Documentation](https://docs.docker.com/compose/)
- [Docker Swarm Documentation](https://docs.docker.com/engine/swarm/)
- [Docker Prune Documentation](https://docs.docker.com/config/pruning/)
- [Docker Network Documentation](https://docs.docker.com/network/)

## Tips for Maintainers

- Test the full process locally before using in GitHub Actions
- Always use GitHub Secrets for sensitive values like SSH keys
- Make sure firewall rules allow SSH access from GitHub runners

## Contributing

Contributions are welcome. If you would like to improve this action, please feel free to open a pull request or raise an issue. I appreciate your input.

## Feature Requests

Have an idea or need something this action doesn't support yet?  
Please [start a discussion](https://github.com/alcharra/docker-deploy-action-go/discussions/new?category=ideas) under the **Ideas** category.

This helps keep feature requests organised and visible to others who may want the same thing.

## License

This project is licensed under the [MIT License](LICENSE).
