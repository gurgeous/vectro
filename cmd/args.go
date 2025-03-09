package main

import (
	"fmt"
	"os"

	"github.com/gurgeous/vectro/internal"
)

type Args struct {
	noInit bool
}

func ParseArgs(args []string) Args {
	var a Args
	for _, arg := range args {
		switch arg {
		case "-q":
			a.noInit = true
		case "-h", "-v", "--help", "--version":
			help(nil)
		default:
			help(fmt.Errorf("unknown flag %s", arg))
		}
	}
	return a
}

func help(err error) {
	fmt.Printf("vectro - the rpn calculator • %s • %s\n", version, date)
	fmt.Println("There are no command line options. It's a calculator.")

	if err != nil {
		fmt.Println()
		fmt.Println(internal.LG.Foreground(internal.Red500).Render(err.Error()))
		os.Exit(1)
	}

	os.Exit(0)
}
