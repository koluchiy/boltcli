package commands

import (
	"github.com/koluchiy/boltcli/internal/config"
	"io"
)

type Command struct {
	Stdout io.Writer
	Stderr io.Writer
	File string `short:"f" long:"file" description:"Path for target file where all changes will be written"`
	Delimiter string `short:"d" long:"delimiter" description:"Path for target file where all changes will be written"`
}

func (c *Command) PrintOut(datas ...[]byte) {
	for _, data := range datas {
		c.Stdout.Write(data)
	}
	c.Stdout.Write([]byte{'\n'})
}

func (c *Command) PrintError(datas ...[]byte) {
	for _, data := range datas {
		c.Stderr.Write(data)
	}
	c.Stderr.Write([]byte{'\n'})
}

func (c *Command) BuildVals() error {
	if len(c.File) == 0 || len(c.Delimiter) == 0 {
		cfg, err := config.GetConfig()
		if err != nil {
			return err
		}

		if len(c.File) == 0 {
			c.File = cfg.File
		}
		if len(c.Delimiter) == 0 {
			c.Delimiter = cfg.Delimiter
		}
	}

	return nil
}
