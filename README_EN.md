# forks-cli

**[English](#features)** | **[中文](README.md)**

> CLI tool for [Forks](https://github.com/cicbyte/forks) server — accelerated cloning, batch backup, and GitHub Trending browsing

> Forks 服务端命令行工具 — 镜像加速克隆、批量备份、GitHub Trending 浏览，配合 [Forks](https://github.com/cicbyte/forks) 服务端使用

## Features

| Command | Description |
|---------|-------------|
| `forks-cli clone` | Clone Git repos via Forks mirror acceleration |
| `forks-cli backup` | Batch backup repos from server to local |
| `forks-cli trending` | Browse GitHub Trending repos |
| `forks-cli config` | Manage app config (list/get/set) |
| `forks-cli website` | Open Forks Web UI in browser |
| `forks-cli version` | Show version info |

## Screenshots

### Trending

![trending](images/trending.png)

### Clone

![clone](images/clone.gif)

### Backup

![backup](images/backup.gif)

## Installation

Download from [Releases](https://github.com/cicbyte/forks-cli/releases) or build from source:

```bash
git clone https://github.com/cicbyte/forks-cli.git
cd forks-cli
go build -o forks-cli
```

Requires: Go 1.23+

## Quick Start

```bash
# Configure Forks server address
forks-cli config set server http://192.168.1.100:8080

# Configure API Token
forks-cli config set token <your-token>

# Clone a repo with mirror acceleration
forks-cli clone golang/go

# Browse today's trending repos
forks-cli trending

# Batch backup to a directory
forks-cli backup -d /data/backup
```

## Usage

### clone — Mirror-accelerated Clone

Supports three URL formats:

```bash
# Shorthand (recommended)
forks-cli clone author/repo

# Original URL
forks-cli clone https://github.com/author/repo

# Mirror URL
forks-cli clone http://host:port/git/github/author/repo.git
```

| Flag | Description |
|------|-------------|
| `-t, --token` | Token for this request (not saved) |
| `-s, --server` | Server address for this request |
| `-f, --force` | Force refresh mirror cache |

### backup — Batch Backup

Fetches repo list from Forks server and batch clones or updates to a local directory.

```bash
# Backup to specified directory (required)
forks-cli backup -d /data/backup

# Use path from config
forks-cli config set backup_dir /data/backup
forks-cli backup

# Custom concurrency
forks-cli backup -d /data/backup -c 10
```

| Flag | Description |
|------|-------------|
| `-d, --dir` | Backup directory (required, or set via config) |
| `-c, --concurrency` | Concurrency (default 5) |
| `-t, --token` | Token for this request |
| `-s, --server` | Server address for this request |

### trending — GitHub Trending

```bash
# All languages, today
forks-cli trending

# Filter by language and time range
forks-cli trending -l go -s weekly

# Filter by spoken language
forks-cli trending -S zh

# View historical data
forks-cli trending -d 2026-05-04

# JSON output
forks-cli trending --format json
```

| Flag | Description |
|------|-------------|
| `-l, --language` | Programming language (go/python/rust/...) |
| `-s, --since` | Time range: daily/weekly/monthly (default daily) |
| `-S, --spoken` | Spoken language (zh/en) |
| `-d, --date` | Specific date (2026-05-04) |
| `--refresh` | Skip cache and re-fetch |

### config — Configuration

```bash
# List all config
forks-cli config list

# Set config values
forks-cli config set server http://192.168.1.100:8080
forks-cli config set token              # Interactive input (hidden)
forks-cli config set log.level debug

# Get a single value
forks-cli config get server
forks-cli config get token --show       # Show plaintext
```

Config keys:

| Key | Description |
|-----|-------------|
| `server` | Forks server address |
| `token` | API Token (sensitive) |
| `backup_dir` | Backup directory |
| `log.level` | Log level (debug/info/warn/error) |
| `log.max_size` | Max log file size in MB |
| `log.max_backups` | Number of log backups to keep |
| `log.max_age` | Log retention days |
| `log.compress` | Compress log files |

### website — Open Web UI

```bash
forks-cli website
```

## Global Flags

| Flag | Description |
|------|-------------|
| `--format` | Output format: table/json/jsonl (default table) |

## Priority

Token and Server support multi-level priority:

```
CLI flags > Environment variables > Config file
```

Environment variable: `FORKS_TOKEN`

## Config File

Path: `~/.cicbyte/forks-cli/config/config.yaml`

```yaml
server: http://192.168.1.100:8080
token: your-api-token
backup_dir: /data/backup
log:
  level: info
  maxSize: 10
  maxBackups: 30
  maxAge: 30
  compress: true
```

## License

[MIT](LICENSE) © 2025 cicbyte
