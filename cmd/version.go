package cmd

import (
	"fmt"

	"github.com/cicbyte/forks-cli/cmd/version"
	"github.com/spf13/cobra"
)

func getVersionCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "显示版本信息",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("forks-cli %s\n", version.Version)
			fmt.Printf("  commit:  %s\n", version.GitCommit)
			fmt.Printf("  built:   %s\n", version.BuildTime)
		},
	}
}

func init() {
	rootCmd.AddCommand(getVersionCommand())
}
