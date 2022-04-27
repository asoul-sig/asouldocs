---
title: 配置反向代理
---


## Caddy 2

通过 [Caddy](https://caddyserver.com/) 可以获得永久免费的 HTTPS：

```caddyfile
{
        http_port 80
        https_port 443
}

asouldocs.dev {
        reverse_proxy * localhost:5555
}
```

## NGINX

如下展示了 NGINX 的配置，但具体值需要根据实际情况修改：

```nginx
server {
    listen 80;
    server_name asouldocs.dev;

    location / {
        proxy_pass http://localhost:5555;
    }
}
```
