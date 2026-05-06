# forks-cli trending

浏览 GitHub Trending 趋势仓库。数据来自 Forks 服务端缓存，支持按编程语言、时间范围、自然语言筛选。

## 语法

```bash
forks-cli trending [flags]
```

## 前置条件

```bash
forks-cli config set server http://192.168.1.100:8080
```

## 选项

| 选项 | 说明 |
|------|------|
| `-l, --language` | 编程语言筛选（go/python/rust/...） |
| `-s, --since` | 时间范围：daily/weekly/monthly（默认 daily） |
| `-S, --spoken` | 自然语言筛选（zh/en） |
| `-d, --date` | 指定日期（格式 2026-05-06） |
| `--refresh` | 跳过缓存重新获取 |
| `-t, --token` | 本次使用的 API Token（不保存） |
| `--server` | 本次使用的服务端地址（不保存） |

## 示例

```bash
# 今日全部语言
forks-cli trending

# 指定语言
forks-cli trending -l go

# 指定语言和时间范围
forks-cli trending -l python -s weekly

# 中文趋势
forks-cli trending -S zh

# 查看历史数据
forks-cli trending -d 2026-05-04

# 刷新缓存
forks-cli trending --refresh
```

## 输出格式

支持全局 `--format` 选项：

```bash
# 表格（默认）
forks-cli trending

# JSON
forks-cli trending --format json

# JSONL（每行一条）
forks-cli trending --format jsonl
```

### table 输出示例

```
今日 GitHub Trending (go) — 2026-05-06
  #  REPO                DESCRIPTION                          LANG  STARS  TODAY
  1  golang/go           The Go programming language          Go    128k   +42
  2  stretchr/testify    A toolkit with common assertions     Go    23.5k  +18
```

### JSON 输出示例

```json
{
  "date": "2026-05-06",
  "count": 25,
  "items": [
    {
      "author": "golang",
      "repo": "go",
      "url": "https://github.com/golang/go",
      "description": "The Go programming language",
      "language": "Go",
      "stars": 128000,
      "forks": 18000,
      "current_period_stars": 42
    }
  ]
}
```
