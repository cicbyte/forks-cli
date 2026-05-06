# trending 命令详细参考

## 数据来源

通过 Forks 服务端的 `/api/trending` 接口获取 GitHub Trending 数据。服务端负责爬取 GitHub 页面并缓存到本地 JSON 文件，CLI 只需调用 API。

## 选项

| Flag | 说明 |
|------|------|
| `-l, --language` | 编程语言：go, python, rust, typescript 等 |
| `-s, --since` | 时间范围：daily（默认）、weekly、monthly |
| `-S, --spoken` | 自然语言代码：zh（中文）、en（英文） |
| `-d, --date` | 指定历史日期，格式 2026-05-06 |
| `--refresh` | 跳过服务端缓存，重新爬取 |
| `-t, --token` | 本次请求 Token |
| `--server` | 本次使用服务器 |

## 输出格式

通过全局 `--format` 控制：

### table（默认）

```
今日 GitHub Trending (go) — 2026-05-06
  #  REPO              DESCRIPTION              LANG   STARS  TODAY
  1  golang/go         The Go programming...    Go     128k   +42
```

### json

```json
{
  "date": "2026-05-06",
  "count": 25,
  "items": [{ "author": "golang", "repo": "go", "stars": 128000, ... }]
}
```

### jsonl

每行一个仓库对象，适合管道处理：

```bash
forks-cli trending --format jsonl | jq '.stars > 10000'
```

## 常见场景

```bash
# 日常浏览
forks-cli trending

# 关注特定语言
forks-cli trending -l rust

# 周报统计用
forks-cli trending -s weekly --format json > weekly-trending.json

# 查看历史某天的数据
forks-cli trending -d 2026-05-01
```
