package main

import (
	"github.com/pubgo/g/errors"
	"github.com/pubgo/mycli/cmds"
	"os"
)

func main() {
	defer errors.Debug()
	errors.Panic(cmds.Execute("MY", os.ExpandEnv("$PWD")))
}
