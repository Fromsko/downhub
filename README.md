# DownHub - Github Release 快捷下载工具

DownHub 是一个用于从 GitHub 快速下载发布版本（Release）的命令行工具，支持单仓库和批量下载，支持代理，支持多文件并发下载与美观的多进度条展示。

---

## ✨ 功能特性

- 支持下载指定仓库所有 Release 的 zip/tar.gz 包
- 支持批量下载（通过文件列表）
- 支持 HTTP/HTTPS 代理
- 多文件并发下载，进度条美观直观
- 下载完成后统计成功/失败数与存放目录
- 彩色日志输出，时间+级别清晰
- 支持YAML配置文件，可自定义下载行为和默认设置
- 支持文档文件下载（.md, .txt等）
- 支持通过YAML配置文件管理多个仓库的下载任务
- 支持文件过滤器，可自定义包含或排除特定文件
- 支持断点续传和下载重试机制

---

## 📦 安装

```sh
git clone https://github.com/Fromsko/downhub.git
cd downhub
go build -o download main.go
```

---

## 🚀 快速开始

### 单仓库下载

```sh
./downhub https://github.com/gin-gonic/gin
```

### 使用代理

```sh
./downhub --proxy http://localhost:7890 https://github.com/gin-gonic/gin
```

### 批量下载

准备一个包含多个仓库地址的文本文件（每行一个）：

```txt
https://github.com/gin-gonic/gin
https://github.com/labstack/echo
...
```

执行批量下载：

```sh
./downhub batch -f repo-list.txt
```

### 文档下载

下载指定仓库的文档文件：

```sh
./downhub docs https://github.com/gin-gonic/gin
```

指定输出目录：

```sh
./downhub docs https://github.com/gin-gonic/gin -o ./my-docs
```

指定文档路径：

```sh
./downhub docs https://github.com/gin-gonic/gin -d documentation
```

---

## ⚙️ 命令行参数

- `-p, --proxy` 指定代理地址（如 http://localhost:7890）
- `batch -f` 批量下载，指定包含仓库地址的文件
- `-h, --help` 查看帮助

![command](res/command.png)

---

## 🖥️ 进度与日志

- 每个文件下载均有独立进度条，支持多文件并发美观展示
- 日志输出带时间戳，级别彩色区分，便于排查问题
- 下载结束后自动统计总数、成功、失败、存放目录

![show](res/show.gif)

---

## 🛠️ 开发&贡献

欢迎提交 PR 或 Issue！

1. 克隆仓库
   ```sh
   git clone https://github.com/Fromsko/downhub.git
   cd downhub
   ```
2. 构建

   ```sh
   go build -o download main.go
   ```

3. 运行

   ```sh
   ./download --help
   ```

---

## 🙏 鸣谢

- [Colly](https://github.com/gocolly/colly) 网页爬取
- [Cobra](https://github.com/spf13/cobra) 命令行解析
- [mpb](https://github.com/vbauerster/mpb) 多进度条
- [fatih/color](https://github.com/fatih/color) 彩色日志

---

## 🛠️ 配置文件

DownHub 支持通过 YAML 配置文件进行自定义设置。首次运行时，如果不存在配置文件，程序会自动创建一个默认配置文件 `downhub.yaml`。

### 配置文件示例

```yaml
# Downhub Configuration File

defaults:
  output_dir: "./downloads"
  docs_path: "docs"
  max_concurrent_downloads: 5
  proxy: "http://localhost:7890"

repositories:
  - name: "go-git"
    url: "https://github.com/go-git/go-git"
    download_docs: true
    download_source: false
    output_dir: "./downloads/go-git"
    docs_path: "docs"

  - name: "trpc-agent-go"
    url: "https://github.com/trpc-group/trpc-agent-go"
    download_docs: true
    download_source: false
    output_dir: "./downloads/trpc-agent-go"
    docs_path: "docs"

file_filters:
  include:
    - "*.md"
    - "*.txt"
    - "*.yaml"
    - "*.yml"
  exclude:
    - "node_modules/*"
    - ".git/*"
    - "vendor/*"

download:
  timeout: 300
  retries: 3
  retry_delay: 5
  user_agent: "Downhub/1.0"

logging:
  level: "info"
  format: "text"
  output: "stdout"

advanced:
  preserve_structure: true
  create_readme: true
  validate_checksums: false
```

### 配置选项详解

- `defaults`: 默认设置
  - `output_dir`: 默认输出目录，所有下载文件将保存到此目录
  - `docs_path`: 默认文档路径，在仓库中查找文档的默认路径
  - `max_concurrent_downloads`: 最大并发下载数，控制同时下载的文件数量
  - `proxy`: 默认代理地址，如果未通过命令行指定代理，则使用此设置

- `repositories`: 仓库配置列表
  - `name`: 仓库名称
  - `url`: 仓库URL
  - `download_docs`: 是否下载文档文件
  - `download_source`: 是否下载源代码包
  - `output_dir`: 该仓库的输出目录
  - `docs_path`: 该仓库的文档路径

- `file_filters`: 文件过滤器
  - `include`: 包含的文件模式，只有匹配这些模式的文件才会被下载
  - `exclude`: 排除的文件模式，匹配这些模式的文件将被忽略

- `download`: 下载设置
  - `timeout`: 下载超时时间（秒）
  - `retries`: 下载失败时的重试次数
  - `retry_delay`: 重试之间的延迟时间（秒）
  - `user_agent`: HTTP请求使用的用户代理字符串

- `logging`: 日志配置
  - `level`: 日志级别（debug, info, warn, error）
  - `format`: 日志格式（text, json）
  - `output`: 日志输出位置（stdout, file）

- `advanced`: 高级设置
  - `preserve_structure`: 是否保持文件夹结构
  - `create_readme`: 是否为每个下载的仓库创建README文件
  - `validate_checksums`: 是否验证文件校验和

### 文档下载命令

除了原有的下载功能，DownHub 还支持专门的文档下载命令：

```sh
# 下载指定仓库的文档文件
./download docs https://github.com/user/repo

# 指定输出目录
./download docs https://github.com/user/repo -o ./my-docs

# 指定文档路径
./download docs https://github.com/user/repo -d documentation
```

## 📄 License

MIT
