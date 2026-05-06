package clone

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// PrepareResult 镜像服务 prepare API 返回结果
type PrepareResult struct {
	UseMirror bool
	Message   string
}

// PrepareFromServer 调用服务端 prepare API，返回是否可从镜像克隆
func PrepareFromServer(serverOrigin, token, source, author, repo string, force bool) PrepareResult {
	prepareURL := serverOrigin + "/api/git/prepare"

	body, _ := json.Marshal(map[string]interface{}{
		"source": source,
		"author": author,
		"repo":   repo,
		"force":  force,
	})

	client := &http.Client{Timeout: 300 * time.Second}
	req, err := http.NewRequest("POST", prepareURL, bytes.NewReader(body))
	if err != nil {
		return PrepareResult{Message: fmt.Sprintf("无法创建请求: %v", err)}
	}
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err := client.Do(req)
	if err != nil {
		return PrepareResult{Message: fmt.Sprintf("无法连接镜像服务: %v", err)}
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)

	if resp.StatusCode == 401 {
		return PrepareResult{Message: "认证失败，请设置 token: forks-cli config token <your-token>"}
	}

	if resp.StatusCode != 200 {
		msg := ""
		if m, ok := result["message"].(string); ok {
			msg = m
		}
		return PrepareResult{Message: fmt.Sprintf("镜像服务返回错误 (%d): %s", resp.StatusCode, msg)}
	}

	msg := ""
	if m, ok := result["message"].(string); ok {
		msg = m
	}
	return PrepareResult{UseMirror: true, Message: msg}
}
