# Search

Search is a plugin for [Kite](https://github.com/kite-plus/kite) that lets
readers search a site's posts and pages. It needs no server and no service:
the index is written when the site is built, and searched in the reader's
browser.

Readers open it with the button in the corner of the page, or by pressing
`/` or `Ctrl K` (`⌘ K` on a Mac). Words match anywhere in a title, a tag or
the text, which also works for Chinese and Japanese, where words are not
spaced.

## Using it

Drop the zip of a release on the upload tile under **Plugins** in Kite's
studio, or add it from the command line inside a site:

```sh
kite plugin add search-0.1.0.zip
kite plugin enable search
```

The index holds every post and page, whole. On a large site, turning off
**Search whole posts** keeps it to titles, tags and summaries. A post stays
out of it with `search: false` in its front matter.

The plugin asks for Kite 1.0 or later.

## In a theme

A theme with a search button of its own marks it, and the corner button is
left out:

```html
<button type="button" data-kite-search>Search</button>
```

## How it works

Once the site is built, the plugin's module, `plugin.wasm`, is handed every
page and writes the index to `plugins/search/index.json`. `search.js` loads
it the first time a reader opens the search.

## Developing it

The module is Go, built for WebAssembly with Go 1.24 or later:

```sh
make test     # the index, on this computer
make build    # plugin.wasm
kite plugin verify .
```

## Releasing

`make zip` builds the module and packs `dist/search-<version>.zip`, the
package without the module's source, which the studio installs as it is.
Tag the release and attach the zip.

## License

[Apache License 2.0](LICENSE).
