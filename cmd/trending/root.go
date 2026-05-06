package trending

import (
	"fmt"
	"os"
	"strings"

	"github.com/cicbyte/forks-cli/internal/common"
	logictrending "github.com/cicbyte/forks-cli/internal/logic/trending"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
)

var (
	trendingFlagToken          string
	trendingFlagServer         string
	trendingFlagLanguage       string
	trendingFlagSince          string
	trendingFlagSpokenLanguage string
	trendingFlagDate           string
	trendingFlagRefresh        bool
)

// GetTrendingCommand 返回 trending 命令
func GetTrendingCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "trending [flags]",
		Short: "浏览 GitHub Trending 趋势仓库",
		Long: `从 Forks 服务端获取 GitHub Trending 趋势仓库列表。

支持按编程语言、时间范围、自然语言筛选。

示例:
  forks-cli trending                          # 今日全部语言趋势
  forks-cli trending -l go                    # Go 语言趋势
  forks-cli trending -l python -s weekly      # Python 周趋势
  forks-cli trending -S zh                    # 中文趋势
  forks-cli trending -d 2026-05-04            # 查看历史趋势
  forks-cli trending --refresh                # 刷新缓存`,
		Args: cobra.NoArgs,
		RunE: runTrending,
	}

	cmd.Flags().StringVarP(&trendingFlagLanguage, "language", "l", "", "编程语言筛选 (如 go, python, rust)")
	cmd.Flags().StringVarP(&trendingFlagSince, "since", "s", "daily", "时间范围: daily/weekly/monthly")
	cmd.Flags().StringVarP(&trendingFlagSpokenLanguage, "spoken", "S", "", "自然语言筛选 (如 zh, en)")
	cmd.Flags().StringVarP(&trendingFlagDate, "date", "d", "", "指定日期 (格式 2026-05-06)")
	cmd.Flags().BoolVar(&trendingFlagRefresh, "refresh", false, "跳过缓存重新获取")
	cmd.Flags().StringVarP(&trendingFlagToken, "token", "t", "", "本次使用的 API Token（不保存）")
	cmd.Flags().StringVarP(&trendingFlagServer, "server", "", "", "本次使用的服务端地址（不保存）")

	return cmd
}

func runTrending(cmd *cobra.Command, args []string) error {
	cfg := common.AppConfigModel

	// server 优先级: 命令行 > 配置文件
	server := cfg.Server
	if trendingFlagServer != "" {
		server = strings.TrimSuffix(trendingFlagServer, "/")
	}
	if server == "" {
		return fmt.Errorf("请先配置服务端地址: forks-cli config set server <url>")
	}

	// token 优先级: 命令行 > 环境变量 > 配置文件
	token := trendingFlagToken
	if token == "" {
		token = os.Getenv("FORKS_TOKEN")
	}
	if token == "" {
		token = cfg.Token
	}

	config := &logictrending.TrendingConfig{
		Server:         server,
		Token:          token,
		Language:       trendingFlagLanguage,
		Since:          trendingFlagSince,
		SpokenLanguage: trendingFlagSpokenLanguage,
		Date:           trendingFlagDate,
		Refresh:        trendingFlagRefresh,
	}

	repos, date, err := logictrending.FetchTrending(config)
	if err != nil {
		return err
	}

	if len(repos) == 0 {
		fmt.Println("没有找到趋势仓库")
		return nil
	}

	// 输出标题
	sinceLabel := map[string]string{"daily": "今日", "weekly": "本周", "monthly": "本月"}
	label := sinceLabel[config.Since]
	if label == "" {
		label = config.Since
	}
	fmt.Printf("\n%s GitHub Trending", label)
	if config.Language != "" {
		fmt.Printf(" (%s)", config.Language)
	}
	fmt.Printf(" — %s\n", date)

	// 构建表格
	t := table.NewWriter()
	t.SetStyle(table.StyleDefault)
	t.Style().Options.DrawBorder = false
	t.Style().Options.SeparateColumns = false

	t.AppendHeader(table.Row{"#", "REPO", "DESCRIPTION", "LANG", "STARS", "TODAY"})
	for i, repo := range repos {
		desc := repo.Description
		if len([]rune(desc)) > 50 {
			desc = string([]rune(desc)[:47]) + "..."
		}

		stars := formatNum(repo.Stars)
		today := ""
		if repo.CurrentPeriodStars > 0 {
			today = fmt.Sprintf("+%s", formatNum(repo.CurrentPeriodStars))
		}

		t.AppendRow(table.Row{
			i + 1,
			fmt.Sprintf("%s/%s", repo.Author, repo.Repo),
			desc,
			repo.Language,
			stars,
			today,
		})
	}

	fmt.Println(t.Render())
	fmt.Printf("\n共 %d 个仓库\n", len(repos))
	return nil
}

func formatNum(n int) string {
	if n >= 1000 {
		return fmt.Sprintf("%.1fk", float64(n)/1000)
	}
	return fmt.Sprintf("%d", n)
}
