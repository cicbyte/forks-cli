package configcmd

import (
	"fmt"
	"strconv"

	"github.com/cicbyte/forks-cli/internal/common"
	"github.com/cicbyte/forks-cli/internal/utils"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "列出所有配置项",
	Long: `列出所有配置项及当前值。敏感字段（token）会显示脱敏后的值。

示例:
  forks-cli config list`,
	Args: cobra.NoArgs,
	Run:  runList,
}

type configEntry struct {
	key       string
	value     string
	sensitive bool
}

func runList(cmd *cobra.Command, args []string) {
	cfg := common.AppConfigModel

	entries := []configEntry{
		// General
		{key: "server", value: cfg.Server},
		{key: "token", value: cfg.Token, sensitive: true},
		{key: "backup_dir", value: cfg.BackupDir},
		// Log
		{key: "log.level", value: cfg.Log.Level},
		{key: "log.max_size", value: strconv.Itoa(cfg.Log.MaxSize)},
		{key: "log.max_backups", value: strconv.Itoa(cfg.Log.MaxBackups)},
		{key: "log.max_age", value: strconv.Itoa(cfg.Log.MaxAge)},
		{key: "log.compress", value: strconv.FormatBool(cfg.Log.Compress)},
	}

	fmt.Printf("配置文件: %s\n\n", utils.ConfigInstance.GetConfigPath())

	fmt.Println("[General]")
	for i := 0; i < 3; i++ {
		printEntry(entries[i])
	}
	fmt.Println()

	fmt.Println("[Log]")
	for i := 3; i < len(entries); i++ {
		printEntry(entries[i])
	}
}

func printEntry(e configEntry) {
	displayVal := e.value
	if e.sensitive {
		displayVal = maskValue(e.value)
	}
	if displayVal == "" {
		displayVal = "(未设置)"
	}
	fmt.Printf("  %-16s %s\n", e.key+":", displayVal)
}
