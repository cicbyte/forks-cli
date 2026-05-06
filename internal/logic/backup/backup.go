package backup

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/cicbyte/forks-cli/internal/utils"
)

// BackupConfig backup 命令配置
type BackupConfig struct {
	Server      string
	Token       string
	Dir         string
	Concurrency int
}

// RepoItem 服务端返回的仓库信息
type RepoItem struct {
	Source   string `json:"source"`
	Author   string `json:"author"`
	Repo     string `json:"repo"`
	URL      string `json:"url"`
	IsCloned int    `json:"is_cloned"`
}

type apiResponse struct {
	Code    int      `json:"code"`
	Message string   `json:"message"`
	Data    pageData `json:"data"`
}

type pageData struct {
	List       []RepoItem `json:"list"`
	Total      int        `json:"total"`
	Page       int        `json:"page"`
	PageSize   int        `json:"page_size"`
	TotalPages int        `json:"total_pages"`
}

// BackupProcessor backup 处理器
type BackupProcessor struct {
	config *BackupConfig
}

func NewBackupProcessor(config *BackupConfig) *BackupProcessor {
	return &BackupProcessor{config: config}
}

// FetchRepoList 从服务端获取所有仓库列表（自动分页）
func FetchRepoList(server, token string) ([]RepoItem, error) {
	var allRepos []RepoItem
	page := 1
	pageSize := 100

	for {
		url := fmt.Sprintf("%s/api/repos?page=%d&page_size=%d", server, page, pageSize)

		client := &http.Client{Timeout: 30 * time.Second}
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return nil, fmt.Errorf("创建请求失败: %w", err)
		}
		if token != "" {
			req.Header.Set("Authorization", "Bearer "+token)
		}

		resp, err := client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("连接服务端失败: %w", err)
		}

		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			return nil, fmt.Errorf("读取响应失败: %w", err)
		}

		if resp.StatusCode == 401 {
			return nil, fmt.Errorf("认证失败，请设置 token: forks-cli config token <your-token>")
		}
		if resp.StatusCode != 200 {
			return nil, fmt.Errorf("服务端返回错误 (%d): %s", resp.StatusCode, string(body))
		}

		var result apiResponse
		if err := json.Unmarshal(body, &result); err != nil {
			return nil, fmt.Errorf("解析响应失败: %w", err)
		}

		pageRepos := result.Data.List
		if len(pageRepos) == 0 {
			break
		}
		allRepos = append(allRepos, pageRepos...)

		if len(pageRepos) < pageSize || page >= result.Data.TotalPages {
			break
		}
		page++
	}

	return allRepos, nil
}

// Execute 执行批量备份
func (p *BackupProcessor) Execute() error {
	server := p.config.Server
	absDir := utils.ResolveAbsDir(p.config.Dir)

	fmt.Printf("正在从 %s 获取仓库列表...\n", server)
	repos, err := FetchRepoList(server, p.config.Token)
	if err != nil {
		return err
	}
	fmt.Printf("获取到 %d 个仓库\n", len(repos))

	total := len(repos)
	if total == 0 {
		fmt.Println("没有需要备份的仓库")
		return nil
	}

	if err := os.MkdirAll(absDir, 0755); err != nil {
		return fmt.Errorf("创建备份目录失败: %w", err)
	}

	fmt.Printf("共 %d 个仓库，备份到 %s（并发 %d）\n\n", total, absDir, p.config.Concurrency)

	type backupTask struct {
		index  int
		repo   RepoItem
		server string
	}

	type backupResult struct {
		index  int
		action string // "cloned", "pulled", "failed"
		repo   string
		err    error
	}

	tasks := make(chan backupTask, total)
	results := make(chan backupResult, total)

	var wg sync.WaitGroup
	for w := 0; w < p.config.Concurrency; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for task := range tasks {
				repo := task.repo
				repoDir := filepath.Join(absDir, repo.Source, repo.Author, repo.Repo)
				gitDir := filepath.Join(repoDir, ".git")
				repoName := fmt.Sprintf("%s/%s/%s", repo.Source, repo.Author, repo.Repo)
				serverURL := fmt.Sprintf("%s/git/%s/%s/%s.git", task.server, repo.Source, repo.Author, repo.Repo)

				if _, err := os.Stat(gitDir); err == nil {
					// 已存在，pull
					if err := utils.RunGitInDirSilent(repoDir, "pull", "--ff-only"); err != nil {
						results <- backupResult{index: task.index, action: "failed", repo: repoName, err: fmt.Errorf("pull 失败: %v", err)}
						continue
					}
					results <- backupResult{index: task.index, action: "pulled", repo: repoName}
					continue
				}

				// 不存在，clone
				parentDir := filepath.Dir(repoDir)
				if err := os.MkdirAll(parentDir, 0755); err != nil {
					results <- backupResult{index: task.index, action: "failed", repo: repoName, err: fmt.Errorf("创建目录失败: %v", err)}
					continue
				}

				cloneErr := utils.RunGitNoProxy("clone", serverURL, repoDir)
				if cloneErr != nil {
					// 回退到原始 URL
					cloneErr = utils.RunGitSilent("clone", repo.URL, repoDir)
				}
				if cloneErr != nil {
					msg := strings.TrimSpace(cloneErr.Error())
					// 只保留最后几行
					lines := strings.Split(msg, "\n")
					if len(lines) > 3 {
						msg = strings.Join(lines[len(lines)-3:], "\n")
					}
					results <- backupResult{index: task.index, action: "failed", repo: repoName, err: fmt.Errorf("clone 失败: %s", msg)}
					continue
				}
				results <- backupResult{index: task.index, action: "cloned", repo: repoName}
			}
		}()
	}

	go func() {
		for i, repo := range repos {
			tasks <- backupTask{index: i, repo: repo, server: server}
		}
		close(tasks)
	}()

	go func() {
		wg.Wait()
		close(results)
	}()

	var cloned, pulled, failed int32
	done := int32(0)

	for r := range results {
		atomic.AddInt32(&done, 1)
		switch r.action {
		case "cloned":
			atomic.AddInt32(&cloned, 1)
			fmt.Printf("[%d/%d] %s ... \x1b[32m✓ 已克隆\x1b[0m\n", atomic.LoadInt32(&done), total, r.repo)
		case "pulled":
			atomic.AddInt32(&pulled, 1)
			fmt.Printf("[%d/%d] %s ... \x1b[32m✓ 已更新\x1b[0m\n", atomic.LoadInt32(&done), total, r.repo)
		case "failed":
			atomic.AddInt32(&failed, 1)
			fmt.Printf("[%d/%d] %s ... \x1b[31m✗ %v\x1b[0m\n", atomic.LoadInt32(&done), total, r.repo, r.err)
		}
	}

	fmt.Printf("\n完成: %d 已克隆, %d 已更新", cloned, pulled)
	if failed > 0 {
		fmt.Printf(", \x1b[31m%d 失败\x1b[0m", failed)
	}
	fmt.Println()

	return nil
}
