# fylr-plugin-server-pdf

> [!WARNING]
> **Deprecated.** [PDF Creator](https://github.com/programmfabrik/fylr-plugin-pdf-creator)
> renders PDFs itself since v1.1.0 and no longer uses this plugin, which was its only
> consumer in the product. Nothing in fylr calls this plugin any more.
>
> It is **not** removed by the 6.35 disk-to-url migration — that only touches plugins fylr
> shipped — so if you installed it for PDF Creator you can uninstall it. Keep it only if
> something of your own posts to its `html2pdf` endpoint.
>
> No new features are planned. The renderer lives on inside PDF Creator.

Provides fylr with a general HTML-to-PDF service: it exposes an `html2pdf` extension endpoint
(`plugin/extension/server-pdf/html2pdf`) that renders HTML posted to it into a PDF, using a
headless Chrome/Chromium on the exec server.

> **Not needed for PDF Creator.** The current
> [PDF Creator](https://github.com/programmfabrik/fylr-plugin-pdf-creator) ships its own
> renderer and does not use this plugin. Install server-pdf when you want the HTML-to-PDF
> endpoint for your own integration; the easydb-era PDF Creator used to depend on it.

* Check that you exec server is prepared to run plugin binaries, like so:

```yaml
    execserver+:
      ...
      env+:
        # The environment variable has to point to the Chrome binary. Without configuration "chromium"
        # is expected in the PATH.
        # Chrome / chromium is started using the parameters:
        #  --headless --disable-gpu --no-sandbox --remote-debugging-port=0
        # These parameters are hardcoded and cannot be changed at this point.
        - SERVER_PDF_CHROME=/Applications/Google Chrome.app/Contents/MacOS/Google Chrome
```

* By default the server-pdf plugin expect the binary `chromium` in the exec server. If you prefer to use "Chrome" or want to provide an absolute path to the binary, use environment.

* Install this Plugin into fylr (URL, ZIP or disk mode)
* The [easydb-pdf-creator-plugin](https://github.com/programmfabrik/easydb-pdf-creator-plugin) will auto recognize this plugin and call the fylr server pdf creator automatically.
