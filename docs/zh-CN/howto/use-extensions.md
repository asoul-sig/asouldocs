---
title: 使用扩展
---

每个 _**一魂文档**_ 服务器都提供了一些极其便捷的内置扩展。

## Plausbile

Plausbile 扩展用于集成 https://plausible.io/，可以通过如下配置启用：

```ini
[extension.plausible]
ENABLED = true
; 用于指定 data-domain 属性的可选值
DOMAIN =
```

## Google Analytics

Google Analytics 扩展用于集成 [Google Analytics 4](https://developers.google.com/analytics/devguides/collection/ga4)，可以通过如下配置启用：

```ini
[extension.google_analytics]
ENABLED = true
; 资产所对应的 Measurement ID
MEASUREMENT_ID = G-VXXXYYYYZZ
```

## Disqus

Disqus 扩展用于集成 [Disqus](https://disqus.com/)，可以通过如下配置启用：

```ini
[extension.disqus]
ENABLED = true
; 站点的 shortname
SHORTNAME = ellien
```

## utterances

utterances 扩展用于集成 [utterances](https://utteranc.es/)，可以通过如下配置启用：

```ini
[extension.utterances]
ENABLED = true
; GitHub 仓库名称
REPO = owner/repo
; Issue 映射模式
ISSUE_TERM = pathname
; Issue 标签
LABEL = utterances
; 组件主题
THEME = github-light
```
