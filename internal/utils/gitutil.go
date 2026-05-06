package utils

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// RunGitSilent 执行 git 命令（静默），失败时返回 git stderr 内容
func RunGitSilent(args ...string) error {
	cmd := exec.Command("git", args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		msg := strings.TrimSpace(string(out))
		if msg != "" {
			return fmt.Errorf("%s", msg)
		}
		return err
	}
	return nil
}

// RunGitNoProxy 执行 git 命令，显式禁用代理
func RunGitNoProxy(args ...string) error {
	fullArgs := []string{"-c", "http.proxy=", "-c", "https.proxy="}
	fullArgs = append(fullArgs, args...)
	cmd := exec.Command("git", fullArgs...)
	cmd.Env = append(os.Environ(),
		"NO_PROXY=localhost,127.0.0.1,0.0.0.0,::1",
		"no_proxy=localhost,127.0.0.1,0.0.0.0,::1",
	)
	out, err := cmd.CombinedOutput()
	if err != nil {
		msg := strings.TrimSpace(string(out))
		if msg != "" {
			return fmt.Errorf("%s", msg)
		}
		return err
	}
	return nil
}

// RunGitInDirSilent 在指定目录执行 git 命令（静默），失败时返回 git stderr 内容
func RunGitInDirSilent(dir string, args ...string) error {
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		msg := strings.TrimSpace(string(out))
		if msg != "" {
			return fmt.Errorf("%s", msg)
		}
		return err
	}
	return nil
}

// ResolveAbsDir 将目录转为绝对路径
func ResolveAbsDir(dir string) string {
	if !filepath.IsAbs(dir) {
		absDir, _ := filepath.Abs(dir)
		return absDir
	}
	return dir
}
