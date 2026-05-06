---
name: forks-cli
description: 指导 AI 使用 forks-cli 命令行工具完成 Git 仓库操作。当用户要求克隆仓库、批量备份仓库、浏览 GitHub Trending、打开 Forks Web UI 时触发。也适用于用户提到 "forks-cli"、"镜像克隆"、"备份仓库"、"trending" 等关键词时。注意：config 命令涉及敏感配置修改，仅供用户自行执行，AI 不应主动调用。
---

# forks-cli Skill

forks-cli 是 [Forks](https://github.com/cicbyte/forks) 服务端的命令行工具，提供镜像加速克隆、批量备份、GitHub Trending 浏览等功能。

## 重要约束

**不要执行 `config` 相关命令。** `config set/get/list` 涉及服务端地址、Token 等敏感配置，应由用户自行操作。如果命令报错提示未配置 server 或 token，告知用户执行对应的 config 命令即可，不要代为执行。

## 命令速查

| 命令 | 用途 | AI 可执行 |
|------|------|-----------|
| `clone` | 镜像加速克隆仓库 | 是 |
| `backup` | 批量备份仓库到本地 | 是 |
| `trending` | 浏览 GitHub Trending | 是 |
| `website` | 打开 Forks Web UI | 是 |
| `version` | 版本信息 | 是 |
| `config` | 管理配置 (list/get/set) | **否** |

## clone — 镜像加速克隆

```bash
forks-cli clone <仓库地址> [目标目录]
```

地址格式：`author/repo`（推荐）、`https://github.com/author/repo`、镜像 URL。

```bash
forks-cli clone golang/go                    # 简写
forks-cli clone golang/go ./my-go            # 指定目标目录
forks-cli clone -f golang/go                 # 强制刷新缓存
forks-cli clone -t <token> golang/go         # 临时 Token
```

详见 [references/clone.md](references/clone.md)

## backup — 批量备份

```bash
forks-cli backup -d <备份目录>               # 必须指定目录
forks-cli backup -d <目录> -c 10            # 自定义并发
```

已存在的仓库 `git pull --ff-only`，不存在的 `git clone`。损坏仓库自动清理重建。

详见 [references/backup.md](references/backup.md)

## trending — GitHub Trending

```bash
forks-cli trending                           # 今日全部
forks-cli trending -l go                     # Go 语言
forks-cli trending -l python -s weekly       # Python 周趋势
forks-cli trending -S zh                     # 中文趋势
forks-cli trending -d 2026-05-04            # 历史数据
forks-cli trending --format json            # JSON 输出
```

详见 [references/trending.md](references/trending.md)

## website — 打开 Web UI

```bash
forks-cli website
```

## config — 配置管理（仅供参考，不执行）

当命令因缺少配置报错时，提示用户：

```bash
# 未配置服务端地址时提示
请先配置服务端地址: forks-cli config set server <url>

# 未配置 Token 时提示
请先配置 Token: forks-cli config set token <your-token>
```

配置键：`server` / `token` / `backup_dir` / `log.level` / `log.max_size` / `log.max_backups` / `log.max_age` / `log.compress`

## 全局选项

`--format table|json|jsonl`（默认 table）

## Token/Server 优先级

```
命令行参数 > FORKS_TOKEN 环境变量 > 配置文件
```
