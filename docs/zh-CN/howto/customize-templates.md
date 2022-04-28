---
title: 自定义模板
---

每个[模板文件](set-up-documentation.md#模板文件)都可以通过覆写实现自定义，自定义模板文件需要被放置在 `custom/templates` 目录下（可以通过配置选项 `[page] CUSTOM_DIRECTORY` 指定其它目录）。

一般只建议自定义 `home.html` 和 `common/navbar.html` 这两个模板文件来实现最大程度上的向后兼容。

## UI 框架

页面样式模式使用 [Tailwind CSS](https://tailwindcss.com/) 的 JIT 编译器进行渲染，因此你可以在自定义模板直接使用该框架的所有样式。

当然了，并不强求所有服务器都使用 Tailwind CSS 作为 UI 框架，只要在 `custom/common/head.html` 中导入你的自选资源即可。

## 本地化

服务器使用 Flamego 的 [i18n](https://flamego.cn/middleware/i18n.html) 中间件实现本地化，进行[本地化配置](set-up-documentation.md#本地化配置)的本地化文件需要被放置在 `custom/locale` 目录下（可以通过配置选项 `[i18n] CUSTOM_DIRECTORY` 指定其它目录）。

在模板文件中调用本地化函数的语法为 `{{call .Tr "footer::copyright"}}`，其中, `footer` 为分区名，`copyright` 为键名。

## 静态资源

自定义静态资源需要被放置在 `custom/public` 目录下（可以通过配置选项 `[asset] CUSTOM_DIRECTORY` 指定其它目录）并在模板中导入。

假设你有一个路径为 `custom/public/css/my.css` 的自定义静态资源，并在 `custom/common/head.tmpl` 文件中添加如下内容：

```go-html-template
<link href="/css/my.css" rel="stylesheet">
```

注意 `href` 属性并不包含 `public` 前缀。
