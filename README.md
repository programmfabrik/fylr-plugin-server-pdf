# fylr-plugin-server-pdf

Plugin to provide HTML to PDF functionality for fylr.

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
