// cli project main.go
package main

import (
	"flag"
	"log"
	"os"

	"github.com/cockroachdb/cli/cmd"

	"code.google.com/p/go-commander"
)

func main() {

	c := commander.Commander{
		Name: "cli",
		Commands: []*commander.Command{
			cmd.CmdTxn,

			{
				UsageLine: "listparams",
				Short:     "list all available parameters and their default values",
				Long: `
List all available parameters and their default values.
Note that parameter parsing stops after the first non-
option after the command name. Hence, the options need
to precede any additional arguments,

  cli <command> [options] [arguments].`,
				Run: func(cmd *commander.Command, args []string) {
					flag.CommandLine.PrintDefaults()
				},
			},
		},
	}

	if len(os.Args) == 1 {
		os.Args = append(os.Args, "help")
	}
	if err := c.Run(os.Args[1:]); err != nil {
		log.Fatalf("Failed running command %q: %v\n", os.Args[1:], err)
	}
}
