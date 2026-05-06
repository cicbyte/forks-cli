# forks-cli backup

从 Forks 服务端批量备份仓库到本地。

已存在的仓库执行 `git pull --ff-only`，不存在的仓库执行 `git clone`。遇到损坏的仓库会自动清理后重新克隆。

## 语法

```bash
forks-cli backup [flags]
```

## 前置条件

```bash
forks-cli config set server http://192.168.1.100:8080
```

必须指定备份目录（通过 `-d` 参数或 `config set backup_dir`）。

## 选项

| 选项 | 说明 |
|------|------|
| `-d, --dir` | 本地备份目录（必填，或通过 config 设置） |
| `-c, --concurrency` | 并发数（默认 5） |
| `-t, --token` | 本次使用的 API Token（不保存） |
| `-s, --server` | 本次使用的服务端地址（不保存） |

## 示例

```bash
# 备份到指定目录
forks-cli backup -d /data/backup

# 使用配置文件中的路径
forks-cli config set backup_dir /data/backup
forks-cli backup

# 自定义并发数
forks-cli backup -d /data/backup -c 10

# 临时指定服务器
forks-cli backup -d /data/backup -s http://192.168.1.100:8080
```

## 输出示例

```
正在从 http://192.168.1.100:8080 获取仓库列表...
获取到 455 个仓库
共 455 个仓库，备份到 /data/backup（并发 5）

[1/455] github/golang/go ... ✓ 已更新
[2/455] github/torvalds/linux ... ✓ 已克隆
[3/455] github/some/repo ... ✗ clone 失败: ...

完成: 400 已克隆, 50 已更新, 5 失败
```

## Token 优先级

```
--token 参数 > FORKS_TOKEN 环境变量 > 配置文件
```

## 异常处理

| 场景 | 处理方式 |
|------|----------|
| 仓库已存在且有效 | `git pull --ff-only` |
| 仓库已存在但 pull 失败 | 清理目录后重新 clone |
| 目录存在但无 `.git` | 清理目录后 clone |
| 镜像 clone 失败 | 回退到原始 GitHub URL clone |
