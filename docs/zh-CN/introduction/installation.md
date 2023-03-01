---
title: 下载安装
---

## 二进制安装

请通过 [GitHub releases](https://github.com/asoul-sig/asouldocs/releases) 页面获取各个版本的二进制文件。

## 源码安装

源码安装要求具有本地的 [Go 语言](https://go.dev/)开发环境，可以通过以下命令检查：

```bash
$ go version
```

最低的 Go 语言版本要求为 **1.19**。

然后通过以下命令构建二进制：

```bash
$ go build
```

最后启动服务器：

```bash
$ ./asouldocs web
```

如需进行本地开发，请阅读[搭建本地开发环境](https://github.com/asoul-sig/asouldocs/blob/main/docs/dev/local_development.md)（英文）。
