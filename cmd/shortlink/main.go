// Command shortlink is a small URL shortener: `shortlink serve` runs the HTTP
// API, and the other subcommands are a CLI client for it.
package main

import (
	"fmt"
	"io"
	"os"

	"github.com/BotJakeDemo/shortlink/internal/version"
)

const usage = `shortlink — a small URL shortener

Usage:
  shortlink <command> [flags]

Commands:
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
	case "help", "-h", "--help":
		fmt.Fprint(stdout, usage)
		return 0
	default:
		fmt.Fprintf(stderr, "shortlink: unknown command %q\n\n%s", args[0], usage)
		return 2
	}
}
