package configcmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/cicbyte/forks-cli/internal/common"
	"github.com/spf13/cobra"
)

var getShowFlag bool

var getCmd = &cobra.Command{
	Use:   "get <key>",
	Short: "查看单个配置项的值",
	Long: `查看指定配置项的当前值。敏感字段（token）默认脱敏显示，使用 --show 查看明文。

示例:
  forks-cli config get server
  forks-cli config get token
  forks-cli config get token --show
  forks-cli config get log.level`,
	Args: cobra.ExactArgs(1),
	Run:  runGet,
}

func init() {
	getCmd.Flags().BoolVar(&getShowFlag, "show", false, "显示敏感字段的明文值")
}

func runGet(cmd *cobra.Command, args []string) {
	key := args[0]

	value, ok, sensitive := getConfigValue(key)
	if !ok {
		fmt.Printf("错误: 未知配置项 '%s'\n", key)
		fmt.Println("使用 'forks-cli config list' 查看所有配置项")
		os.Exit(1)
	}

	if value == "" {
		fmt.Printf("%s: (未设置)\n", key)
		return
	}

	if sensitive && !getShowFlag {
		fmt.Printf("%s: %s\n", key, maskValue(value))
		fmt.Println("使用 --show 查看明文")
		return
	}

	fmt.Printf("%s: %s\n", key, value)
}

// getConfigValue 根据键名获取配置值。
// 返回 (值, 是否存在, 是否敏感)。
func getConfigValue(key string) (string, bool, bool) {
	cfg := common.AppConfigModel

	switch key {
	case "server":
		return cfg.Server, true, false
	case "token":
		return cfg.Token, true, true
	case "backup_dir":
		return cfg.BackupDir, true, false
	case "log.level":
		return cfg.Log.Level, true, false
	case "log.max_size":
		return strconv.Itoa(cfg.Log.MaxSize), true, false
	case "log.max_backups":
		return strconv.Itoa(cfg.Log.MaxBackups), true, false
	case "log.max_age":
		return strconv.Itoa(cfg.Log.MaxAge), true, false
	case "log.compress":
		return strconv.FormatBool(cfg.Log.Compress), true, false
	default:
		return "", false, false
	}
}
