package handler

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/storage/memory"
)

// shouldIncludeFile checks if a file should be included based on the configuration
func shouldIncludeFile(filename string, docsPath string) bool {
	// Use the global cfg variable from handler.go
	// If file filters are configured, use them
	if cfg != nil && len(cfg.FileFilters.Include) > 0 {
		// Check include patterns
		included := false
		for _, pattern := range cfg.FileFilters.Include {
			matched, err := filepath.Match(pattern, filename)
			if err == nil && matched {
				included = true
				break
			}
			// Also check if filename contains the pattern (for simple patterns)
			if strings.Contains(filename, strings.TrimPrefix(pattern, "*")) {
				included = true
				break
			}
		}

		// If not included, skip
		if !included {
			return false
		}

		// Check exclude patterns
		for _, pattern := range cfg.FileFilters.Exclude {
			matched, err := filepath.Match(pattern, filename)
			if err == nil && matched {
				return false
			}
			// Also check if filename contains the pattern (for simple patterns)
			if strings.Contains(filename, strings.TrimPrefix(pattern, "*")) {
				return false
			}
		}

		return true
	}

	// Default behavior: check if file is in docs path or is txt/md file
	return (docsPath != "" && (strings.HasPrefix(filename, docsPath+"/") || filename == docsPath)) ||
		   (strings.HasSuffix(filename, ".txt") || strings.HasSuffix(filename, ".md"))
}

// DownloadDocs downloads txt and md files from a GitHub repository
func DownloadDocs(repoURL, outputDir, docsPath, proxy string) {
	fmt.Printf("Cloning repository: %s\n", repoURL)

	// Extract owner and repo name from URL for directory structure
	var owner, repo string
	if strings.Contains(repoURL, "github.com/") {
		parts := strings.Split(repoURL, "github.com/")
		if len(parts) > 1 {
			repoParts := strings.Split(parts[1], "/")
			if len(repoParts) >= 2 {
				owner = repoParts[0]
				repo = repoParts[1]
			}
		}
	}

	// If we have owner and repo, create subdirectory structure
	if owner != "" && repo != "" {
		outputDir = filepath.Join(outputDir, owner, repo)
	}

	// Clone the repository into memory
	r, err := git.Clone(memory.NewStorage(), nil, &git.CloneOptions{
		URL: repoURL,
	})
	if err != nil {
		fmt.Printf("Error cloning repository: %v\n", err)
		return
	}

	// Get the HEAD reference
	ref, err := r.Head()
	if err != nil {
		fmt.Printf("Error getting HEAD reference: %v\n", err)
		return
	}

	// Get the commit object
	commit, err := r.CommitObject(ref.Hash())
	if err != nil {
		fmt.Printf("Error getting commit object: %v\n", err)
		return
	}

	// Get the tree object
	tree, err := commit.Tree()
	if err != nil {
		fmt.Printf("Error getting tree object: %v\n", err)
		return
	}

	// Walk the tree to find txt and md files
	var filesToDownload []string
	var filePaths []string

	// First, try to find files in the specified docs path
	err = tree.Files().ForEach(func(f *object.File) error {
		// Check if file should be included based on configuration
		if shouldIncludeFile(f.Name, docsPath) {
			filesToDownload = append(filesToDownload, f.Name)
			filePaths = append(filePaths, f.Name)
		}
		return nil
	})

	if err != nil {
		fmt.Printf("Error walking tree: %v\n", err)
		return
	}

	if len(filesToDownload) == 0 {
		fmt.Println("No txt or md files found in repository")
		return
	}

	fmt.Printf("Found %d files to download\n", len(filesToDownload))

	// Download each file
	for i, filePath := range filePaths {
		err := downloadFileFromRepo(tree, filePath, outputDir, filesToDownload[i])
		if err != nil {
			fmt.Printf("Error downloading %s: %v\n", filePath, err)
		} else {
			fmt.Printf("Downloaded: %s\n", filePath)
		}
	}

	fmt.Printf("Download completed. Files saved to: %s\n", outputDir)
}

// DownloadDocsToDataDir downloads txt and md files from a GitHub repository to the data directory structure
func DownloadDocsToDataDir(repoURL, docsPath, proxy, docsBaseDir string) {
	// Extract owner and repo name from URL
	// e.g., https://github.com/gin-gonic/gin -> gin-gonic/gin
	var owner, repo string
	if strings.Contains(repoURL, "github.com/") {
		parts := strings.Split(repoURL, "github.com/")
		if len(parts) > 1 {
			repoParts := strings.Split(parts[1], "/")
			if len(repoParts) >= 2 {
				owner = repoParts[0]
				repo = repoParts[1]
			}
		}
	}

	// Create output directory with owner/repo structure
	outputDir := filepath.Join(docsBaseDir, owner, repo)
	fmt.Printf("Cloning repository: %s\n", repoURL)

	// Clone the repository into memory
	r, err := git.Clone(memory.NewStorage(), nil, &git.CloneOptions{
		URL: repoURL,
	})
	if err != nil {
		fmt.Printf("Error cloning repository: %v\n", err)
		return
	}

	// Get the HEAD reference
	ref, err := r.Head()
	if err != nil {
		fmt.Printf("Error getting HEAD reference: %v\n", err)
		return
	}

	// Get the commit object
	commit, err := r.CommitObject(ref.Hash())
	if err != nil {
		fmt.Printf("Error getting commit object: %v\n", err)
		return
	}

	// Get the tree object
	tree, err := commit.Tree()
	if err != nil {
		fmt.Printf("Error getting tree object: %v\n", err)
		return
	}

	// Walk the tree to find txt and md files
	var filesToDownload []string
	var filePaths []string

	// First, try to find files in the specified docs path
	err = tree.Files().ForEach(func(f *object.File) error {
		// Check if file should be included based on configuration
		if shouldIncludeFile(f.Name, docsPath) {
			filesToDownload = append(filesToDownload, f.Name)
			filePaths = append(filePaths, f.Name)
		}
		return nil
	})

	if err != nil {
		fmt.Printf("Error walking tree: %v\n", err)
		return
	}

	if len(filesToDownload) == 0 {
		fmt.Println("No txt or md files found in repository")
		return
	}

	fmt.Printf("Found %d files to download\n", len(filesToDownload))

	// Download each file
	for i, filePath := range filePaths {
		err := downloadFileFromRepo(tree, filePath, outputDir, filesToDownload[i])
		if err != nil {
			fmt.Printf("Error downloading %s: %v\n", filePath, err)
		} else {
			fmt.Printf("Downloaded: %s\n", filePath)
		}
	}

	fmt.Printf("Download completed. Files saved to: %s\n", outputDir)
}

// downloadFileFromRepo downloads a file from the repository tree
func downloadFileFromRepo(tree *object.Tree, filePath, outputDir, savePath string) error {
	// Get the file from the tree
	file, err := tree.File(filePath)
	if err != nil {
		return fmt.Errorf("error getting file %s: %v", filePath, err)
	}

	// Create the output directory structure
	fullOutputPath := filepath.Join(outputDir, savePath)
	dir := filepath.Dir(fullOutputPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("error creating directory %s: %v", dir, err)
	}

	// Create the output file
	outFile, err := os.Create(fullOutputPath)
	if err != nil {
		return fmt.Errorf("error creating file %s: %v", fullOutputPath, err)
	}
	defer outFile.Close()

	// Get the file contents
	reader, err := file.Reader()
	if err != nil {
		return fmt.Errorf("error reading file %s: %v", filePath, err)
	}
	defer reader.Close()

	// Copy the contents to the output file
	_, err = io.Copy(outFile, reader)
	if err != nil {
		return fmt.Errorf("error writing file %s: %v", fullOutputPath, err)
	}

	return nil
}

// downloadFileOverHTTP downloads a file over HTTP with optional proxy
func downloadFileOverHTTP(fileURL, outputPath, proxy string) error {
	var client *http.Client

	if proxy != "" {
		proxyURL, err := url.Parse(proxy)
		if err != nil {
			return fmt.Errorf("invalid proxy URL: %v", err)
		}
		client = &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyURL(proxyURL),
			},
		}
	} else {
		client = &http.Client{}
	}

	resp, err := client.Get(fileURL)
	if err != nil {
		return fmt.Errorf("error downloading file: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP error %d: %s", resp.StatusCode, resp.Status)
	}

	outFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("error creating file: %v", err)
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, resp.Body)
	if err != nil {
		return fmt.Errorf("error writing file: %v", err)
	}

	return nil
}
