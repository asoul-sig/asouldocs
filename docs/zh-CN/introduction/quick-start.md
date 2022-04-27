---
title: 快速开始
---

我们将会通过渲染 _**一魂文档**_ 自身的用户文档来对服务器的用法进行了解。

每个服务器都要求一个路径为 `custom/app.ini` 的配置文件和一个包含用户文档的本地目录。

1. 克隆 `asoul-sig/asouldocs` 仓库到本地：

    ```bash
    $ git clone --depth 1 https://github.com/asoul-sig/asouldocs.git
    ```

1. 创建配置文件 `custom/app.ini`：

    ```bash
    $ mkdir custom
    $ touch custom/app.ini
    ```

1. 编辑配置文件 `custom/app.ini` 并指明文档目录 (docs) 的路径：

    ```ini
    [docs]
    TARGET = ./asouldocs/docs
    ```

1. 启动服务器并访问 [http://localhost:5555](http://localhost:5555)：

    ```bash
    $ asouldocs web
    YYYY/MM/DD 00:00:00 [ INFO] ASoulDocs 1.0.0
    YYYY/MM/DD 00:00:00 [ INFO] Listen on http://localhost:5555
    ```

完美！下一步，让我们学习[如何创建文档仓库](../howto/set-up-documentation.md)。
