---
title: 创建文档仓库
---

_**一魂文档**_ 的文档仓库包含以下四个要素：

1. 配置文件
1. 模板文件
1. 本地化文件
1. 用户文档

## 配置文件

每一个文档项目都需要一个配置文件，默认路径为 `custom/app.ini`。

每个服务器都会内置一个[默认配置文件](https://github.com/asoul-sig/asouldocs/blob/main/conf/app.ini)为各项配置提供了默认选项，从而使你的文档项目只需要进行几项必要的修改。

命令行参数 `--config`（或 `-c`）可以用于指定飞默认路径的配置文件。

## 模板文件

模板文件使用 [Go HTML 模板](https://pkg.go.dev/html/template)语法渲染各个页面。

每个服务器都会内置一批可直接使用的[默认模板文件](https://github.com/asoul-sig/asouldocs/tree/main/templates)，但大多数情况下你可能都需要修改 `home.html` 和 `common/navbar.html` 这两个模板文件用于展示与你项目相关的内容。

你可以通过在 `custom/templates` 目录下放置对应名称的模板文件来覆写内置模板文件。

例如，下面展示了如何将 `common/navbar.html` 中的 GitHub 字样替换为一个笑脸的图案：

```go-html-template {hl_lines=["7-10"]}
<ul class="flex items-center space-x-5">
  <li>
    <a class="hover:text-sky-500 dark:hover:text-sky-400" href="{{.Page.DocsBasePath}}">{{call .Tr "navbar::docs"}}</a>
  </li>
  <li>
    <a class="hover:text-sky-500 dark:hover:text-sky-400" href="https://github.com/asoul-sig/asouldocs">
      <!-- Heroicon: outline/emoji-happy -->
      <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
        <path stroke-linecap="round" stroke-linejoin="round" d="M14.828 14.828a4 4 0 01-5.656 0M9 10h.01M15 10h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
      </svg>
    </a>
  </li>
</ul>
```

如需了解更多，请阅读[如何自定义模板](customize-templates.md)。

## 文档结构

单语言和多语言的文档结构是一致的，在文档的根目录中需要放置一个 `toc.ini` 文件，其子目录的命名也需要遵循 [IETF BCP 47 language tags](https://en.wikipedia.org/wiki/IETF_language_tag) 规范，如 en-US、zh-CN。

如下所示：

```
├── en-US
│   ├── howto
│   │   ├── README.md
│   │   └── set_up_documentation.md
│   └── introduction
│       ├── README.md
│       ├── installation.md
│       └── quick_start.md
├── zh-CN
│   ├── howto
│   │   ├── README.md
│   │   └── set_up_documentation.md
│   └── introduction
│       ├── README.md
│       ├── installation.md
│       └── quick_start.md
└── toc.ini
```

`toc.ini` 文件用于描述各个文档将以如何顺序和组织结构展现在站点上，下面展示了与上例所匹配的 `toc.ini` 文件内容：

```ini
-: introduction
-: howto

[introduction]
-: README
-: installation
-: quick_start

[howto]
-: README
-: set_up_documentation
```

默认分区用于定义文档目录的展现顺序，各个值与目录名称一一对应：

```ini
-: introduction
-: howto
```

各个目录也分别需要一个同名的分区：

```ini
[introduction]
...

[howto]
...
```

与目录同名的分区用于定义文件的展现顺序，各个值与文件名称一一对应（但不包括 `.md` 后缀名）：

```ini
[introduction]
-: README
-: installation
-: quick_start
```

其它事项：

1. 目前仅支持单层目录结构
1. 每个目录都必须包含一个 `README.md` 文件，并且该文件至少需要包含[前置配置](write-document.md#前置配置)部分的内容。

### 本地化配置

每个服务器都会根据[默认配置](https://github.com/asoul-sig/asouldocs/blob/39b59c4159e4a2b0e0a290c79f85c46a3e1faf0b/conf/app.ini#L26-L30)认为需要为用户提供英语 (en-US) 和简体中文 (zh-CN) 两种语言版本的文档：

```ini
[i18n]
; 支持的文档语言列表
LANGUAGES = en-US,zh-CN
; 各个语言用户友好的名称
NAMES = English,简体中文
```

`LANGUAGES` 值中的第一个语言会被作为默认语言，当服务器无法找到（根据浏览器的 `Accept-Language` 请求头识别）偏好语言的文档时则会显示默认语言的内容。

如果你的文档仅包含简体中文，则需要进行以下配置的修改：

```ini
[i18n]
; 支持的文档语言列表
LANGUAGES = zh-CN
; 各个语言用户友好的名称
NAMES = 简体中文
```

## 文档目标

文档目标可以是本地目录的路径或远程 Git 地址。

以下配置用于本地目录的文档：

```ini
[docs]
TYPE = local
TARGET = ./docs
```

以下配置用于远程 Git 地址的文档：

```ini
[docs]
TYPE = remote
TARGET = https://github.com/asoul-sig/asouldocs.git
```

如果文档目录存在于子目录，则可以通过 `TARGET_DIR` 选项指明：

```ini
[docs]
TYPE = remote
TARGET = https://github.com/asoul-sig/asouldocs.git
TARGET_DIR = docs
```

## 编辑本页内容

如果你希望用户帮助改进文档，则可以提供一个快捷链接方便用户快速定位到当前所浏览文档的源文件：

```ini
[docs]
; 编辑链接的格式化字符串，留空表示禁用该功能
EDIT_PAGE_LINK_FORMAT = https://github.com/asoul-sig/asouldocs/blob/main/docs/{blob}
```
