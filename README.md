# fylr-plugin-server-pdf

Plugin to provide HTML to PDF functionality for fylr.

* Check that you exec server is prepared to run plugin binaries, like so:

```yaml
    execserver:
      ...
      waitgroups:
        a:
          processes: 2
      services:
        # this service allows to execute arbitrary binaries
        exec:
          # choose a waitgroup which can take heavy tasks. Chromium on
          # Linux needs about 500 MB to operate. Each parallel produce PDF
          # will take that memory.
          waitgroup: c
          commands:
            exec:
              env:
                # defaults to "chromium"
                - SERVER_PDF_CHROME=chromium # change to Chrome or absolute paths

```

* By default the server-pdf plugin expect the binary `chromium` in the exec server. If you prefer to use "Chrome" or want to provide an absolute path to the binary, use environment.

* Install this Plugin into fylr (URL, ZIP or disk mode)
* The [easydb-pdf-creator-plugin](https://github.com/programmfabrik/easydb-pdf-creator-plugin) will auto recognize this plugin and call the fylr server pdf creator automatically.
