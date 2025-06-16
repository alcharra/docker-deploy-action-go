# 🐳 Docker Deploy Action (Go)

[![Run Tests](https://github.com/alcharra/docker-deploy-action-go/actions/workflows/go-test.yml/badge.svg)](https://github.com/alcharra/docker-deploy-action-go/actions/workflows/go-test.yml)
[![GitHub tag](https://img.shields.io/github/tag/alcharra/docker-deploy-action-go.svg)](https://github.com/alcharra/docker-deploy-action-go/releases)
[![Go Version](https://img.shields.io/github/go-mod/go-version/alcharra/docker-deploy-action-go)](https://golang.org/)
[![Go Report Card](https://goreportcard.com/badge/github.com/alcharra/docker-deploy-action-go)](https://goreportcard.com/report/github.com/alcharra/docker-deploy-action-go)
[![CodeQL](https://github.com/alcharra/docker-deploy-action-go/actions/workflows/codeql-analysis.yml/badge.svg)](https://github.com/alcharra/docker-deploy-action-go/actions/workflows/codeql-analysis.yml)
[![GoDoc](https://godoc.org/github.com/alcharra/docker-deploy-action-go?status.svg)](https://godoc.org/github.com/alcharra/docker-deploy-action-go)

A **fast and dependable GitHub Action** written in Go for deploying **Docker Compose** and **Docker Swarm** apps over SSH.

It handles everything from **file uploads** and **network setup** to **health checks**, **rollback** and **clean-up** — so your deployments stay simple, safe and consistent.

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

| Input Parameter             | Description                                                                             | Required | Default Value        |
| --------------------------- | --------------------------------------------------------------------------------------- | :------: | -------------------- |
| `ssh_host`                  | The hostname or IP address of the remote server you’re deploying to                     |    ✅    |                      |
| `ssh_port`                  | The port used to connect via SSH                                                        |    ❌    | `22`                 |
| `ssh_user`                  | The SSH username used to connect to the server                                          |    ✅    |                      |
| `ssh_key`                   | Your private SSH key for authenticating with the server                                 |    ✅    |                      |
| `ssh_key_passphrase`        | (If applicable) The passphrase used to unlock the SSH key                               |    ❌    |                      |
| `ssh_known_hosts`           | The contents of your `known_hosts` file, used to verify the server’s identity           |    ❌    |                      |
| `ssh_fingerprint`           | The server’s SSH fingerprint in SHA256 format (alternative to `known_hosts`)            |    ❌    |                      |
| `ssh_timeout`               | SSH connection timeout duration (e.g. `10s`, `30s`, `1m`)                               |    ❌    | `10s`                |
| `project_path`              | The full path on the server where files will be uploaded and deployed                   |    ✅    |                      |
| `deploy_file`               | The name of your main deployment file (e.g. `docker-compose.yml` or `docker-stack.yml`) |    ✅    | `docker-compose.yml` |
| `extra_files`               | A list of extra files or folders to upload. Use a multi-line format — one path per line |    ❌    |                      |
| `mode`                      | Deployment method: either `compose` or `stack`                                          |    ❌    | `compose`            |
| `stack_name`                | Name of the Docker stack (required if using `stack` mode)                               |    ❌    |                      |
| `compose_pull`              | Pull the latest images before starting services (`true` or `false`)                     |    ❌    | `true`               |
| `compose_build`             | Build images before starting services (`true` or `false`)                               |    ❌    | `false`              |
| `compose_no_deps`           | Skip starting linked services (`true` or `false`)                                       |    ❌    | `false`              |
| `compose_target_services`   | A list of specific services to restart. Use a multi-line format — one service per line  |    ❌    |                      |
| `docker_network`            | The name of the Docker network to use or create if missing                              |    ❌    |                      |
| `docker_network_driver`     | The network driver to use (`bridge`, `overlay`, etc.)                                   |    ❌    | `bridge`             |
| `docker_network_attachable` | Allow standalone containers to attach to the network (`true` or `false`)                |    ❌    | `false`              |
| `docker_prune`              | Type of Docker clean-up to run after deployment (e.g. `system`, `volumes`, `none`)      |    ❌    | `none`               |
| `registry_host`             | The container registry hostname (e.g. `ghcr.io`) if login is required                   |    ❌    |                      |
| `registry_user`             | Username for the registry                                                               |    ❌    |                      |
| `registry_pass`             | Password or token for the registry                                                      |    ❌    |                      |
| `enable_rollback`           | Automatically roll back if deployment fails (`true` or `false`)                         |    ❌    | `false`              |
| `env_vars`                  | Environment variables to include in a `.env` file uploaded to the server                |    ❌    |                      |
| `verbose`                   | Show extra internal command details and debug output (`true` or `false`)                |    ❌    | `false`              |

## SSH Host Key Verification

To securely verify the identity of your SSH server, you can use **either** of the following:

- A `known_hosts` entry (compatible with OpenSSH)
- A SHA256 `fingerprint` of the server's host key

You only need to provide **one** — not both.

> [!IMPORTANT]  
> If neither `ssh_known_hosts` nor `fingerprint` is set, the tool disables host key verification.  
> This exposes your connection to man-in-the-middle attacks and is **not safe for production**.  
> Always use one of the verification options and store it securely as a GitHub secret.

For most setups:

- Use `known_hosts` if you're familiar with SSH or need compatibility with multiple key types.
- Use `fingerprint` for a simpler, one-line setup — ideal for single-server use.

## Supported Prune Types

You can choose what to clean up on the server after deployment by setting the `docker_prune` option. The following types are supported:

- `none` – No pruning (default)
- `system` – Remove unused images, containers, volumes and networks
- `volumes` – Remove unused volumes
- `networks` – Remove unused networks
- `images` – Remove unused images
- `containers` – Remove stopped containers

## Controlling Upload Paths

By default, any file or folder listed under `extra_files` will be uploaded to the server **with its folder structure preserved**. For example, if you include `configs/settings.conf`, it will be uploaded to `project_path/configs/settings.conf`.

If you want to upload a file **directly to the root of the deployment folder**, ignoring its original directory, you can use the `flatten` keyword before the path.

### How It Works

- Files are uploaded with their **relative path preserved** by default.
- Prefix a path with `flatten` to upload file(s) **directly into the root** of the `project_path`, removing any folder structure.
- Use `source:destination` syntax to specify a custom destination path.
- Entire directories are supported and uploaded **recursively**, preserving structure.
- Supports both **individual files** and **glob patterns** such as `folder/*.env` or `assets/**/*`.

> [!NOTE]
> If multiple files flatten to the same name, the action will throw an error to prevent overwriting. Ensure flattened filenames are unique.

### Examples

```yaml
extra_files: |
  .env.production                        # → project-root/.env.production (preserved)
  configs/*                              # → project-root/configs/*.*
  flatten configs/*                      # → project-root/*.*
  flatten configs/db.env                 # → project-root/db.env
  configs/**/*.conf                      # → project-root/configs/**/*.conf
  flatten configs/**/*.conf              # → project-root/**/*.conf
  flatten configs/legacy.conf            # → project-root/legacy.conf
  scripts/init.sh                        # → project-root/scripts/init.sh
  flatten scripts/init.sh                # → project-root/init.sh
  assets/**/*                            # → project-root/assets/**/* (preserved structure)
  assets/:resources/                     # → project-root/resources/**/* (preserved under custom target)
  flatten assets/images/*.png:img/       # → project-root/img/*.png (flattened into folder)
```

### Best Practice

- Use `flatten` **only when necessary** to remove folder structure.
- Avoid flattening entire directories unless you're confident there are no filename conflicts.
- Default to preserved paths to ensure clarity and maintainability in your deployment layout.

## Docker Network Management

This step ensures that the required Docker network exists before deployment begins. If it does not exist, it will be created automatically using the specified driver and relevant options.

### How It Works

- If `docker_network` is not set, this step is skipped.
- If the network already exists, its driver is verified.
  - A warning is displayed if the driver does not match the expected value.
- If the network does not exist, it is created using:
  - A default driver if none is provided:
    - `overlay` in `stack` mode
    - `bridge` in all other modes
  - Optional flags:
    - `--attachable` if `docker_network_attachable: true`
    - `--scope swarm` when using `overlay` in `stack` mode

> [!TIP]  
> You do not need to specify the driver manually unless you want to override the defaults.

### Example

```yaml
docker_network: my_network
docker_network_attachable: true
mode: stack
# docker_network_driver: overlay  # Optional; defaults to 'overlay' in stack mode
```

## Rollback Behaviour

If something goes wrong during deployment, this action can automatically roll back to a previous working state.

### How It Works

- **Compose mode**  
  Before deployment, the full project folder is backed up.  
  If containers fail to start, the backup is restored and deployment is retried automatically.

- **Stack mode**  
  If any services fail to start or scale correctly, the tool attempts to roll back only the affected services using  
  `docker service update --rollback`.

> [!NOTE]  
> Rollback only runs if `enable_rollback` is set to `true`.  
> If rollback is attempted but fails, the process stops with an error message.

### When Rollback Happens

- Containers fail to start correctly in Compose mode
- Services in the stack fail to reach their expected replica count

### When It Doesn’t

Rollback will not trigger if:

- The deployment fails, but the previous version is still running (e.g. the new one never started)
- The app starts but has internal issues (e.g. logic errors or misconfiguration)
- Services are stopped or altered manually outside of deployment
- `enable_rollback` is set to `false`

### Example

```yaml
enable_rollback: true
mode: compose
```

or:

```yaml
enable_rollback: true
mode: stack
```

## YAML Validation (Beta)

This action now includes built-in validation for your Docker stack YAML file before deployment. It helps catch mistakes early and gives clear, readable feedback.

> [!WARNING]  
> YAML validation is a new feature and currently in **beta**. It may not catch every edge case or match `docker stack deploy` exactly, but it's actively being improved.

### What It Checks

- The `version` field is present and supported
- Each service defines an `image`
- `build` is not used (not supported in `docker stack deploy`)
- `deploy.replicas` is a positive number
- Port mappings use the correct `HOST:CONTAINER` format
- `command` values are valid (a string or list of strings)
- Placement constraints use valid syntax (like `node.role == manager`)
- Referenced `configs` and `secrets` exist at the top level
- `volumes`, `networks`, `configs` and `secrets` are defined as maps
- Duplicate keys are flagged to avoid unexpected behaviour

> [!NOTE]  
> This validation only applies when using **stack** mode.  
> If you're using **compose** mode, Docker's built-in `docker compose config` already performs deep validation automatically before deployment.

### Why It Matters

This check helps prevent common problems that could break your deployment, like missing fields, misconfigured services or invalid syntax. It makes issues easier to spot and fix before they cause failures.

> [!NOTE]  
> More checks and improvements will be added over time. If you find a validation error that seems wrong, feel free to open an issue.

## Example Workflows

### 🚀 Deploy Using Docker Stack

```yaml
name: Deploy Stack

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

      - name: 🚀 Deploy using Docker Stack
        uses: alcharra/docker-deploy-action-go@v2
        with:
          ssh_host: ${{ secrets.SSH_HOST }}
          ssh_user: ${{ secrets.SSH_USER }}
          ssh_key: ${{ secrets.SSH_KEY }}
          ssh_key_passphrase: ${{ secrets.SSH_KEY_PASSPHRASE }}
          ssh_known_hosts: ${{ secrets.SSH_KNOWN_HOSTS }}

          project_path: /opt/myapp
          deploy_file: docker-stack.yml
          mode: stack
          stack_name: myapp

          extra_files: |
            traefik.yml

          docker_network: myapp_network
          docker_network_driver: overlay

          docker_prune: system

          registry_host: ghcr.io
          registry_user: ${{ github.actor }}
          registry_pass: ${{ secrets.GITHUB_TOKEN }}
```

### 🐳 Deploy Using Docker Compose

```yaml
name: Deploy Compose

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

      - name: 🚀 Deploy using Docker Compose
        uses: alcharra/docker-deploy-action-go@v2
        with:
          ssh_host: ${{ secrets.SSH_HOST }}
          ssh_user: ${{ secrets.SSH_USER }}
          ssh_key: ${{ secrets.SSH_KEY }}
          ssh_fingerprint: ${{ secrets.SSH_FINGERPRINT }}

          project_path: /opt/myapp
          deploy_file: docker-compose.yml
          mode: compose

          env_vars: |
            DB_HOST=localhost
            DB_USER=myuser
            DB_PASS=${{ secrets.DB_PASS }}

          extra_files: |
            database.env
            nginx.conf

          compose_pull: true
          compose_build: true
          compose_no_deps: true
          compose_target_services: |
            web
            db

          enable_rollback: true

          docker_network: myapp_network
          docker_network_driver: bridge

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
