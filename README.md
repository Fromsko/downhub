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
./download https://github.com/gin-gonic/gin
```

### 使用代理

```sh
./download --proxy http://localhost:7890 https://github.com/gin-gonic/gin
```

### 批量下载

准备一个包含多个仓库地址的文本文件（每行一个）：

```
https://github.com/gin-gonic/gin
https://github.com/labstack/echo
...
```

执行批量下载：

```sh
./download batch -f repo-list.txt
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

## 📄 License

MIT
