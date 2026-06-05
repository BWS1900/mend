package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/BWS1900/mend/internal/render"
	"github.com/BWS1900/mend/internal/serve"
	"github.com/BWS1900/mend/internal/theme"
	"github.com/BWS1900/mend/internal/version"
	"github.com/BWS1900/mend/internal/watch"
)

type config struct {
	output      string
	themeName   string
	noHighlight bool
	noCSS       bool
	fragment    bool
	watchFlag   bool
	serveAddr   string
	printVer    bool
}

func main() {
	cfg := parseFlags()

	if cfg.printVer {
		fmt.Println(version.String())
		return
	}

	args := flag.Args()

	if cfg.watchFlag {
		if len(args) == 0 {
			fmt.Fprintln(os.Stderr, "mend: -w requires an input file")
			os.Exit(2)
		}
		runWatch(args[0], cfg)
		return
	}

	if len(args) == 0 {
		stat, _ := os.Stdin.Stat()
		if (stat.Mode() & os.ModeCharDevice) != 0 {
			usage()
			os.Exit(2)
		}
	}

	if err := runOnce(args, cfg); err != nil {
		die(err)
	}
}

func parseFlags() config {
	cfg := config{}
	flag.StringVar(&cfg.output, "o", "", "write output to file (default: stdout)")
	flag.StringVar(&cfg.themeName, "theme", "github", "theme: "+theme.Names())
	flag.BoolVar(&cfg.noHighlight, "no-highlight", false, "disable syntax highlighting in code blocks")
	flag.BoolVar(&cfg.noCSS, "no-css", false, "do not embed page CSS (fragment only)")
	flag.BoolVar(&cfg.fragment, "fragment", false, "emit HTML fragment (no <html>/<head>/<body>)")
	flag.BoolVar(&cfg.watchFlag, "w", false, "watch input file and re-render on change")
	flag.StringVar(&cfg.serveAddr, "serve", "", "serve rendered HTML on addr (e.g. :8080) with live reload")
	flag.BoolVar(&cfg.printVer, "version", false, "print version and exit")
	flag.Usage = usage
	flag.Parse()
	return cfg
}

func runOnce(args []string, cfg config) error {
	var (
		input []byte
		err   error
	)
	if len(args) > 0 {
		input, err = os.ReadFile(args[0])
	} else {
		input, err = io.ReadAll(os.Stdin)
	}
	if err != nil {
		return err
	}
	out, err := assemble(input, cfg)
	if err != nil {
		return err
	}
	if cfg.output != "" {
		return os.WriteFile(cfg.output, out, 0o644)
	}
	_, err = os.Stdout.Write(out)
	return err
}

func assemble(input []byte, cfg config) ([]byte, error) {
	body, err := render.Markdown(input, render.Options{
		Highlight: !cfg.noHighlight,
		Chroma:    theme.CodeStyle(cfg.themeName),
	})
	if err != nil {
		return nil, err
	}
	if cfg.fragment {
		return body, nil
	}
	return theme.Wrap(body, theme.WrapOptions{
		Theme:    cfg.themeName,
		EmbedCSS: !cfg.noCSS,
	}), nil
}

func runWatch(path string, cfg config) {
	abs, err := filepath.Abs(path)
	if err != nil {
		die(err)
	}

	var srv *serve.Server
	if cfg.serveAddr != "" {
		srv = serve.New(cfg.serveAddr, filepath.Base(abs))
		go func() {
			if err := srv.Start(); err != nil {
				fmt.Fprintf(os.Stderr, "mend: serve: %v\n", err)
			}
		}()
		fmt.Fprintf(os.Stderr, "mend: serving on http://localhost%s\n", cfg.serveAddr)
	}

	rebuild := func() {
		in, err := os.ReadFile(abs)
		if err != nil {
			fmt.Fprintf(os.Stderr, "mend: %v\n", err)
			return
		}
		out, err := assemble(in, cfg)
		if err != nil {
			fmt.Fprintf(os.Stderr, "mend: %v\n", err)
			return
		}
		if cfg.output != "" {
			if err := os.WriteFile(cfg.output, out, 0o644); err != nil {
				fmt.Fprintf(os.Stderr, "mend: %v\n", err)
				return
			}
			fmt.Fprintf(os.Stderr, "mend: wrote %s\n", cfg.output)
		} else {
			fmt.Fprintln(os.Stderr, "mend: rebuilt (use -o to persist)")
		}
		if srv != nil {
			srv.Set(out)
		}
	}

	rebuild()

	w, err := watch.New(abs)
	if err != nil {
		die(err)
	}
	defer w.Close()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	events := w.Events()
	fmt.Fprintf(os.Stderr, "mend: watching %s (ctrl-c to stop)\n", abs)
	for {
		select {
		case ev, ok := <-events:
			if !ok {
				return
			}
			if ev {
				rebuild()
			}
		case <-sigCh:
			fmt.Fprintln(os.Stderr, "mend: stopped")
			return
		}
	}
}

func die(err error) {
	fmt.Fprintln(os.Stderr, "mend:", err)
	os.Exit(1)
}

func usage() {
	fmt.Fprintf(os.Stderr, `mend %s - markdown to HTML, fast and standalone

Usage:
  mend [flags] file.md           render file to stdout
  cat file.md | mend             render stdin to stdout
  mend -o out.html file.md       write to file
  mend -w -o out.html file.md    watch and re-render on change
  mend -w -serve :8080 file.md   watch and serve with live reload
  mend -theme monokai file.md    pick a theme
  mend -fragment file.md         emit HTML fragment only

Flags:
`, version.String())
	flag.PrintDefaults()
}
