package cmd

import (
	"fmt"

	"github.com/Fromsko/downhub/handler"

	"github.com/spf13/cobra"
)

var (
	proxy string
)

func checkAndWarnGithubAccess(proxy string) bool {
	if !handler.CheckGithubAccess(proxy) {
		fmt.Println("[警告] 无法直接访问 github.com，建议通过 --proxy 参数配置命令行代理。")
		return false
	}
	return true
}

var RootCmd = &cobra.Command{
	Use:   "downhub",
	Short: "DownHub is a tool for downloading GitHub repositories",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 1 {
			if !checkAndWarnGithubAccess(proxy) {
				fmt.Println("终止操作，请先通过 --proxy 参数配置好代理后再重试。")
				return
			}
			handler.DownloadRepo(
				args[0],
				proxy,
			)
		} else {
			// show help
			cmd.Help()
		}
	},
}

func init() {
	RootCmd.CompletionOptions.DisableDefaultCmd = true
	RootCmd.SetHelpCommand(&cobra.Command{
		Use:    "no-help",
		Hidden: true,
	})
	RootCmd.Flags().StringVarP(&proxy, "proxy", "p", "", "Proxy URL (如 http://localhost:7890)")
	batchCmd.Flags().StringVarP(&proxy, "proxy", "p", "", "Proxy URL (如 http://localhost:7890)")
	RootCmd.AddCommand(batchCmd)
}

var batchCmd = &cobra.Command{
	Use:   "batch -f [file with repo URLs]",
	Short: "Batch download repositories from a file",
	Run: func(cmd *cobra.Command, args []string) {
		filePath := ""
		if len(args) > 0 {
			filePath = args[0]
		} else {
			filePath, _ = cmd.Flags().GetString("file")
		}
		if filePath == "" {
			fmt.Println("请指定包含仓库地址的文件，如: ./download batch -f repo-list.txt 或 ./download batch repo-list.txt")
			return
		}
		if !checkAndWarnGithubAccess(proxy) {
			fmt.Println("终止操作，请先通过 --proxy 参数配置好代理后再重试。")
			return
		}
		handler.DownloadRepos(
			handler.ReadFromFile(filePath),
			proxy,
		)
	},
}
