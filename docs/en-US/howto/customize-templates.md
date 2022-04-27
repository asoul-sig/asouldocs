---
title: Customize templates
---

Every [template file](set-up-documentation.md#template-files) can be customized through overwrite, and custom template files should be placed under the `custom/templates` directory (the path of direcotry can be changed via `[page] CUSTOM_DIRECTORY`).

It is only recommended to customize the `home.html` and `common/navbar.html` files to maintain the maximum backward compatibility.

## UI framework

By default, the JIT of [Tailwind CSS](https://tailwindcss.com/) is included so you do not need any build step for using any of its classes in your custom templates.

However, it is not required to use Tailwind CSS as you have all the freedom to use any UI framework in your `custom/common/head.html` file.

## Localization

The Flamego's [i18n](https://flamego.dev/middleware/i18n.html) middleware is used to handle localization, create locale files under the `custom/locale` directory (the path of direcotry can be changed via `[i18n] CUSTOM_DIRECTORY`) to customize [localization configuration](set-up-documentation.md#localization-configuration).

The syntax for invoking localization function in template files looks like `{{call .Tr "footer::copyright"}}`, where `footer` is the section name and `copyright` is the key name.

## Static assets

Custom static assets should be placed under the `custom/public` directory (the path of direcotry can be changed via `[asset] CUSTOM_DIRECTORY`), then include them in your template file.

For example, suppose you have a custom static asset in the path `custom/public/css/my.css`, then add the following line in your `custom/common/head.tmpl` file:

```go-html-template
<link href="/css/my.css" rel="stylesheet">
```

Notice there is no `public` prefix in the `href` attribute.
