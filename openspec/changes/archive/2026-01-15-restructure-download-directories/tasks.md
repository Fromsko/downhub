# 任务清单

## 配置结构更新
- [x] 更新 `config/config.go` 中的 `Defaults` 结构体
- [x] 添加 `BaseDataDir`, `DocsDir`, `SourceDir` 字段
- [x] 更新 `GetDefaultConfig()` 函数
- [x] 移除旧的 `BaseDir` 字段

## YAML 配置文件更新
- [x] 更新 `downhub.yaml` 配置文件
- [x] 设置默认的 `base_data_dir: "data"`
- [x] 设置默认的 `docs_dir: "docs"`
- [x] 设置默认的 `source_dir: "source"`

## 代码逻辑更新
- [x] 更新 `cmd/cmd.go` 中的路径生成逻辑
- [x] 更新 `handler/handler.go` 中的源码下载路径
- [x] 更新 `handler/docs.go` 中的文档下载路径
- [x] 确保所有路径都使用新的配置结构

## 测试验证
- [x] 测试单个文档下载功能
- [x] 验证目录结构符合规范 (`data/docs/gin-gonic/gin`)
- [x] 确保配置文件正确生效
- [x] 验证文件过滤器正常工作
