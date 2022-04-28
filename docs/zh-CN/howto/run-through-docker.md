---
title: 通过 Docker 运行
---

_**一魂文档**_ 的 Docker 镜像可以通过 [Docker Hub](https://hub.docker.com/r/unknwon/asouldocs) 或 [GitHub Container Registry](https://github.com/asoul-sig/asouldocs/pkgs/container/asouldocs) 获取。

`latest` 标签指向 [`main` 分支](https://github.com/asoul-sig/asouldocs)上的最新构建版本。

## 注意事项

配置选项 `HTTP_ADDR` 需要被修改为监听 Docker 容器中的网络地址：

```ini
HTTP_ADDR = 0.0.0.0
```

## 启动容器

你需要挂在 `custom` 目录才能使 Docker 容器成功启动（`/app/asouldocs/custom` 为容器内的对应路径）：

```bash
$ docker run \
    --name=asouldocs \
    -p 15555:5555 \
    -v $(pwd)/custom:/app/asouldocs/custom \
    unknwon/asouldocs
```

如果你的文档目标并不是[远程 Git 地址](set-up-documentation.md#文档目标)，则还需要挂载 `docs` 目录（`/app/asouldocs/docs` 为容器内的对应路径）：

```bash
$ docker run \
    --name=asouldocs \
    -p 15555:5555 \
    -v $(pwd)/custom:/app/asouldocs/custom \
    -v $(pwd)/docs:/app/asouldocs/docs \
    unknwon/asouldocs
```
