package config

import (
	"errors"
	"fmt"
	"github.com/koluchiy/boltcli/internal/config"
)

type Command struct {
	Global bool `short:"g" long:"global" description:"Operate with global user config"`
	Args struct {
		Key string `positional-arg-name:"key" description:"key"`
		Value string `positional-arg-name:"value" description:"value"`
	} `positional-args:"true"`
}

func NewCommand() *Command {
	return &Command{}
}

func (c *Command) Execute(args []string) error {
	if len(c.Args.Key) == 0 && len(c.Args.Value) == 0 {
		var cfg *config.Config
		var err error
		if c.Global {
			cfg, err = config.GetConfigGlobal()
		} else {
			cfg, err = config.GetConfig()
		}

		if err != nil {
			return err
		}
		fmt.Println(cfg)
		return nil
	}
	patch := config.Patch{}

	if c.Args.Key == "delimiter" {
		patch.Delimiter = &c.Args.Value
	} else if c.Args.Key == "file" {
		patch.File = &c.Args.Value
	} else {
		return errors.New("incorrect key for config: " + c.Args.Key)
	}

	err := config.PatchConfig(&patch, c.Global)
	if err != nil {
		return err
	}

	return nil
}
