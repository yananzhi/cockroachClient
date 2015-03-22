package cmd

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"

	"code.google.com/p/go-commander"
)

// A CmdInit command initializes a new Cockroach cluster.
var CmdTxn = &commander.Command{
	UsageLine: "txn",
	Short:     "init a transactional cli",
	Long: `use a transactional kv client
`,
	Run:  runTxnKV,
	Flag: *flag.CommandLine,
}

type cmd struct {
	name string
	args []string
}

// cmdDict maps from command name to function implementing the command.
// Use only upper case letters for commands. More than one letter is OK.
var cmdDict = map[string]func(c *cmd) error{
	"S": startCmd,
}

func startCmd(c *cmd) error {
	//for test
	fmt.Printf("startcmd: %v", c)
	return nil
}

func getCmd(str string) (*cmd, error) {
	args := strings.Split(str, " ")

	if len(args) < 1 {
		return nil, fmt.Errorf("not enouf arguments")
	}

	c := &cmd{
		name: strings.ToUpper(args[0]),
		args: args[1:],
	}

	return c, nil
}

func runTxnKV(cmd *commander.Command, args []string) {
	fmt.Printf("txn kv client:\n")

	for {
		reader := bufio.NewReader(os.Stdin)
		strBytes, _, err := reader.ReadLine()

		if err == nil {
			//for test
			fmt.Println("test readline:", string(strBytes))

			if c, err := getCmd(string(strBytes)); err == nil {
				if fn, exist := cmdDict[c.name]; exist {
					fn(c)
				} else {
					fmt.Printf("cmd %v not exist", c)
				}

			} else {
				fmt.Printf("getcmd error:%v", err)
			}
		}

	}

}
