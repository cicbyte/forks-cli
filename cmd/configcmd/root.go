package configcmd

import (
	"github.com/cicbyte/forks-cli/internal/common"
	"github.com/cicbyte/forks-cli/internal/utils"
	"github.com/spf13/cobra"
)

// GetConfigCommand 返回 config 命令
func GetConfigCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "管理应用配置",
		Long: `管理 forks-cli 应用配置。

示例:
  forks-cli config list
  forks-cli config get server
  forks-cli config set server http://192.168.1.100:8080
  forks-cli config set token              # 交互式输入（不回显）
  forks-cli config set token sk-xxx       # 直接设置`,
	}
	cmd.AddCommand(listCmd, getCmd, setCmd)
	return cmd
}

// maskValue 对敏感值进行脱敏
func maskValue(value string) string {
	if value == "" {
		return ""
	}
	runes := []rune(value)
	if len(runes) <= 8 {
		return "******"
	}
	return string(runes[:4]) + "..." + string(runes[len(runes)-4:])
}

// saveConfig 保存当前配置到文件
func saveConfig() {
	utils.ConfigInstance.SaveConfig(common.AppConfigModel)
}
