package clone

import (
	"fmt"
	"regexp"
	"strings"
)

// PlatformDomains 平台域名映射
var PlatformDomains = map[string]string{
	"github":    "github.com",
	"gitee":     "gitee.com",
	"gitlab":    "gitlab.com",
	"bitbucket": "bitbucket.org",
}

// RepoInfo 解析后的仓库信息
type RepoInfo struct {
	Source       string // github, gitee, gitlab, bitbucket
	Author       string
	Repo         string
	MirrorURL    string // 镜像克隆地址
	OriginalURL  string // 原始仓库地址
	ServerOrigin string // 镜像服务根地址
}

// ResolveRepoInfo 根据用户输入推算仓库信息
// 支持三种输入格式:
//   1. 完整镜像 URL: http://host:port/git/{source}/{author}/{repo}.git
//   2. 原始仓库 URL: https://github.com/author/repo[.git]
//   3. 简写: author/repo 或 source/author/repo
func ResolveRepoInfo(input, configuredServer string) (*RepoInfo, error) {
	if strings.Contains(input, "/git/") {
		s, a, r, server, e := parseMirrorURL(input)
		if e != nil {
			return nil, e
		}
		origURL := BuildOriginalURL(s, a, r)
		return &RepoInfo{
			Source: s, Author: a, Repo: r,
			MirrorURL: input, OriginalURL: origURL, ServerOrigin: server,
		}, nil
	}

	if strings.HasPrefix(input, "https://") || strings.HasPrefix(input, "http://") {
		s, a, r, e := parseOriginalURL(input)
		if e != nil {
			return nil, e
		}
		origURL := BuildOriginalURL(s, a, r)
		var mURL, sOrigin string
		if configuredServer != "" {
			mURL = BuildMirrorURL(configuredServer, s, a, r)
			sOrigin = strings.TrimSuffix(configuredServer, "/")
		}
		return &RepoInfo{
			Source: s, Author: a, Repo: r,
			MirrorURL: mURL, OriginalURL: origURL, ServerOrigin: sOrigin,
		}, nil
	}

	s, a, r, e := parseShorthand(input)
	if e != nil {
		return nil, e
	}
	origURL := BuildOriginalURL(s, a, r)
	var mURL, sOrigin string
	if configuredServer != "" {
		mURL = BuildMirrorURL(configuredServer, s, a, r)
		sOrigin = strings.TrimSuffix(configuredServer, "/")
	}
	return &RepoInfo{
		Source: s, Author: a, Repo: r,
		MirrorURL: mURL, OriginalURL: origURL, ServerOrigin: sOrigin,
	}, nil
}

func parseOriginalURL(rawURL string) (source, author, repo string, err error) {
	for s, domain := range PlatformDomains {
		prefix := "https://" + domain + "/"
		if strings.HasPrefix(rawURL, prefix) || strings.HasPrefix(rawURL, "http://"+domain+"/") {
			path := rawURL
			path = strings.TrimPrefix(path, "https://"+domain+"/")
			path = strings.TrimPrefix(path, "http://"+domain+"/")
			path = strings.TrimSuffix(path, ".git")
			path = strings.Trim(path, "/")
			parts := strings.SplitN(path, "/", 2)
			if len(parts) < 2 {
				return "", "", "", fmt.Errorf("无效的仓库 URL: %s", rawURL)
			}
			return s, parts[0], parts[1], nil
		}
	}
	return "", "", "", fmt.Errorf("不支持的平台 URL: %s\n支持: github.com, gitee.com, gitlab.com, bitbucket.org", rawURL)
}

func parseShorthand(input string) (source, author, repo string, err error) {
	input = strings.TrimSuffix(input, ".git")
	parts := strings.Split(input, "/")
	switch len(parts) {
	case 2:
		return "github", parts[0], parts[1], nil
	case 3:
		return parts[0], parts[1], parts[2], nil
	default:
		return "", "", "", fmt.Errorf("无效的仓库格式: %s\n应为: author/repo 或 source/author/repo", input)
	}
}

func parseMirrorURL(rawURL string) (source, author, repo, serverOrigin string, err error) {
	re := regexp.MustCompile(`/git/([^/]+)/([^/]+)/([^/?#]+?)(?:\.git)?(?:/|$|\?|#)`)
	matches := re.FindStringSubmatch(rawURL)
	if len(matches) < 4 {
		return "", "", "", "", fmt.Errorf("无法解析镜像 URL，格式应为: http://host:port/git/{source}/{author}/{repo}.git")
	}

	source = matches[1]
	author = matches[2]
	repo = matches[3]

	gitIdx := strings.Index(rawURL, "/git/")
	if gitIdx > 0 {
		serverOrigin = rawURL[:gitIdx]
	}

	return source, author, repo, serverOrigin, nil
}

// BuildMirrorURL 构建镜像克隆 URL
func BuildMirrorURL(server, source, author, repo string) string {
	server = strings.TrimSuffix(server, "/")
	return fmt.Sprintf("%s/git/%s/%s/%s.git", server, source, author, repo)
}

// BuildOriginalURL 根据 source/author/repo 构建原始仓库地址
func BuildOriginalURL(source, author, repo string) string {
	domain, ok := PlatformDomains[source]
	if !ok {
		domain = source + ".com"
	}
	return fmt.Sprintf("https://%s/%s/%s.git", domain, author, repo)
}
