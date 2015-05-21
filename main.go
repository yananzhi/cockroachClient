// Author: yananzhi (zac.zhiyanan@gmail.com)

package main

import (
	"fmt"
	"os"

	"github.com/cockroachdb/cli/cmd"
)

func main() {

	if len(os.Args) == 1 {
		os.Args = append(os.Args, "help")
	}
	if err := cmd.Run(os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "Failed running command %q: %v\n", os.Args[1:], err)
		os.Exit(1)
	}
}
