package website

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"

	"github.com/cicbyte/forks-cli/internal/common"
	"github.com/spf13/cobra"
)

var websiteFlagServer string

// GetWebsiteCommand 返回 website 命令
func GetWebsiteCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "website",
		Short: "在浏览器中打开 Forks Web UI",
		Long: `在默认浏览器中打开 Forks 服务端的 Web 管理界面。

示例:
  forks-cli website
  forks-cli website --server http://192.168.1.100:8080`,
		Args: cobra.NoArgs,
		RunE: runWebsite,
	}

	cmd.Flags().StringVar(&websiteFlagServer, "server", "", "本次使用的服务端地址（不保存）")

	return cmd
}

func runWebsite(cmd *cobra.Command, args []string) error {
	cfg := common.AppConfigModel

	server := cfg.Server
	if websiteFlagServer != "" {
		server = strings.TrimSuffix(websiteFlagServer, "/")
	}
	if server == "" {
		return fmt.Errorf("请先配置服务端地址: forks-cli config set server <url>")
	}

	url := server
	fmt.Printf("正在打开 %s ...\n", url)

	var err error
	switch runtime.GOOS {
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = exec.Command("xdg-open", url).Start()
	}

	if err != nil {
		return fmt.Errorf("打开浏览器失败: %w", err)
	}
	return nil
}
