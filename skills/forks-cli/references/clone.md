# clone 命令详细参考

## 地址解析规则

| 输入格式 | 示例 | 解析结果 |
|----------|------|----------|
| 简写 | `author/repo` | github.com/author/repo |
| 带平台简写 | `github/author/repo` | github.com/author/repo |
| 原始 URL | `https://github.com/author/repo` | 直接使用 |
| 镜像 URL | `http://host:port/git/github/author/repo.git` | 直接克隆 |

简写和原始 URL 会自动通过 Forks 镜像转换后克隆。

## 流程

1. 解析仓库地址 → 提取 source/author/repo
2. 调用 Forks `/api/git/prepare` 准备镜像
3. 从 Forks 镜像 URL 克隆到本地

## 选项

| Flag | 说明 |
|------|------|
| `-t, --token` | 本次请求 Token，不保存到配置 |
| `-s, --server` | 本次使用服务器，不保存到配置 |
| `-f, --force` | 强制 Forks 刷新镜像缓存 |

## 常见场景

```bash
# 首次使用，配置服务器后克隆
forks-cli config set server http://192.168.1.100:8080
forks-cli clone torvalds/linux

# 临时使用另一台 Forks 服务器
forks-cli clone -s http://10.0.0.1:8080 torvalds/linux

# 仓库较大，使用带 Token 的认证请求
forks-cli clone -t my-secret-token microsoft/vscode
```
