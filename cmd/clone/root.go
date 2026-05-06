package clone

import (
	"os"
	"strings"

	"github.com/cicbyte/forks-cli/internal/common"
	logicclone "github.com/cicbyte/forks-cli/internal/logic/clone"
	"github.com/spf13/cobra"
)

var (
	cloneFlagToken  string
	cloneFlagServer string
	cloneFlagForce  bool
)

// GetCloneCommand 返回 clone 命令
func GetCloneCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "clone [flags] <仓库地址> [目标目录]",
		Short: "通过 Forks 镜像加速克隆 Git 仓库",
		Long: `通过 Forks 镜像加速克隆 Git 仓库

仓库地址支持三种格式:
  1) 镜像 URL:  http://host:port/git/github/author/repo.git
  2) 原始 URL:  https://github.com/author/repo
  3) 简写:      author/repo 或 github/author/repo

使用简写或原始 URL 时，需先配置镜像服务器:
  forks-cli config set server http://192.168.1.100:8080

Token 优先级: --token 参数 > FORKS_TOKEN 环境变量 > 配置文件`,
		Args: cobra.MinimumNArgs(1),
		RunE: runClone,
	}

	cmd.Flags().StringVarP(&cloneFlagToken, "token", "t", "", "本次使用的 API Token（不保存）")
	cmd.Flags().StringVarP(&cloneFlagServer, "server", "s", "", "本次使用的镜像服务器（不保存）")
	cmd.Flags().BoolVarP(&cloneFlagForce, "force", "f", false, "强制更新镜像缓存")

	return cmd
}

func runClone(cmd *cobra.Command, args []string) error {
	repoArg := args[0]
	targetDir := ""
	if len(args) > 1 {
		targetDir = args[1]
	}

	cfg := common.AppConfigModel

	// 命令行 --server 临时覆盖
	server := cfg.Server
	if cloneFlagServer != "" {
		server = strings.TrimSuffix(cloneFlagServer, "/")
	}

	// token 优先级: 命令行 > 环境变量 > 配置文件
	token := cloneFlagToken
	if token == "" {
		token = os.Getenv("FORKS_TOKEN")
	}
	if token == "" {
		token = cfg.Token
	}

	force, _ := cmd.Flags().GetBool("force")

	config := &logicclone.CloneConfig{
		Server:    server,
		Token:     token,
		RepoURL:   repoArg,
		TargetDir: targetDir,
		Force:     force,
	}

	processor := logicclone.NewCloneProcessor(config)
	return processor.Execute()
}
