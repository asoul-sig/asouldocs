---
title: Set up documentation
---

There are four essential components of a _**ASoulDocs**_ documentation project:

1. Configuration file
1. Template files
1. Locale files
1. Actual documents

## Configuration file

The configuration file is required for every documentation project. By default, it is expected be available at `custom/app.ini`.

Every server comes with a [builtin configuration file](https://github.com/asoul-sig/asouldocs/blob/main/conf/app.ini) as default options. Therefore, you only need to overwrite few options in your own configuration file.

The `--config` (or `-c`) CLI flag can be used to specify a configuration file that is not located in the default location.

## Template files

Template files are [Go HTML templates](https://pkg.go.dev/html/template) used to render different pages served by the server.

Every server comes with a set of [builtin template files](https://github.com/asoul-sig/asouldocs/tree/main/templates) that works out-of-the-box. However, the content of builtin template files are probably not what you would want for your documentation in most cases, at least for `home.html` and `common/navbar.html`.

Luckily, you can overwrite these template files by placing your template files with same file name under the `custom/templates` directory.

For example, you can replace the "GitHub" in navbar to be a happy face by overwriting the `common/navbar.html`:

```go-html-template {hl_lines=["7-10"]}
<ul class="flex items-center space-x-5">
  <li>
    <a class="hover:text-sky-500 dark:hover:text-sky-400" href="{{.Page.DocsBasePath}}">{{call .Tr "navbar::docs"}}</a>
  </li>
  <li>
    <a class="hover:text-sky-500 dark:hover:text-sky-400" href="https://github.com/asoul-sig/asouldocs">
      <!-- Heroicon: outline/emoji-happy -->
      <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
        <path stroke-linecap="round" stroke-linejoin="round" d="M14.828 14.828a4 4 0 01-5.656 0M9 10h.01M15 10h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
      </svg>
    </a>
  </li>
</ul>
```

You can read more about how to [customize templates](customize-templates.md).

## Document hierarchy

Whether you want to serve your documentation in multiple languages or just one language, the hierarchy is the same. In the root directory of your documentation, you need a `toc.ini` file and subdirectories using [IETF BCP 47 language tags](https://en.wikipedia.org/wiki/IETF_language_tag), e.g. "en-US", "zh-CN".

Here is an example hierarchy:

```
├── en-US
│   ├── howto
│   │   ├── README.md
│   │   └── set_up_documentation.md
│   └── introduction
│       ├── README.md
│       ├── installation.md
│       └── quick_start.md
├── zh-CN
│   ├── howto
│   │   ├── README.md
│   │   └── set_up_documentation.md
│   └── introduction
│       ├── README.md
│       ├── installation.md
│       └── quick_start.md
└── toc.ini
```

The `toc.ini` is used to define how exactly these documents should look like on the site (e.g. how they are ordered). The following example is a corresponding `toc.ini` to the above hierarchy:

```ini
-: introduction
-: howto

[introduction]
-: README
-: installation
-: quick_start

[howto]
-: README
-: set_up_documentation
```

In the default section, document directories are defined in the exact order, and these names are corresponding to the directories' title:

```ini
-: introduction
-: howto
```

Then there are sections for each directory:

```ini
[introduction]
...

[howto]
...
```

Within each section, files are defined in the exact order, and these names are corresponding to the files' name in the document directory (but without the `.md` extension):

```ini
[introduction]
-: README
-: installation
-: quick_start
```

Other notes:

1. Only single-level directories are supported.
1. Every document directory must have a `README.md` file, to at least define its name through the [frontmatter](write-document.md#frontmatter).

### Localization configuration

By default, the server assumes to have documentation both in English (`en-US`) and Simplified Chinese (`zh-CN`), as in the [default configuration](https://github.com/asoul-sig/asouldocs/blob/39b59c4159e4a2b0e0a290c79f85c46a3e1faf0b/conf/app.ini#L26-L30):

```ini
[i18n]
; The list of languages that is supported
LANGUAGES = en-US,zh-CN
; The list of user-friendly names of languages
NAMES = English,简体中文
```

The first language in the `LANGUAGES` is considered as the default language, and the server shows its content if the prefered language (from browser's `Accept-Language` request header) does not exists, or the particular document is not available in the prefered language (but available in the default language).

If you are just writing documentation in English, you would need to overwrite the configuration as follows:

```ini
[i18n]
; The list of languages that is supported
LANGUAGES = en-US
; The list of user-friendly names of languages
NAMES = English
```

## Target

The target of documents can be either a local directory or a remote Git address.

To use a local directory:

```ini
[docs]
TYPE = local
TARGET = ./docs
```

To use a remote Git address:

```ini
[docs]
TYPE = remote
TARGET = https://github.com/asoul-sig/asouldocs.git
```

If documents are residing in a subdirectory of the target, use `TARGET_DIR` as follows:

```ini
[docs]
TYPE = remote
TARGET = https://github.com/asoul-sig/asouldocs.git
TARGET_DIR = docs
```

## Link to edit page

You probably want to welcome small contributions from the community if your documentation repository is open sourced. You can navigate users directly to edit the page on the code host:

```ini
[docs]
; The format to construct a edit page link, leave it empty to disable
EDIT_PAGE_LINK_FORMAT = https://github.com/asoul-sig/asouldocs/blob/main/docs/{blob}
```
