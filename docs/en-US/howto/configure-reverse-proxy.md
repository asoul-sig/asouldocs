---
title: Configure reverse proxy
---

## Caddy 2

You get HTTPS for free and forever just by using [Caddy](https://caddyserver.com/):

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

Here is an example of NGINX config section, but values can be different based on your situation:

```nginx
server {
    listen 80;
    server_name asouldocs.dev;

    location / {
        proxy_pass http://localhost:5555;
    }
}
```
