package main

import (
	"github.com/pubgo/mycli/cmds"
	"os"
)

func main() {
	cmds.Execute("MY", os.ExpandEnv("$PWD"))
}
