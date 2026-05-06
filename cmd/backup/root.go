package backup

import (
	"fmt"
	"os"
	"strings"

	"github.com/cicbyte/forks-cli/internal/common"
	logicbackup "github.com/cicbyte/forks-cli/internal/logic/backup"
	"github.com/spf13/cobra"
)

var (
	backupFlagToken       string
	backupFlagServer      string
	backupFlagDir         string
	backupFlagConcurrency int
)

// GetBackupCommand 返回 backup 命令
func GetBackupCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "backup [flags]",
		Short: "从 Forks 服务端批量备份仓库到本地",
		Long: `从 Forks 服务端批量备份仓库到本地

从 Forks 服务端读取已克隆的仓库列表，然后批量备份到本地。
已存在的仓库执行 git pull --ff-only，不存在的仓库执行 git clone。

使用前需先配置服务端地址:
  forks-cli config set server http://192.168.1.100:8080

Token 优先级: --token 参数 > FORKS_TOKEN 环境变量 > 配置文件`,
		Args: cobra.NoArgs,
		RunE: runBackup,
	}

	cmd.Flags().StringVarP(&backupFlagToken, "token", "t", "", "本次使用的 API Token（不保存）")
	cmd.Flags().StringVarP(&backupFlagServer, "server", "s", "", "本次使用的服务端地址（不保存）")
	cmd.Flags().StringVarP(&backupFlagDir, "dir", "d", "", "本地备份目录（默认 ./backup）")
	cmd.Flags().IntVarP(&backupFlagConcurrency, "concurrency", "c", 5, "并发数")

	return cmd
}

func runBackup(cmd *cobra.Command, args []string) error {
	cfg := common.AppConfigModel

	// 命令行 --server 临时覆盖
	server := cfg.Server
	if backupFlagServer != "" {
		server = strings.TrimSuffix(backupFlagServer, "/")
	}

	if server == "" {
		return fmt.Errorf("请先配置服务端地址: forks-cli config set server <url>")
	}

	// token 优先级: 命令行 > 环境变量 > 配置文件
	token := backupFlagToken
	if token == "" {
		token = os.Getenv("FORKS_TOKEN")
	}
	if token == "" {
		token = cfg.Token
	}

	// 备份目录优先级: 命令行 > 配置文件 > 默认值
	dir := backupFlagDir
	if dir == "" {
		dir = cfg.BackupDir
	}
	if dir == "" {
		dir = "./backup"
	}

	concurrency, _ := cmd.Flags().GetInt("concurrency")

	config := &logicbackup.BackupConfig{
		Server:      server,
		Token:       token,
		Dir:         dir,
		Concurrency: concurrency,
	}

	processor := logicbackup.NewBackupProcessor(config)
	return processor.Execute()
}
