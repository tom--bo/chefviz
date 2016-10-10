package main

import (
	"flag"
	"fmt"
	"io"
	"os"
)

// Exit codes are int values that represent an exit code for a particular error.
const (
	ExitCodeOK    int = 0
	ExitCodeError int = 1 + iota
)

// CLI is the command line object
type CLI struct {
	// outStream and errStream are the stdout and stderr
	// to write message from the CLI.
	outStream, errStream io.Writer
}

// Run invokes the CLI with the given arguments.
func (cli *CLI) Run(args []string) int {
	var (
		rootDir string = ""
		recipe  bool
		version bool
	)
	if len(args) < 2 {
		fmt.Println("No role or recipes are specified.")
		return 1
	}

	// Define option flag parse
	flags := flag.NewFlagSet(Name, flag.ContinueOnError)
	flags.SetOutput(cli.errStream)

	flags.BoolVar(&recipe, "recipe", false, "Specify a recipe.")
	// flags.BoolVar(&recipe, "re", false, "(Short)")
	flags.StringVar(&rootDir, "rootdir", "", "")
	flags.StringVar(&rootDir, "r", "", "(Short)")

	flags.BoolVar(&version, "version", false, "Print version information and quit.")

	// Parse commandline flag
	if err := flags.Parse(args[1:]); err != nil {
		return ExitCodeError
	}

	// Show version
	if version {
		fmt.Fprintf(cli.errStream, "%s version %s\n", Name, Version)
		return ExitCodeOK
	}

	currentDir, _ := os.Getwd()
	if rootDir == "" {
		rootDir = currentDir
	} else if rootDir[:1] != "/" {
		rootDir = currentDir + "/" + rootDir + "/"
	}

	// check rootdir has 'cookbooks' dir
	if _, err := os.Stat(rootDir + "/cookbooks"); err != nil {
		fmt.Println("cookbooks directory is not found.")
		return 1
	}

	_ = recipe
	cv := &Chefviz{outStream: cli.outStream, errStream: cli.errStream}
	cv.newChefviz()
	cv.main(rootDir, flags.Args())

	return ExitCodeOK
}
