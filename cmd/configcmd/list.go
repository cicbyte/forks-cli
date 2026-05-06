package configcmd

import (
	"fmt"
	"strconv"

	"github.com/cicbyte/forks-cli/internal/common"
	"github.com/cicbyte/forks-cli/internal/output"
	"github.com/cicbyte/forks-cli/internal/utils"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "列出所有配置项",
	Long: `列出所有配置项及当前值。敏感字段（token）会显示脱敏后的值。

示例:
  forks-cli config list
  forks-cli config list --format json`,
	Args: cobra.NoArgs,
	Run:  runList,
}

type configEntry struct {
	key       string
	section   string
	value     string
	sensitive bool
}

func runList(cmd *cobra.Command, args []string) {
	cfg := common.AppConfigModel

	entries := []configEntry{
		{key: "server", section: "General", value: cfg.Server},
		{key: "token", section: "General", value: cfg.Token, sensitive: true},
		{key: "backup_dir", section: "General", value: cfg.BackupDir},
		{key: "log.level", section: "Log", value: cfg.Log.Level},
		{key: "log.max_size", section: "Log", value: strconv.Itoa(cfg.Log.MaxSize)},
		{key: "log.max_backups", section: "Log", value: strconv.Itoa(cfg.Log.MaxBackups)},
		{key: "log.max_age", section: "Log", value: strconv.Itoa(cfg.Log.MaxAge)},
		{key: "log.compress", section: "Log", value: strconv.FormatBool(cfg.Log.Compress)},
	}

	if output.IsJSON() {
		printConfigJSON(entries)
		return
	}

	printConfigTable(entries)
}

func printConfigJSON(entries []configEntry) {
	items := make([]map[string]any, 0, len(entries))
	for _, e := range entries {
		displayVal := e.value
		if e.sensitive {
			displayVal = maskValue(e.value)
		}
		if displayVal == "" {
			displayVal = "(未设置)"
		}
		items = append(items, map[string]any{
			"key":       e.key,
			"value":     displayVal,
			"section":   e.section,
			"sensitive": e.sensitive,
		})
	}

	if output.IsJSONL() {
		output.PrintJSONL(items)
		return
	}

	output.PrintJSON(map[string]any{
		"config_file": utils.ConfigInstance.GetConfigPath(),
		"items":       items,
	})
}

func printConfigTable(entries []configEntry) {
	headers := []string{"KEY", "VALUE"}
	rows := make([][]string, 0, len(entries))
	currentSection := ""

	for _, e := range entries {
		if e.section != currentSection {
			currentSection = e.section
			rows = append(rows, []string{fmt.Sprintf("[%s]", currentSection), ""})
		}

		displayVal := e.value
		if e.sensitive {
			displayVal = maskValue(e.value)
		}
		if displayVal == "" {
			displayVal = "(未设置)"
		}
		rows = append(rows, []string{e.key, displayVal})
	}

	fmt.Printf("配置文件: %s\n\n", utils.ConfigInstance.GetConfigPath())
	output.PrintTable(headers, rows)
}
