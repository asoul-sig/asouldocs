---
title: Use extensions
---

Every _**ASoulDocs**_ server comes with some builtin extensions, it's just few edits away to use them!

## Plausbile

The Plausbile extension integrates with https://plausible.io/, it is disabled by default. Use the following configuration to enable it:

```ini
[extension.plausible]
; Whether to enable this extension
ENABLED = false
; The optional value to be specified for the "data-domain" attribute
DOMAIN =
```
