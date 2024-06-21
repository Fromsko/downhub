package cmd

import (
	"github.com/Fromsko/downhub/handler"

	"github.com/spf13/cobra"
)

const defaultProxy = "http://localhost:7890"

var RootCmd = &cobra.Command{
	Use:   "downhub",
	Short: "DownHub is a tool for downloading GitHub repositories",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 1 {
			handler.DownloadRepo(
				args[0],
				defaultProxy,
			)
		} else {
			// show help
			cmd.Help()
		}
	},
}

var batchCmd = &cobra.Command{
	Use:   "batch -f [file with repo URLs]",
	Short: "Batch download repositories from a file",
	Run: func(cmd *cobra.Command, args []string) {
		filePath, _ := cmd.Flags().GetString("file")
		proxy, _ := cmd.Flags().GetString("proxy")

		if proxy == "false" {
			proxy = ""
		}

		handler.DownloadRepos(
			handler.ReadFromFile(filePath),
			proxy,
		)
	},
}

func init() {
	RootCmd.CompletionOptions.DisableDefaultCmd = true
	RootCmd.SetHelpCommand(&cobra.Command{
		Use:    "no-help",
		Hidden: true,
	})
	batchCmd.Flags().StringP("file", "f", "", "File containing repository URLs")
	batchCmd.Flags().StringP("proxy", "p", defaultProxy, "Proxy URL (set 'false' to disable)")
	batchCmd.MarkFlagRequired("file")
	RootCmd.AddCommand(batchCmd)
}
