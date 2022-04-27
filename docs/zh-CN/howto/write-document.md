---
title: 编写文档
---

每个文档都需要使用 [Markdown](https://www.markdownguide.org/) 语法编写，[GitHub Flavored Markdown](https://github.github.com/gfm/) 规范当中的大部分功能也已通过 [yuin/goldmark](https://github.com/yuin/goldmark) 实现。

如下所示：

```markdown
---
title: 项目介绍
---

_**一魂文档**_ 是一款支持多语言的 Web 文档服务器。
```

### 前置配置

前置配置是指每个文档的开头部分使用 `---` 包括的 YAML 代码块。

下面展示了目前支持的所有前置配置字段：

```yaml
title: 文档标题
previous:
  title: 前个页面的标题
  link: 前个页面的相对路径
next:
  title: 后个页面的标题
  link: 后个页面的相对路径
```

- 除了 `title` 以为均为可选字段
- `previous` 和 `next` 下 `link` 字段的语法和[链接与图片](#链接与图片)一致。

### 链接与图片

指向其它文档或图片的链接与任何编辑器中的语法并无二致（如 VSCode）：

- 指向同个目录下的文档：`[Customize templates](customize-templates.md)`
- 指向其它目录：`[How to](README.md)`
- 指向不同目录下的文档：`[Quick start](../introduction/quick-start.md)`
- 指向图片：`![](../../assets/workflow.png)`

### 代码块

代码块的高亮使用 [alecthomas/chroma](https://github.com/alecthomas/chroma) 实现，因此需要使用其[支持语言](https://github.com/alecthomas/chroma#supported-languages)（英文）列表中的名称来指定代码块的语言，名称中的空格需要使用横线 (`-`) 替代，如使用 `go-html-template` 而不是 `go html template`（大小写不敏感）。

代码块还支持行数显示和多行高亮：

```markdown
...go-html-template {linenos=true, hl_lines=["7-10", 12]}
```

### 渲染缓存

当服务器使用[默认配置](set-up-documentation.md#configuration-file)运行时（`dev` 环境），每个请求都会重新加载并渲染文档页面：

```ini
ENV = dev
```

在生产环境需要通过配置 `ENV = prod` 启用渲染缓存。
