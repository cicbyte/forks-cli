package configcmd

import (
	"fmt"
	"strings"

	"github.com/cicbyte/forks-cli/internal/common"
	"github.com/cicbyte/forks-cli/internal/utils"
	"github.com/spf13/cobra"
)

// GetConfigCommand 返回 config 命令
func GetConfigCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config [子命令]",
		Short: "查看或修改配置",
		Long: `查看或修改 forks-cli 配置

子命令:
  server <url>       设置 Forks 服务端地址
  token <value>      设置 API Token
  backup-dir <path>  设置备份目录
  show               显示当前配置（默认行为）`,
		Args: cobra.MaximumNArgs(3),
		RunE: runConfigShow,
	}

	cmd.AddCommand(getConfigServerCmd())
	cmd.AddCommand(getConfigTokenCmd())
	cmd.AddCommand(getConfigBackupDirCmd())
	cmd.AddCommand(getConfigShowCmd())

	return cmd
}

func runConfigShow(cmd *cobra.Command, args []string) error {
	cfg := common.AppConfigModel
	path := utils.ConfigInstance.GetConfigPath()

	fmt.Printf("配置文件: %s\n", path)
	if cfg.Server != "" {
		fmt.Printf("server:     %s\n", cfg.Server)
	} else {
		fmt.Println("server:     (未设置)")
	}
	if cfg.Token != "" {
		fmt.Printf("token:      %s\n", maskToken(cfg.Token))
	} else {
		fmt.Println("token:      (未设置)")
	}
	if cfg.BackupDir != "" {
		fmt.Printf("backup_dir: %s\n", cfg.BackupDir)
	} else {
		fmt.Println("backup_dir: ./backup (默认)")
	}
	return nil
}

func maskToken(t string) string {
	if len(t) <= 8 {
		return "****"
	}
	return t[:4] + "****" + t[len(t)-4:]
}

func saveConfig() {
	utils.ConfigInstance.SaveConfig(common.AppConfigModel)
}

func getConfigServerCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "server <url>",
		Short: "设置 Forks 服务端地址",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			common.AppConfigModel.Server = strings.TrimSuffix(args[0], "/")
			saveConfig()
			fmt.Printf("已保存 server: %s\n", common.AppConfigModel.Server)
			return nil
		},
	}
}

func getConfigTokenCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "token <value>",
		Short: "设置 API Token",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			common.AppConfigModel.Token = args[0]
			saveConfig()
			fmt.Printf("已保存 token: %s\n", maskToken(common.AppConfigModel.Token))
			return nil
		},
	}
}

func getConfigBackupDirCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "backup-dir <path>",
		Short: "设置本地备份目录",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			common.AppConfigModel.BackupDir = args[0]
			saveConfig()
			fmt.Printf("已保存 backup_dir: %s\n", common.AppConfigModel.BackupDir)
			return nil
		},
	}
}

func getConfigShowCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "show",
		Short: "显示当前配置",
		Args:  cobra.NoArgs,
		RunE:  runConfigShow,
	}
}
