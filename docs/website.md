# forks-cli website

在默认浏览器中打开 Forks 服务端的 Web 管理界面。

## 语法

```bash
forks-cli website [flags]
```

## 选项

| 选项 | 说明 |
|------|------|
| `--server` | 本次使用的服务端地址（不保存） |

## 示例

```bash
# 打开已配置的服务端
forks-cli website

# 临时指定地址
forks-cli website --server http://192.168.1.100:8080
```

## 前置条件

```bash
forks-cli config set server http://192.168.1.100:8080
```

支持 Windows、macOS、Linux 自动打开默认浏览器。
