# Project Context

## Purpose
DownHub is a command-line tool for quickly downloading GitHub releases, supporting single repository and batch downloads with proxy support, concurrent downloads, and beautiful progress bar displays.

## Tech Stack
- **Language**: Go 1.24.0
- **CLI Framework**: Cobra (github.com/spf13/cobra)
- **Web Scraping**: Colly (github.com/gocolly/colly/v2)
- **Git Operations**: go-git (github.com/go-git/go-git/v5)
- **Progress Display**: mpb (github.com/vbauerster/mpb/v8)
- **Color Output**: fatih/color (github.com/fatih/color)
- **Configuration**: YAML (gopkg.in/yaml.v3)

## Project Conventions

### Code Style
- Standard Go formatting and conventions
- Clear, descriptive function and variable names
- Proper error handling with meaningful messages
- Modular package structure under `github.com/Fromsko/downhub`

### Architecture Patterns
- CLI-first design with Cobra command structure
- Configuration-driven behavior via YAML files
- Concurrent download management with goroutines
- Clean separation of concerns: CLI, download logic, configuration, and utilities

### Testing Strategy
- Unit tests for core download functionality
- Integration tests for CLI commands
- Mock external dependencies (GitHub API, file system)
- Test coverage for error scenarios and edge cases

### Git Workflow
- Feature branches for new functionality
- Semantic versioning for releases
- Clear commit messages following conventional commits
- PR reviews required for merging

## Domain Context
- GitHub repository structure and release artifacts
- HTTP/HTTPS proxy configuration and usage
- File download patterns with resume capabilities
- Progress tracking and user feedback mechanisms
- YAML configuration management

## Important Constraints
- Must handle GitHub rate limits gracefully
- Support for various file types and sizes
- Network resilience with retry mechanisms
- Cross-platform compatibility (Linux, macOS, Windows)
- Memory efficient for large downloads

## External Dependencies
- GitHub API for repository and release information
- HTTP proxy servers for network access
- File system for local storage and configuration
- Network connectivity for download operations
