package main

import (
	"github.com/koluchiy/boltcli/internal/commands"
	"github.com/koluchiy/boltcli/internal/commands/cat"
	config2 "github.com/koluchiy/boltcli/internal/commands/config"
	"github.com/koluchiy/boltcli/internal/commands/ls"
	"os"

	"github.com/jessevdk/go-flags"
)

func main() {
	parser := flags.NewParser(nil, flags.Default)
	baseCommand := commands.Command{
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}

	_, err := parser.AddCommand("ls", "", "", ls.NewCommand())
	if err != nil {
		panic(err)
	}
	_, err = parser.AddCommand("cat", "", "", cat.NewCommand(baseCommand))
	if err != nil {
		panic(err)
	}
	_, err = parser.AddCommand("config", "", "", config2.NewCommand())
	if err != nil {
		panic(err)
	}

	_, err = parser.Parse()

	if err != nil {
		switch flagsErr := err.(type) {
		case flags.ErrorType:
			if flagsErr == flags.ErrHelp {
				os.Exit(0)
			}
			os.Exit(1)
		default:
			os.Exit(1)
		}
	}
}
