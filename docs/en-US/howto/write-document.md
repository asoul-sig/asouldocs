---
title: Write document
---

Every document is a [Markdown](https://www.markdownguide.org/) file, most of features from [GitHub Flavored Markdown](https://github.github.com/gfm/) are supported through [yuin/goldmark](https://github.com/yuin/goldmark).

Here is an simple example:

```markdown
---
title: Introduction
---

_**ASoulDocs**_ is a stupid web server for multilingual documentation.
```

### Frontmatter

The frontmatter is a block of YAML snippet, `---` are used to both indicate the start and the end of the snippet.

Here is an full example of supported fields:

```yaml
title: The title of the document
previous:
  title: The title of the previous page
  link: the relative path to the page
next:
  title: The title of the next page
  link: the relative path to the page
```

- Only the `title` field is required, others all have reasonable default.
- The `link` syntax of both `previous` and `next` sections is exactly same as described in [Links and images](#links-and-images).

### Links and images

Links to other documents or images just works like you would do in any editor (e.g. VSCode):

- Link to a document under the same directory: `[Customize templates](customize-templates.md)`
- Link to the directory page: `[How to](README.md)`
- Link to a document in another directory: `[Quick start](../introduction/quick-start.md)`
- Link to an image: `![](../../assets/workflow.png)`

### Code blocks

The [alecthomas/chroma](https://github.com/alecthomas/chroma) is used to syntax highlighting your code blocks.

Use name from its [supported languages](https://github.com/alecthomas/chroma#supported-languages) to specify the language of the code block, be sure to replace whitespaces with hyphenes (`-`) in the language name, e.g. use `go-html-template` not `go html template` (names are case insensitive).

To enable line numbers and highlighting lines:

```markdown
...go-html-template {linenos=true, hl_lines=["7-10", 12]}
```

### Render caching

Each documents is reloaded and re-rendered for every request if the server is running with `dev` envrionment, as defined in the [configuration file](set-up-documentation.md#configuration-file):

```ini
ENV = dev
```

Set `ENV = prod` to enable render caching when deploy to production.
