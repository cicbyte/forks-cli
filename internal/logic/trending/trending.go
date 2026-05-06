package trending

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

// TrendingConfig trending 命令配置
type TrendingConfig struct {
	Server         string
	Token          string
	Language       string
	Since          string
	SpokenLanguage string
	Date           string
	Refresh        bool
}

// TrendingRepo 趋势仓库信息
type TrendingRepo struct {
	Author             string `json:"author"`
	Repo               string `json:"repo"`
	URL                string `json:"url"`
	Description        string `json:"description"`
	Language           string `json:"language"`
	LanguageColor      string `json:"language_color,omitempty"`
	Stars              int    `json:"stars"`
	Forks              int    `json:"forks"`
	CurrentPeriodStars int    `json:"current_period_stars"`
}

type trendingResponse struct {
	Code    int            `json:"code"`
	Message string         `json:"message"`
	Data    trendingData   `json:"data"`
}

type trendingData struct {
	Date  string         `json:"date"`
	Count int            `json:"count"`
	Items []TrendingRepo `json:"items"`
}

// FetchTrending 从 Forks 服务端获取 GitHub Trending 数据
func FetchTrending(config *TrendingConfig) ([]TrendingRepo, string, error) {
	params := url.Values{}
	if config.Language != "" {
		params.Set("language", config.Language)
	}
	if config.Since != "" {
		params.Set("since", config.Since)
	}
	if config.SpokenLanguage != "" {
		params.Set("spoken_language_code", config.SpokenLanguage)
	}
	if config.Date != "" {
		params.Set("date", config.Date)
	}
	if config.Refresh {
		params.Set("refresh", "true")
	}

	reqURL := fmt.Sprintf("%s/api/trending?%s", config.Server, params.Encode())

	client := &http.Client{Timeout: 30 * time.Second}
	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return nil, "", fmt.Errorf("创建请求失败: %w", err)
	}
	if config.Token != "" {
		req.Header.Set("Authorization", "Bearer "+config.Token)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, "", fmt.Errorf("连接服务端失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", fmt.Errorf("读取响应失败: %w", err)
	}

	if resp.StatusCode == 401 {
		return nil, "", fmt.Errorf("认证失败，请设置 token: forks-cli config set token <your-token>")
	}
	if resp.StatusCode != 200 {
		return nil, "", fmt.Errorf("服务端返回错误 (%d): %s", resp.StatusCode, string(body))
	}

	var result trendingResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, "", fmt.Errorf("解析响应失败: %w", err)
	}

	return result.Data.Items, result.Data.Date, nil
}
