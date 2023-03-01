---
title: Installation
---

## From binary

Release binaries are available on [GitHub releases](https://github.com/asoul-sig/asouldocs/releases).

## From source code

Install from source code requires you having a working local environment of [Go](https://go.dev/).

Use the following command to check:

```bash
$ go version
```

The minimum requirement version of Go is **1.19**.

Then build the binary:

```bash
$ go build
```

Finally, start the server:

```bash
$ ./asouldocs web
```

Please refer to [Set up your development environment](https://github.com/asoul-sig/asouldocs/blob/main/docs/dev/local_development.md) for local development guide.
