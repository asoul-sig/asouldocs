---
title: Introduction
previous:
  title: Home
  path: ../../
---

_**ASoulDocs**_ is a stupid web server for multilingual documentation.

## Motivation

Project documentation tools is an already-crowded place yet not being both mature and affordable for individual especially OSS developers. Countless static site generators, documentation servers, SaaS products, unfortunately, we are not happy with any of them. More importantly, we love and are capable of hacking on this area.

It has been years we struggled with picking one of them, that's basically why _**ASoulDocs**_ was created (previously "Peach" pre-1.0).

The following table illustrates the features (that we care) comparisons between _**ASoulDocs**_ and other existing tools (that we investigated, concepts may be different from what we understand):

|Name/Feature                 |_**ASoulDocs**_|[Mkdocs](https://www.mkdocs.org/)|[Hugo](https://gohugo.io/)|[VuePress](https://v2.vuepress.vuejs.org/)/[VitePress](https://vitepress.vuejs.org/)|[GitBook](https://www.gitbook.com/)|
|:---------------------------:|:-------------:|:----:|:--:|:----------------:|:----:|
|Self-hosted                  | ✅ | ✅ | ✅ | ✅ | ❌ |
|Multilingual<sup>1</sup>     | ✅ | ✅ | ✅ | ✅ | ❌ |
|Builtin push-to-sync         | ✅ | ❌ | ❌ | ❌ | ✅ |
|DocSearch                    | 🎯 | ❌ | ✅ | ✅ | ❌ |
|Builtin search               | 🎯 | ✅ | ❌ | ✅ | ✅ |
|Commenting system            | ✅ | ❌ | ✅ | ❌ | ❌ |
|Versionable                  | 🎯 | ❌ | ❌ | ❌ | ❌ |
|Protected resources          | 🎯 | ❌ | ❌ | ❌ | ❌ |
|Dark mode                    | ✅ | ❌ | ✅ | ✅ | ❌ |
|Customizable<sup>2</sup>     | ✅ | ❌ | ✅ | ❌ | ❌ |
|Language fallback<sup>3</sup>| ✅ | ❌ | ❌ | ❌ | ❌ |

- <sup>1</sup>: None of others support multilingual without changing the URL, which to us is a bizarre because we have to share different URLs for different groups of users.
- <sup>2</sup>: In such way that visitors couldn't recognize what is powering the site behind the scene.
- <sup>3</sup>: When a page does not exist in the preferred language, fallback to show the version from the default language.
- 🎯: Features that are on the roadmap.

## History

The project was initially named as "Peach Docs" in pre-1.0 releases, which is now branded as _**ASoulDocs**_ starting from the 1.0 release.

The tech stack has evolved since 2015, [Macaron](https://go-macaron.com) and [Semantic UI](https://semantic-ui.com/) was the new and hot things, and the latest golden partners are [Flamego](https://flamego.dev) and [Tailwind CSS](https://tailwindcss.com/).

The project is now part of [A-SOUL SIG](https://github.com/asoul-sig) (previously "github.com/peachdocs"), consists a group of A-SOUL fans.

## OK, then what?

[Install the server](installation.md) or go ahead to [Quick start](quick-start.md)!
