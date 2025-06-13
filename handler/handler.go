package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/Fromsko/downhub/common"

	"github.com/gocolly/colly/v2"
	"github.com/vbauerster/mpb/v8"
	"github.com/vbauerster/mpb/v8/decor"
)

const DefaultProxy = "http://localhost:7890"

func ReadFromFile(file string) (urlList []string) {
	basePath, err := filepath.Abs(file)
	if err != nil {
		common.Log.Error("Convert path: ", err)
	}

	exists := func(filename string) bool {
		info, err := os.Stat(filename)
		if os.IsNotExist(err) {
			return false
		}
		return !info.IsDir()
	}

	if !exists(basePath) {
		common.Log.Error("%s not found!", filepath.Base(basePath))
	} else {
		f, err := os.ReadFile(basePath)
		if err != nil {
			common.Log.Error("ReadFile Error: ", err)
		}
		urlList = strings.Split(strings.TrimSpace(string(f)), "\n")
	}
	return
}

func Repo(hub *common.DownHub) {

	matchRegex := regexp.MustCompile(hub.Link())
	hub.DownDir = filepath.Join(hub.DownDir, hub.RepoName)
	if err := os.MkdirAll(hub.DownDir, 0755); err != nil {
		common.Log.Error("Create directory error!", err)
	}

	hub.Spider.OnHTML("a.Link--muted[href]", func(e *colly.HTMLElement) {
		if matchRegex.MatchString(e.Attr("href")) {
			link := e.Request.AbsoluteURL(e.Attr("href"))
			hub.Filter(link)
		}
	})

	hub.Spider.OnHTML("h2.f4.d-inline a.Link--primary", func(e *colly.HTMLElement) {
		version := e.Text
		hub.LastTag = version
	})

	hub.Spider.OnScraped(func(r *colly.Response) {
		nxp := hub.BaseUrl + "/tags?after=" + hub.LastTag
		if r.Request.URL.String() != nxp {
			common.Log.Info("Next repo :> %s", nxp)
			hub.Spider.Visit(nxp)
		}
	})

	err := hub.Spider.Visit(hub.BaseUrl + "/tags")
	if err != nil {
		common.Log.Error("Visiting URL:", hub.BaseUrl, "-", err)
	} else {
		common.Log.Info("Visiting  :> ", hub.BaseUrl)
	}

	hub.Spider.Wait()

	defer func() {
		common.Log.Info("Finished tasks!")
	}()
}

func downFile(fetchUrl, dir string, bar *mpb.Bar, proxy string) error {
	tokens := strings.Split(fetchUrl, "/")
	fileName := tokens[len(tokens)-1]

	out, err := os.Create(filepath.Join(dir, fileName))
	if err != nil {
		bar.Abort(false)
		return err
	}

	var httpClient *http.Client
	if proxy != "" {
		proxyURL, err := url.Parse(proxy)
		if err != nil {
			bar.Abort(false)
			return fmt.Errorf("invalid proxy URL: %v", err)
		}
		httpClient = &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyURL(proxyURL),
			},
		}
	} else {
		httpClient = &http.Client{}
	}

	resp, err := httpClient.Get(fetchUrl)
	if err != nil {
		bar.Abort(false)
		return fmt.Errorf("request failed: %v", err)
	}
	defer resp.Body.Close()
	defer out.Close()

	total := resp.ContentLength
	if total > 0 {
		bar.SetTotal(total, false)
	}
	reader := bar.ProxyReader(resp.Body)
	_, err = io.Copy(out, reader)
	if err != nil {
		bar.Abort(false)
		return err
	}
	bar.SetTotal(bar.Current(), true)
	return nil
}

func saveResult(hub *common.DownHub) {
	type save struct {
		Name string           `json:"repo_name"`
		Url  string           `json:"repo_url"`
		Tags *common.DownType `json:"download_sources"`
	}

	f, _ := os.Create("output-" + hub.RepoName + ".json")

	json.NewEncoder(f).Encode(save{
		Name: hub.RepoName,
		Url:  hub.BaseUrl,
		Tags: hub.DownType,
	})

	defer f.Close()
}

func DownloadRepo(url string, proxy string, opts ...common.Option) {
	hub := common.NewDownHub(common.WithBaseUrl(url), common.WithProxy(proxy), common.WithDefaultSpider())
	for _, opt := range opts {
		opt(hub)
	}

	Repo(hub)
	total := len(hub.Zip) + len(hub.TarGz)
	if total == 0 {
		common.Log.Info("No files to download")
		return
	}

	var success, failed int
	var mu sync.Mutex
	p := mpb.New(mpb.WithWidth(60))
	wg := sync.WaitGroup{}

	downloadFiles := func(fileList []string) {
		for _, fileURL := range fileList {
			wg.Add(1)
			fileURL := fileURL
			fileName := filepath.Base(fileURL)
			bar := p.New(0,
				mpb.BarStyle().Rbound("⠿").Filler("⠶").Tip("⠿").Padding(" "),
				mpb.PrependDecorators(
					decor.Name(fileName+" ", decor.WC{W: 30, C: decor.DSyncWidth}),
				),
				mpb.AppendDecorators(
					decor.Percentage(decor.WC{W: 5}),
					decor.CountersKibiByte("% .1f / % .1f"),
				),
			)
			go func(url, dir, proxy string, bar *mpb.Bar) {
				defer wg.Done()
				err := downFile(url, dir, bar, proxy)
				mu.Lock()
				if err == nil {
					success++
				} else {
					failed++
					common.Log.Error("下载失败: %s, %v", url, err)
				}
				mu.Unlock()
				saveResult(hub)
			}(fileURL, hub.DownDir, hub.ProxyUrl, bar)
		}
	}
	downloadFiles(hub.Zip)
	downloadFiles(hub.TarGz)
	wg.Wait()
	p.Wait()
	common.Log.Info("下载完成，总数: %d，成功: %d，失败: %d，存放目录: %s", total, success, failed, hub.DownDir)
}

func DownloadRepos(repos []string, proxy string) {
	for _, repoUrl := range repos {
		if repoUrl != "" {
			DownloadRepo(repoUrl, proxy)
		}
	}
}

// 检查能否访问 github.com
func CheckGithubAccess(proxy string) bool {
	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	if proxy != "" {
		proxyURL, err := url.Parse(proxy)
		if err == nil {
			client.Transport = &http.Transport{Proxy: http.ProxyURL(proxyURL)}
		}
	}
	resp, err := client.Get("https://github.com/")
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode == 200
}
