package configcmd

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/cicbyte/forks-cli/internal/common"
	"github.com/cicbyte/forks-cli/internal/models"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var setCmd = &cobra.Command{
	Use:   "set <key> [value]",
	Short: "设置配置项的值",
	Long: `设置指定配置项的值。

敏感字段（token）如果不提供 value 参数，会以不回显方式交互式输入。

示例:
  forks-cli config set server http://192.168.1.100:8080
  forks-cli config set token sk-xxx
  forks-cli config set token              # 交互式输入（不回显）
  forks-cli config set backup_dir /data/backup
  forks-cli config set log.level debug
  forks-cli config set log.compress false`,
	Args: cobra.RangeArgs(1, 2),
	Run:  runSet,
}

// sensitiveKeys 敏感配置项集合
var sensitiveKeys = map[string]bool{
	"token": true,
}

func runSet(cmd *cobra.Command, args []string) {
	key := args[0]

	// 检查 key 是否有效
	_, ok, _ := getConfigValue(key)
	if !ok {
		fmt.Printf("错误: 未知配置项 '%s'\n", key)
		fmt.Println("使用 'forks-cli config list' 查看所有配置项")
		os.Exit(1)
	}

	var value string

	if len(args) >= 2 {
		value = args[1]
	} else if sensitiveKeys[key] {
		// 敏感字段交互式输入（不回显）
		fmt.Printf("请输入 %s: ", key)
		raw, err := term.ReadPassword(int(os.Stdin.Fd()))
		fmt.Println()
		if err != nil {
			fmt.Println("错误: 读取输入失败")
			os.Exit(1)
		}
		value = string(raw)
	} else {
		// 普通字段交互式输入
		fmt.Printf("请输入 %s: ", key)
		reader := bufio.NewReader(os.Stdin)
		line, _ := reader.ReadString('\n')
		value = strings.TrimSpace(line)
	}

	if value == "" {
		fmt.Println("错误: 值不能为空")
		os.Exit(1)
	}

	// server 字段自动去除末尾斜杠
	if key == "server" {
		value = strings.TrimSuffix(value, "/")
	}

	// 类型校验并设置值
	if err := setConfigValue(common.AppConfigModel, key, value); err != nil {
		fmt.Printf("错误: %v\n", err)
		os.Exit(1)
	}

	saveConfig()
	fmt.Printf("%s 已更新\n", key)
}

// setConfigValue 设置配置值，包含类型校验
func setConfigValue(c *models.AppConfig, key, value string) error {
	switch key {
	case "server":
		c.Server = value
	case "token":
		c.Token = value
	case "backup_dir":
		c.BackupDir = value
	case "log.level":
		c.Log.Level = value
	case "log.max_size":
		v, err := strconv.Atoi(value)
		if err != nil {
			return fmt.Errorf("无效的整数值: %s", value)
		}
		c.Log.MaxSize = v
	case "log.max_backups":
		v, err := strconv.Atoi(value)
		if err != nil {
			return fmt.Errorf("无效的整数值: %s", value)
		}
		c.Log.MaxBackups = v
	case "log.max_age":
		v, err := strconv.Atoi(value)
		if err != nil {
			return fmt.Errorf("无效的整数值: %s", value)
		}
		c.Log.MaxAge = v
	case "log.compress":
		v, err := strconv.ParseBool(value)
		if err != nil {
			return fmt.Errorf("无效的布尔值: %s (true/false)", value)
		}
		c.Log.Compress = v
	}
	return nil
}
