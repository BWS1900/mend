package theme

import (
	_ "embed"
	"fmt"
	"strings"
)

const (
	ThemeGitHub     = "github"
	ThemeGitHubDark = "github-dark"
	ThemeMonokai    = "monokai"
	ThemeNord       = "nord"
	ThemeNone       = "none"
)

var themes = []string{ThemeGitHub, ThemeGitHubDark, ThemeMonokai, ThemeNord}

//go:embed style.css
var styleCSS []byte

func Names() string {
	return strings.Join(themes, ", ") + ", " + ThemeNone
}

func IsValid(name string) bool {
	if name == ThemeNone {
		return true
	}
	for _, t := range themes {
		if t == name {
			return true
		}
	}
	return false
}

func CodeStyle(theme string) string {
	switch theme {
	case ThemeGitHubDark:
		return "github-dark"
	case ThemeMonokai:
		return "monokai"
	case ThemeNord:
		return "nord"
	case ThemeNone:
		return ""
	default:
		return "github"
	}
}

type WrapOptions struct {
	Theme    string
	EmbedCSS bool
}

func Wrap(body []byte, opt WrapOptions) []byte {
	if !IsValid(opt.Theme) {
		opt.Theme = ThemeGitHub
	}
	cssBlock := ""
	if opt.EmbedCSS {
		cssBlock = "<style>\n" + string(styleCSS) + "\n</style>"
	}
	out := fmt.Sprintf(`<!doctype html>
<html lang="en">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width,initial-scale=1">
<title>%s</title>
%s
</head>
<body class="theme-%s">
<main class="mend">
%s
</main>
</body>
</html>
`, defaultTitle(opt.Theme), cssBlock, opt.Theme, body)
	return []byte(out)
}

func defaultTitle(theme string) string {
	return "mend - " + theme
}
