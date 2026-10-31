# PicClaw 快速开始（简体中文）

本指南帮助你快速运行 PicClaw。命令与标识符保持英文原样，请勿翻译。

## 环境要求

- Go 1.24.0 或更高版本
- Git

## 获取代码

git clone https://github.com/sipeed/picoclaw.git
cd picoclaw

## 构建

go build ./...

## 运行

go run .

## 常见问题

- 若依赖下载失败，请检查网络代理设置。
- 若编译失败，请确认 Go 版本满足 `go.mod` 要求。

## 未验证项

- Windows 环境下的完整端到端流程尚未在本机验证。