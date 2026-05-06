# forks-cli config

管理 forks-cli 应用配置。

## 子命令

| 命令 | 说明 |
|------|------|
| `config list` | 列出所有配置项 |
| `config get <key>` | 查看单个配置项 |
| `config set <key> [value]` | 设置配置项 |

## config list

列出所有配置项及当前值。敏感字段（token）会显示脱敏后的值。

```bash
# 表格格式（默认）
forks-cli config list

# JSON 格式
forks-cli config list --format json

# JSONL 格式
forks-cli config list --format jsonl
```

## config get

查看指定配置项的当前值。敏感字段默认脱敏显示。

```bash
# 查看配置项
forks-cli config get server
forks-cli config get token
forks-cli config get log.level

# 查看敏感字段明文
forks-cli config get token --show
```

| 选项 | 说明 |
|------|------|
| `--show` | 显示敏感字段的明文值 |

## config set

设置配置项的值。敏感字段（token）如果不提供 value，会以不回显方式交互式输入。

```bash
# 设置服务端地址
forks-cli config set server http://192.168.1.100:8080

# 直接设置 Token
forks-cli config set token sk-xxx

# 交互式输入 Token（不回显）
forks-cli config set token

# 设置备份目录
forks-cli config set backup_dir /data/backup

# 设置日志级别
forks-cli config set log.level debug

# 设置布尔值
forks-cli config set log.compress false

# 设置数值
forks-cli config set log.max_size 20
```

## 配置项列表

| 键名 | 类型 | 说明 |
|------|------|------|
| `server` | string | Forks 服务端地址 |
| `token` | string | API Token（敏感字段） |
| `backup_dir` | string | 备份目录 |
| `log.level` | string | 日志级别（debug/info/warn/error） |
| `log.max_size` | int | 单个日志文件最大 MB |
| `log.max_backups` | int | 保留日志备份数 |
| `log.max_age` | int | 日志保留天数 |
| `log.compress` | bool | 是否压缩日志 |

## 配置文件

路径：`~/.cicbyte/forks-cli/config/config.yaml`

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
