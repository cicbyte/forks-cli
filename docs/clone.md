# forks-cli clone

通过 Forks 镜像加速克隆 Git 仓库。

## 语法

```bash
forks-cli clone [flags] <仓库地址> [目标目录]
```

## 仓库地址格式

支持三种格式：

| 格式 | 示例 |
|------|------|
| 简写（推荐） | `author/repo` 或 `github/author/repo` |
| 原始 URL | `https://github.com/author/repo` |
| 镜像 URL | `http://host:port/git/github/author/repo.git` |

使用简写或原始 URL 时，需先配置镜像服务器：

```bash
forks-cli config set server http://192.168.1.100:8080
```

## 选项

| 选项 | 说明 |
|------|------|
| `-t, --token` | 本次使用的 API Token（不保存） |
| `-s, --server` | 本次使用的镜像服务器（不保存） |
| `-f, --force` | 强制更新镜像缓存 |

## 示例

```bash
# 简写克隆
forks-cli clone golang/go

# 指定目标目录
forks-cli clone golang/go ./my-go

# 临时指定服务器
forks-cli clone -s http://192.168.1.100:8080 golang/go

# 强制刷新镜像缓存
forks-cli clone -f golang/go
```

## Token 优先级

```
--token 参数 > FORKS_TOKEN 环境变量 > 配置文件
```
