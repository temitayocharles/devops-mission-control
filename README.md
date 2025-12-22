# ops-tool

A lightweight, agnostic DevOps CLI tool unifying Kubernetes, Docker, cloud infrastructure, and more.

Think of it as **k9s for the entire DevOps ecosystem**.

## Features

### Phase 1 (MVP)
- ✅ **Kubernetes**: Pods, Deployments, Services, Logs, Context switching
- ✅ **Docker**: Containers, Images, Compose
- ✅ **AWS**: EC2, S3, RDS, Cost Optimization basics
- ✅ **Git**: Status, Branches, Commits, Push/Pull
- ✅ **Lightweight**: Single Go binary, ~25MB, minimal dependencies
- ✅ **Plugin Architecture**: Extensible for custom commands

### Planned Features
- **Phase 2**: GCP, Azure, Helm, Terraform, Network diagnostics
- **Phase 3**: Observability dashboard, Slack integration, Audit logging
- **Phase 4**: RBAC, CI/CD integrations, Multi-tenancy, DR workflows

## Installation

### From Source
```bash
git clone https://github.com/yourusername/ops-tool.git
cd ops-tool
go build -o ops-tool
sudo mv ops-tool /usr/local/bin/
```

### From Docker
```bash
docker pull yourusername/ops-tool:latest
alias ops-tool='docker run --rm -v ~/.kube:/.kube -v ~/.aws:/.aws yourusername/ops-tool'
```

## Usage

```bash
# Show help
ops-tool help

# Kubernetes operations
ops-tool k8s pods list
ops-tool k8s pods logs <pod-name>
ops-tool k8s deployments list
ops-tool k8s context list
ops-tool k8s context switch <context>

# Docker operations
ops-tool docker containers list
ops-tool docker containers stop <container-id>
ops-tool docker images list
ops-tool docker compose up
ops-tool docker compose down

# AWS operations
ops-tool aws ec2 list
ops-tool aws s3 list
ops-tool aws rds list
ops-tool aws cost optimize

# Git operations
ops-tool git status
ops-tool git branch list
ops-tool git branch create <branch-name>
ops-tool git commit "<message>"
ops-tool git push
ops-tool git pull
```

## Design Philosophy

- **Lightweight**: Single binary, fast startup, minimal memory footprint
- **Intuitive**: Agnostic command naming across all domains
- **Extensible**: Plugin system for custom commands
- **Safe**: Dry-run mode, confirmations on destructive operations
- **Audit-friendly**: Optional audit logging for compliance

## Project Structure

```
ops-tool/
├── main.go           # Entry point
├── cmd/              # Command implementations
│   ├── root.go       # Root command
│   ├── k8s.go        # Kubernetes commands
│   ├── docker.go     # Docker commands
│   ├── aws.go        # AWS commands
│   └── git.go        # Git commands
├── pkg/              # Shared packages (future)
│   ├── config/       # Configuration management
│   ├── cache/        # Caching layer
│   └── plugins/      # Plugin system
├── go.mod            # Go module file
├── Dockerfile        # Container image
├── .github/          # GitHub workflows
│   └── workflows/
│       ├── ci.yml    # Testing & building
│       └── release.yml # Release automation
├── README.md         # This file
├── LICENSE           # MIT License
└── docs/             # Documentation

```

## Development

### Prerequisites
- Go 1.21+
- Docker (for containerization)
- kubectl, docker, aws-cli (for testing)

### Setup
```bash
git clone https://github.com/yourusername/ops-tool.git
cd ops-tool
go mod download
go build
```

### Testing
```bash
go test ./...
```

### Building Docker Image
```bash
docker build -t yourusername/ops-tool:latest .
```

## Contributing

Contributions are welcome! Please:

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit changes (`git commit -m 'Add amazing feature'`)
4. Push to branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## Roadmap

- [ ] Phase 1: Core CLI framework (MVP)
- [ ] Phase 2: Advanced cloud features & multi-cloud
- [ ] Phase 3: Observability & monitoring dashboard
- [ ] Phase 4: Enterprise features (RBAC, audit, multi-tenancy)
- [ ] Community plugins marketplace

## License

MIT License - see LICENSE file for details

## Support

- 📖 Documentation: [docs/](./docs/)
- 🐛 Issues: [GitHub Issues](https://github.com/yourusername/ops-tool/issues)
- 💬 Discussions: [GitHub Discussions](https://github.com/yourusername/ops-tool/discussions)

---

Made with ❤️ for DevOps engineers
