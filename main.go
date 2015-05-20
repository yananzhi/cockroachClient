//// cli project main.go
//package main

//import (
//	"flag"
//	"log"
//	"os"

//	"github.com/cockroachdb/cli/cmd"

//	"code.google.com/p/go-commander"
//)

//func main() {

//	c := commander.Commander{
//		Name: "cli",
//		Commands: []*commander.Command{
//			cmd.CmdTxn,

//			{
//				UsageLine: "listparams",
//				Short:     "list all available parameters and their default values",
//				Long: `
//List all available parameters and their default values.
//Note that parameter parsing stops after the first non-
//option after the command name. Hence, the options need
//to precede any additional arguments,

//  cli <command> [options] [arguments].`,
//				Run: func(cmd *commander.Command, args []string) {
//					flag.CommandLine.PrintDefaults()
//				},
//			},
//		},
//	}

//	if len(os.Args) == 1 {
//		os.Args = append(os.Args, "help")
//	}
//	if err := c.Run(os.Args[1:]); err != nil {
//		log.Fatalf("Failed running command %q: %v\n", os.Args[1:], err)
//	}
//}

// Copyright 2014 The Cockroach Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or
// implied. See the License for the specific language governing
// permissions and limitations under the License. See the AUTHORS file
// for names of contributors.
//
// Author: Spencer Kimball (spencer.kimball@gmail.com)

package main

import (
	"fmt"
	"os"

	"github.com/cockroachdb/cli/cmd"

	"github.com/spf13/cobra"
)

func init() {
	//	// If log directory has not been set, set -alsologtostderr to true.
	//	var hasLogDir, hasAlsoLogStderr bool
	//	for _, arg := range os.Args[1:] {
	//		switch arg {
	//		case "-log_dir", "--log_dir":
	//			hasLogDir = true
	//		case "-alsologtostderr", "--alsologtostderr":
	//			hasAlsoLogStderr = true
	//		}
	//	}
	//	if !hasLogDir && !hasAlsoLogStderr {
	//		if err := flag.CommandLine.Set("alsologtostderr", "true"); err != nil {
	//			log.Fatal(err)
	//		}
	//	}
}

func main() {

	if len(os.Args) == 1 {
		os.Args = append(os.Args, "help")
	}
	if err := Run(os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "Failed running command %q: %v\n", os.Args[1:], err)
		os.Exit(1)
	}
}

var cockroachCmd = &cobra.Command{
	Use: "cockroach",
}

func init() {
	cockroachCmd.AddCommand(
		cmd.CmdTxn,
	)

	// The default cobra usage and help templates have some
	// ugliness. For example, the "Additional help topics:" section is
	// shown unnecessarily and it doesn't place a newline before the
	// "Flags:" section if there are no subcommands. We should really
	// get these tweaks merged upstream.
	cockroachCmd.SetUsageTemplate(`{{ $cmd := . }}Usage: {{if .Runnable}}
  {{.UseLine}}{{if .HasFlags}} [flags]{{end}}{{end}}{{if .HasSubCommands}}
  {{ .CommandPath}} [command]{{end}}{{if gt .Aliases 0}}

Aliases:
  {{.NameAndAliases}}
{{end}}{{if .HasExample}}

Examples:
{{ .Example }}{{end}}{{ if .HasRunnableSubCommands}}

Available Commands: {{range .Commands}}{{if and (.Runnable) (not .Deprecated)}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{end}}
{{ if .HasLocalFlags}}
Flags:
{{.LocalFlags.FlagUsages}}{{end}}{{ if .HasInheritedFlags}}
Global Flags:
{{.InheritedFlags.FlagUsages}}{{end}}{{if .HasHelpSubCommands}}
Additional help topics:
{{if .HasHelpSubCommands}}{{range .Commands}}{{if and (not .Runnable) (not .Deprecated)}} {{rpad .CommandPath .CommandPathPadding}} {{.Short}}{{end}}{{end}}{{end}}
{{end}}{{ if .HasSubCommands }}
Use "{{.Root.Name}} help [command]" for more information about a command.
{{end}}`)
	cockroachCmd.SetHelpTemplate(`{{with or .Long .Short }}{{. | trim}}

{{end}}{{if or .Runnable .HasSubCommands}}{{.UsageString}}
{{end}}`)
}

// Run ...
func Run(args []string) error {
	cockroachCmd.SetArgs(args)
	return cockroachCmd.Execute()
}

// runInit initializes the engine based on the first
// store. The bootstrap engine may not be an in-memory type.
func runInit(cmd *cobra.Command, args []string) {
	// First initialize the Context as it is used in other places.
	fmt.Println("runInit execute")

}
