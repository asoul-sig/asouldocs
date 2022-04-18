# Set up your development environment

_**ASoulDocs**_ is written in [Go](https://golang.org/), please take [A Tour of Go](https://tour.golang.org/) if you haven't done so!

## Outline

- [Environment](#environment)
- [Step 1: Install dependencies](#step-1-install-dependencies)
- [Step 2: Get the code](#step-2-get-the-code)
- [Step 3: Start the server](#step-3-start-the-server)
- [Other nice things](#other-nice-things)

## Environment

_**ASoulDocs**_ is built and runs as a single binary and meant to be cross platform. Therefore, you should be able to develop _**ASoulDocs**_ in any major platforms you prefer. However, this guide will focus on macOS only.

## Step 1: Install dependencies

_**ASoulDocs**_ has the following dependencies:

- [Git](https://git-scm.com/book/en/v2/Getting-Started-Installing-Git) (v2 or higher)
- [Go](https://go.dev/doc/install) (v1.17 or higher)
- [Task](https://taskfile.dev/) (v3)

1. Install [Homebrew](https://brew.sh/).
1. Install dependencies:

    ```bash
    brew install go git go-task/tap/go-task
    ```

## Step 2: Get the code

Generally, you don't need a full clone, so set `--depth` to `10`:

```bash
git clone --depth 10 https://github.com/asoul-sig/asouldocs.git

# or

git clone --depth 10 git@github.com:asoul-sig/asouldocs.git
```

**NOTE** The repository has Go modules enabled, please clone to somewhere outside of your `$GOPATH`.

## Step 3: Start the server

The following command will start the web server and automatically recompile and restart the server if any watched files changed:

```bash
task web --watch
```

## Other nice things

### Load HTML templates and static files from disk

When you are actively working on HTML templates and static files during development, you would want to ensure the following configuration to avoid recompiling and restarting _**ASoulDocs**_ every time you make a change to files under `templates/` and `public/` directories:

```ini
ENV = dev
```
