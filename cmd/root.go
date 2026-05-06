/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/cicbyte/forks-cli/internal/common"
	"github.com/cicbyte/forks-cli/internal/log"
	"github.com/cicbyte/forks-cli/internal/utils"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
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
	// 初始化日志
	if err := log.Init(utils.ConfigInstance.GetLogPath()); err != nil {
		fmt.Printf("日志初始化失败: %v\n", err)
		os.Exit(1)
	}

	log.Info("初始化完成")
}
