// Command shortlink is a small URL shortener: `shortlink serve` runs the HTTP
// API, and the other subcommands are a CLI client for it.
package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/BotJakeDemo/shortlink/internal/server"
	"github.com/BotJakeDemo/shortlink/internal/version"
)

const usage = `shortlink — a small URL shortener

Usage:
  shortlink <command> [flags]

Commands:
  serve     Run the HTTP API (--addr, --data)
  version   Print the shortlink version
  help      Print this help
`

func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr))
}

// run dispatches a subcommand and returns the process exit code.
func run(args []string, stdout, stderr io.Writer) int {
	if len(args) == 0 {
		fmt.Fprint(stderr, usage)
		return 2
	}
	switch args[0] {
	case "version", "--version":
		fmt.Fprintf(stdout, "shortlink %s\n", version.Version)
		return 0
	case "serve":
		return serve(args[1:], stderr)
	case "help", "-h", "--help":
		fmt.Fprint(stdout, usage)
		return 0
	default:
		fmt.Fprintf(stderr, "shortlink: unknown command %q\n\n%s", args[0], usage)
		return 2
	}
}

// serve runs the HTTP API until the listener fails, and returns the exit code.
func serve(args []string, stderr io.Writer) int {
	fs := flag.NewFlagSet("serve", flag.ContinueOnError)
	fs.SetOutput(stderr)
	addr := fs.String("addr", "127.0.0.1:8787", "address to listen on")
	data := fs.String("data", "links.json", "path to the JSON link store")
	if err := fs.Parse(args); err != nil {
		return 2
	}

	handler, err := server.New(*data)
	if err != nil {
		fmt.Fprintf(stderr, "shortlink serve: %v\n", err)
		return 1
	}
	if err := http.ListenAndServe(*addr, handler); err != nil {
		fmt.Fprintf(stderr, "shortlink serve: %v\n", err)
		return 1
	}
	return 0
}
