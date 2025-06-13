package common

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/Fromsko/downhub/logs"

	"github.com/gocolly/colly/v2"
)

type (
	DownHub struct {
		*DownType
		DownDir  string
		BaseUrl  string
		ProxyUrl string
		LastTag  string
		RepoName string
		Spider   *colly.Collector
	}
	DownType struct {
		TarGz []string `json:"tar_list,omitempty"`
		Zip   []string `json:"zip_list,omitempty"`
	}
	Option func(*DownHub)
)

var Log = struct {
	Info  func(string, ...interface{})
	Warn  func(string, ...interface{})
	Error func(string, ...interface{})
}{
	Info:  logs.Info,
	Warn:  logs.Warn,
	Error: logs.Error,
}

func (d *DownType) Filter(url string) {
	if strings.HasSuffix(url, ".zip") {
		d.Zip = append(d.Zip, url)
	} else if strings.HasSuffix(url, ".tar.gz") {
		d.TarGz = append(d.TarGz, url)
	}
}

func (hub *DownHub) Link() string {
	baseLink := strings.Split(hub.BaseUrl, "https://github.com")[1]
	repo := strings.Split(hub.BaseUrl, "/")
	hub.RepoName = repo[len(repo)-1]
	return baseLink + "/archive/refs/tags/.*"
}

func WithProxy(proxy string) Option {
	return func(dh *DownHub) {
		dh.ProxyUrl = proxy
	}
}

func WithBaseUrl(url string) Option {
	return func(dh *DownHub) {
		dh.BaseUrl = url
	}
}

func WithDownDir(dir ...string) Option {
	return func(dh *DownHub) {
		var (
			dirPath     string
			basePath, _ = os.Getwd()
		)

		if dirs := len(dir); dirs > 0 {
			dirPath = filepath.Join(
				basePath,
				filepath.Join(dir...),
			)
			if err := os.MkdirAll(dirPath, 0755); err != nil {
				Log.Error("Create folders error!", err)
			}
		} else {
			dirPath = filepath.Join(
				basePath,
				dh.DownDir,
			)
			if err := os.MkdirAll(dirPath, 0755); err != nil {
				Log.Error("Create directory error!", err)
			}
		}
		dh.DownDir = dirPath
	}
}

func WithDefaultSpider() Option {
	return func(dh *DownHub) {
		if dh.Spider == nil {
			dh.Spider = colly.NewCollector()
		}

		dh.Spider.Async = true
		dh.Spider.Limit(&colly.LimitRule{
			DomainGlob:  "*",
			Parallelism: 20,
		})
	}
}

func NewDownHub(opts ...Option) *DownHub {
	dh := &DownHub{
		DownDir:  "source",
		DownType: new(DownType),
	}

	for _, opt := range opts {
		opt(dh)
	}

	// Create spider if not already created
	if dh.Spider == nil {
		dh.Spider = colly.NewCollector()
	}

	// Set proxy if provided
	if dh.ProxyUrl != "" {
		dh.Spider.SetProxy(dh.ProxyUrl)
	}

	return dh
}
