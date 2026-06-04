<div align="center">

# mend

**Markdown → HTML, one static binary, zero dependencies.**

[![CI](https://github.com/will/mend/actions/workflows/ci.yml/badge.svg)](https://github.com/will/mend/actions/workflows/ci.yml)
[![Release](https://img.shields.io/github/v/release/will/mend)](https://github.com/will/mend/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/will/mend)](https://goreportcard.com/report/github.com/will/mend)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)
[![Go](https://img.shields.io/badge/go-1.24+-00ADD8)](https://go.dev)

```sh
$ mend -theme monokai README.md > readme.html
```

</div>

---

`mend` is a tiny, fast markdown-to-HTML converter written in Go. It produces
**self-contained HTML files** with embedded CSS and inline-styled syntax
highlighting — no JavaScript, no web fonts, no external requests.

- ✅ CommonMark + GFM (tables, task lists, strikethrough, autolinks)
- ✅ Footnotes, definition lists, smart typography
- ✅ Syntax highlighting via [chroma](https://github.com/alecthomas/chroma) — 200+ languages
- ✅ 4 built-in themes: `github`, `github-dark`, `monokai`, `nord`
- ✅ Watch mode + HTTP live reload for previewing docs locally
- ✅ Single static binary — ships via Homebrew, Scoop, `go install`, or just curl

## Demo

<p align="center">
  <img src="docs/demo.gif" alt="mend demo" width="720">
</p>

Don't have a demo GIF yet? Generate one with [vhs](https://github.com/charmbracelet/vhs):

```sh
vhs docs/demo.tape
```

## Install

### Homebrew

```sh
brew install will/tap/mend
```

### Scoop

```pwsh
scoop bucket add will https://github.com/will/scoop-bucket
scoop install mend
```

### `go install`

```sh
go install github.com/will/mend@latest
```

### Binary download

Grab a release for your platform from [the releases page](https://github.com/will/mend/releases/latest).

### Build from source

```sh
git clone https://github.com/will/mend
cd mend
go build -o mend .
sudo mv mend /usr/local/bin/    # optional
```

## Usage

```sh
mend file.md                    # render to stdout
mend -o out.html file.md        # write to file
cat file.md | mend              # render stdin to stdout
mend -fragment file.md          # emit just the <body> content
mend -theme monokai file.md     # pick a theme
mend -no-highlight file.md      # disable syntax highlighting
mend -no-css file.md            # don't embed the page CSS
mend -w -o out.html file.md     # watch & re-render on change
mend -w -serve :8080 file.md    # watch & serve with live reload
mend --version
```

### Themes

| name          | description                              |
| ------------- | ---------------------------------------- |
| `github`      | Light, default. GitHub-flavoured.         |
| `github-dark` | Dark counterpart.                        |
| `monokai`     | Vibrant dark theme.                      |
| `nord`        | Cool blue-grey nord palette.             |
| `none`        | Disable chroma code styling.             |

## Watch + live reload

```sh
mend -w -serve :8080 README.md
```

Open <http://localhost:8080> and edit `README.md` in your editor. The browser
reloads on save.

## Library use

```go
import "github.com/will/mend/internal/render"

html, err := render.Markdown([]byte("# hi"), render.Options{
    Highlight: true,
    Chroma:    "github",
})
```

## Why `mend`?

> mend (v.) _to repair; to make whole; to finish._

`mend` finishes your markdown — turns a half-written `.md` file into a
ready-to-ship `.html` file. The name is short, low-colision, and verb-shaped:
it does a thing to a file.

## Contributing

PRs welcome. Run the tests:

```sh
go test ./...
```

For releases, push a tag:

```sh
git tag v0.1.0
git push origin v0.1.0
```

`.github/workflows/release.yml` will build and publish binaries for
linux/darwin/windows on amd64 and arm64.

## License

[MIT](LICENSE)
