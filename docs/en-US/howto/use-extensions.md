---
title: Use extensions
---

Every _**ASoulDocs**_ server comes with some builtin extensions, it's just few edits away to use them!

## Plausbile

The Plausbile extension integrates with https://plausible.io/, it is disabled by default. Use the following configuration to enable it:

```ini
[extension.plausible]
ENABLED = true
; The optional value to be specified for the "data-domain" attribute
DOMAIN =
```

## Google Analytics

The Google Analytics extension integrates with [Google Analytics 4](https://developers.google.com/analytics/devguides/collection/ga4), it is disabled by default. Use the following configuration to enable it:

```ini
[extension.google_analytics]
ENABLED = true
; The measurement ID of your property
MEASUREMENT_ID = G-VXXXYYYYZZ
```

## Disqus

The Disqus extension integrates with [Disqus](https://disqus.com/), it is disabled by default. Use the following configuration to enable it:

```ini
[extension.disqus]
ENABLED = true
; The shortname of your site
SHORTNAME = ellien
```

## utterances

The utterances extension integrates with [utterances](https://utteranc.es/), it is disabled by default. Use the following configuration to enable it:

```ini
[extension.utterances]
ENABLED = true
; The GitHub repository
REPO = owner/repo
; The issue mapping pattern
ISSUE_TERM = pathname
; The issue label for comments
LABEL = utterances
; The theme of the component
THEME = github-light
```
