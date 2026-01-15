package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Fromsko/downhub/config"
	"github.com/Fromsko/downhub/handler"

	"github.com/spf13/cobra"
)

var (
	proxy string
	cfg   *config.Config
)

// SetConfig sets the configuration for the command
func SetConfig(c *config.Config) {
	cfg = c
}

func checkAndWarnGithubAccess(proxy string) bool {
	// First try direct access
	if handler.CheckGithubAccess("") {
		return true
	}

	// If direct access fails and proxy is provided, try with proxy
	if proxy != "" {
		if handler.CheckGithubAccess(proxy) {
			return true
		}
		fmt.Println("[警告] 无法直接访问 github.com，且代理配置无效。")
	} else {
		fmt.Println("[警告] 无法直接访问 github.com，建议通过 --proxy 参数配置命令行代理。")
	}
	return false
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
			handler.DownloadRepoToDataDir(
				args[0],
				proxy,
			)
		} else {
			// show help
			if err := cmd.Help(); err != nil {
				fmt.Fprintln(os.Stderr, err)
			}
		}
	},
}

func init() {
	RootCmd.CompletionOptions.DisableDefaultCmd = true
	RootCmd.SetHelpCommand(&cobra.Command{
		Use:    "no-help",
		Hidden: true,
	})
	// Set default proxy from config if available
	if cfg != nil && cfg.Defaults.Proxy != "" {
		proxy = cfg.Defaults.Proxy
	}
	RootCmd.Flags().StringVarP(&proxy, "proxy", "p", proxy, "Proxy URL (如 http://localhost:7890)")
	batchCmd.Flags().StringVarP(&proxy, "proxy", "p", proxy, "Proxy URL (如 http://localhost:7890)")
	RootCmd.AddCommand(batchCmd)
	RootCmd.AddCommand(docsCmd)
	RootCmd.AddCommand(commonCmd)
}

var commonCmd = &cobra.Command{
	Use:   "common",
	Short: "Download repositories configured in YAML file",
	Run: func(cmd *cobra.Command, args []string) {
		if cfg == nil {
			fmt.Println("Configuration not loaded")
			return
		}

		if !checkAndWarnGithubAccess(proxy) {
			fmt.Println("终止操作，请先通过 --proxy 参数配置好代理后再重试。")
			return
		}

		// Download repositories configured in YAML
		for _, repo := range cfg.Repositories {
			fmt.Printf("Downloading repository: %s\n", repo.Name)
			if repo.DownloadDocs {
				// Download docs to base_data_dir/docs_dir/owner/repo structure
				baseDataDir := cfg.Defaults.BaseDataDir
				docsDir := cfg.Defaults.DocsDir
				if baseDataDir == "" {
					baseDataDir = "data"
				}
				if docsDir == "" {
					docsDir = "docs"
				}
				handler.DownloadDocsToDataDir(repo.URL, repo.DocsPath, proxy, filepath.Join(baseDataDir, docsDir))
			}
			if repo.DownloadSource {
				// Download source to data/source/owner/repo structure
				handler.DownloadRepoToDataDir(repo.URL, proxy)
			}
		}
	},
}

func init() {
	commonCmd.Flags().StringVarP(&proxy, "proxy", "p", proxy, "Proxy URL (如 http://localhost:7890)")
}

var docsCmd = &cobra.Command{
	Use:   "docs [repo-url]",
	Short: "Download txt and md files from a GitHub repository",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		repoURL := args[0]
		outputDir, _ := cmd.Flags().GetString("output")
		docsPath, _ := cmd.Flags().GetString("docs-path")

		if outputDir == "" {
			// Use default docs directory: base_data_dir/docs_dir
			if cfg != nil && cfg.Defaults.BaseDataDir != "" && cfg.Defaults.DocsDir != "" {
				outputDir = filepath.Join(cfg.Defaults.BaseDataDir, cfg.Defaults.DocsDir)
			} else {
				outputDir = "data/docs"
			}
		}

		if docsPath == "" {
			// Use default docs path from config if available
			if cfg != nil && cfg.Defaults.DocsPath != "" {
				docsPath = cfg.Defaults.DocsPath
			} else {
				docsPath = "docs"
			}
		}

		// Create output directory if it doesn't exist
		if err := os.MkdirAll(outputDir, 0755); err != nil {
			fmt.Printf("Error creating output directory: %v\n", err)
			return
		}

		// Download docs using the handler
		handler.DownloadDocs(repoURL, outputDir, docsPath, proxy)
	},
}

func init() {
	docsCmd.Flags().StringP("output", "o", "", "Output directory for downloaded files")
	docsCmd.Flags().StringP("docs-path", "d", "docs", "Path to docs directory in repository")
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
