package cli

import (
	"fmt"
	"log"
	"os"
	"xona/src/editor"

	"github.com/urfave/cli/v2"
)

const VERSION string = "0.0.1"

var helpTemplate string = fmt.Sprintf(`
🌌 A minimalist and fast terminal text editor.

Usage:
  xona [command] [command options] [arguments...]

Commands:
  asas
`)

func InitCli() {
	app := &cli.App{
		Name:                  "xona",
		Description:           "🌌 A minimalist and fast terminal text editor.",
		CustomAppHelpTemplate: helpTemplate,
		Authors:               []*cli.Author{{Name: "Nehuén / Neth"}},
		Version:               VERSION,
		Action: func(*cli.Context) error {
			editor.Editor()
			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
