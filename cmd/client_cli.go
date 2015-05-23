package cmd

import (
	"flag"

	"log"
	"os"

	"github.com/spf13/cobra"
)

func init() {
	// If log directory has not been set, set -alsologtostderr to true.
	var hasLogDir, hasAlsoLogStderr bool
	for _, arg := range os.Args[1:] {
		switch arg {
		case "-log_dir", "--log_dir":
			hasLogDir = true
		case "-alsologtostderr", "--alsologtostderr":
			hasAlsoLogStderr = true
		}
	}
	if !hasLogDir && !hasAlsoLogStderr {
		if err := flag.CommandLine.Set("alsologtostderr", "true"); err != nil {
			log.Fatal(err)
		}
	}
}

var cockroachCmd = &cobra.Command{
	Use: "cli",
}

func init() {
	cockroachCmd.AddCommand(
		CmdTxn,
		CmdTpcc,
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
