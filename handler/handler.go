package handler

import (
	"github.com/Fromsko/downhub/common"
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

	"github.com/gocolly/colly/v2"
	"github.com/panjf2000/ants/v2"
	"github.com/schollz/progressbar/v3"
)

var wg sync.WaitGroup

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
		common.Log.Errorf(
			"%s not found!", filepath.Base(basePath),
		)
	} else {
		f, err := os.ReadFile(basePath)
		if err != nil {
			common.Log.Error("ReadFile Error: ", err)
		}
		urlList = strings.Split(strings.TrimSpace(string(f)), "\n")
	}
	return
}

func Repo(hub *common.downhub) {

	matchRegex := regexp.MustCompile(hub.Link())
	hub.DownDir = filepath.Join(hub.DownDir, hub.RepoName)
	_ = os.Mkdir(hub.DownDir, 0644)

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
			common.Log.Info("Next repo :> ", nxp)
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

func downFile(fetchUrl, dir string, bar *progressbar.ProgressBar) error {
	tokens := strings.Split(fetchUrl, "/")
	fileName := tokens[len(tokens)-1]

	out, err := os.Create(filepath.Join(dir, fileName))
	if err != nil {
		return err
	}

	// TODO: It should be dynamic, not fixed.
	proxyURL, err := url.Parse("http://localhost:7890")
	if err != nil {
		return fmt.Errorf("invalid proxy URL: %v", err)
	}

	httpClient := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxyURL),
		},
	}

	resp, err := httpClient.Get(fetchUrl)
	if err != nil {
		return fmt.Errorf("request failed: %v", err)
	}

	writer := io.MultiWriter(out, bar)

	_, err = io.Copy(writer, resp.Body)
	if err != nil {
		return err
	}

	defer out.Close()
	defer resp.Body.Close()

	return nil
}

func saveResult(hub *common.downhub) {
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

func taskBar(url, dir string) {
	bar := progressbar.DefaultBytes(-1, "Download: "+url)

	if err := downFile(url, dir, bar); err != nil {
		common.Log.Error("The network caused the request to fail.")
		return
	}
}

func DownloadRepo(url, proxy string, opts ...common.Option) {
	hub := common.Newdownhub(
		common.WithBaseUrl(url),
		common.WithProxy(proxy),
		common.WithDownDir("source"),
		common.WithDefaultSpider(),
	)

	Repo(hub)

	poolSize := len(hub.Zip)
	pool, _ := ants.NewPool(poolSize)

	for _, zipUrl := range hub.Zip {
		wg.Add(1)
		// Prevent resource competition
		url := zipUrl
		task := func() {
			defer wg.Done()
			taskBar(url, hub.DownDir)
			saveResult(hub)
			if poolSize != 1 {
				fmt.Println()
			}
		}

		_ = pool.Submit(task)
	}

	wg.Wait()
	defer pool.Release()
}

func DownloadRepos(repos []string, proxy string) {
	for _, repoUrl := range repos {
		if repoUrl != "" {
			DownloadRepo(repoUrl, proxy)
		}
	}
}
