---
title: 通过 Webhook 同步
---

每次更新文档都要登录服务器拉取和重启总是令人不爽，而 _**一魂文档**_ 仅需通过以下两步就能实现文档的自动更新：

1. [使用远程 Git 地址作为文档仓库](set-up-documentation.md#文档目标)
1. 配置在每次推送后都发送 Webhook 到服务器

Webhook 的请求路径为 `/webhook` 且接受任何 HTTP 方法，因此一个简单的 `curl` 命令就可以实现：

```bash
$ curl http://localhost:5555/webhook
```

几乎所有的代码平台都会提供内置的 Webhook 功能，如可以通过 **Settings > Webhooks** 页面中的 **Add webhook** 按钮为 GitHub 仓库配置 Webhook：

![GitHub webhook](../../assets/github-webhook.jpg)
