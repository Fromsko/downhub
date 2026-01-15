# 重新组织下载目录结构

## 概述
重新设计 DownHub 的下载目录结构，使其更加规范化，所有内容默认存放在 `data/` 目录下，并支持通过 YAML 配置自定义基础路径。

## 当前问题
1. 当前配置使用 `./downloads` 作为基础目录，不够规范
2. 文档和源码下载的目录结构不够清晰
3. 缺乏统一的目录组织标准

## 建议的解决方案

### 新的目录结构
```
data/
├── docs/           # 文档下载目录
│   └── gin-gonic/
│       └── gin/
└── source/         # 源码下载目录
    └── gin-gonic/
        └── gin/
```

### 配置结构调整
```yaml
defaults:
  # 基础数据目录（可配置）
  base_data_dir: "data"
  # 文档目录名（相对于 base_data_dir）
  docs_dir: "docs"
  # 源码目录名（相对于 base_data_dir）
  source_dir: "source"
  # 仓库内文档路径
  docs_path: "docs"
  # 其他配置...
```

### 路径生成逻辑
- 文档下载路径：`{base_data_dir}/{docs_dir}/{owner}/{repo}/`
- 源码下载路径：`{base_data_dir}/{source_dir}/{owner}/{repo}/`

## 实现计划
1. 更新配置结构，添加新的目录配置项
2. 修改代码中的路径生成逻辑
3. 更新默认配置文件
4. 确保向后兼容性

## 验收标准
- [ ] 所有下载默认存放在 `data/` 目录下
- [ ] 文档和源码分别存放在 `data/docs/` 和 `data/source/` 下
- [ ] 按仓库 owner/repo 结构组织子目录
- [ ] 可通过 YAML 配置自定义基础目录
- [ ] 保持现有功能完整性
