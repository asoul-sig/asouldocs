---
title: Quick start
---

Let's boot up a server serving the documentation of _**ASoulDocs**_ its own!

At bare minimum, the server requires a `custom/app.ini` file and a local directory where source files of documentation are located.

1. Clone the `asoul-sig/asouldocs` repository locally:

    ```bash
    $ git clone --depth 1 https://github.com/asoul-sig/asouldocs.git
    ```

1. Create a `custom/app.ini` file:

    ```bash
    $ mkdir custom
    $ touch custom/app.ini
    ```

1. Edit the `custom/app.ini` to tell the server where "docs" directory is located:

    ```ini
    [docs]
    TARGET = ./asouldocs/docs
    ```

1. Start the server and visit [http://localhost:5555](http://localhost:5555):

    ```bash
    $ asouldocs web
    YYYY/MM/DD 00:00:00 [ INFO] ASoulDocs 1.0.0
    YYYY/MM/DD 00:00:00 [ INFO] Listen on http://localhost:5555
    ```

Great! Let’s move on to [how to set up documentation](../howto/set-up-documentation.md).
