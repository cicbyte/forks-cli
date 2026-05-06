# backup 命令详细参考

## 工作流程

1. 从 Forks 服务端 `GET /api/repos` 获取仓库列表（自动分页）
2. 并发执行备份任务
3. 已存在目录：`git pull --ff-only`
4. 不存在目录：先尝试 Forks 镜像 clone，失败回退原始 URL clone

## 异常处理

| 场景 | 行为 |
|------|------|
| `.git` 存在且有效 | `git pull --ff-only` |
| `.git` 存在但仓库损坏 | 删除目录后重新 clone |
| 目录存在但无 `.git` | 删除目录后 clone |
| Forks 镜像 clone 失败 | 回退到 GitHub 原始 URL |
| 所有 clone 方式都失败 | 标记为失败，继续下一个 |

## 选项

| Flag | 说明 |
|------|------|
| `-d, --dir` | 备份目标目录（必填，或 config 设置 backup_dir） |
| `-c, --concurrency` | 并发 worker 数（默认 5） |
| `-t, --token` | 本次请求 Token |
| `-s, --server` | 本次使用服务器 |

## 目录结构

备份后的目录结构：

```
<backup_dir>/
└── github/
    ├── golang/
    │   └── go/
    ├── torvalds/
    │   └── linux/
    └── microsoft/
        └── vscode/
```

按 `source/author/repo` 三级目录组织。
