package cmd

import (
	"fmt"
	"os"

	"github.com/cicbyte/forks-cli/cmd/backup"
	"github.com/cicbyte/forks-cli/cmd/clone"
	"github.com/cicbyte/forks-cli/cmd/configcmd"
	"github.com/cicbyte/forks-cli/cmd/trending"
	"github.com/cicbyte/forks-cli/internal/common"
	"github.com/cicbyte/forks-cli/internal/log"
	"github.com/cicbyte/forks-cli/internal/output"
	"github.com/cicbyte/forks-cli/internal/utils"
	"github.com/spf13/cobra"
)

var globalFormat string

var rootCmd = &cobra.Command{
	Use:   "forks-cli",
	Short: "forks-cli",
	Long:  `forks-cli`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// 初始化应用目录
	if err := utils.InitAppDirs(); err != nil {
		fmt.Printf("初始化目录失败: %v\n", err)
		os.Exit(1)
	}
	// 加载配置(会自动创建默认配置)
	common.AppConfigModel = utils.ConfigInstance.LoadConfig()

	// 全局 flag
	rootCmd.PersistentFlags().StringVar(&globalFormat, "format", "table", "输出格式 (table|json|jsonl)")
	rootCmd.PersistentPreRun = func(cmd *cobra.Command, args []string) {
		output.SetFormat(globalFormat)
	}

	// 初始化日志
	if err := log.Init(utils.ConfigInstance.GetLogPath()); err != nil {
		fmt.Printf("日志初始化失败: %v\n", err)
		os.Exit(1)
	}

	// 注册子命令
	rootCmd.AddCommand(clone.GetCloneCommand())
	rootCmd.AddCommand(backup.GetBackupCommand())
	rootCmd.AddCommand(configcmd.GetConfigCommand())
	rootCmd.AddCommand(trending.GetTrendingCommand())

	log.Info("初始化完成")
}
