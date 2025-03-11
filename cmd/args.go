package main

import (
	"flag"
	"fmt"
	"os"
)

type Args struct {
	noInit bool
}

func ParseArgs(args []string) Args {
	var a Args
	var v bool

	f := flag.NewFlagSet("vectro", flag.ExitOnError)
	f.Usage = func() {
		fmt.Printf("vectro - the rpn calculator • %s • %s\n", version, date)
		fmt.Println("There are no command line options. It's a calculator.")
	}
	f.BoolVar(&a.noInit, "q", false, "disable initialization")
	f.BoolVar(&v, "v", false, "show version")
	f.BoolVar(&v, "version", false, "show version")

	_ = f.Parse(args)
	if v {
		f.Usage()
		os.Exit(0)
	}
	return a
}
