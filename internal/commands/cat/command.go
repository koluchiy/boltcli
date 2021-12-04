package cat

import (
	"errors"
	"github.com/koluchiy/boltcli/internal/commands"
	bolt "go.etcd.io/bbolt"
	"os"
	"strings"
)

type Command struct {
	commands.Command
	Args struct {
		Paths []string `positional-arg-name:"paths" required:"1" description:"Paths for cat"`
	} `positional-args:"true"`
}

func NewCommand(cmd commands.Command) *Command {
	return &Command{
		Command: cmd,
	}
}

func (c *Command) Execute(args []string) error {
	err := c.BuildVals()
	if err != nil {
		return err
	}
	db, err := bolt.Open(c.File, 0600, nil)
	if err != nil {
		return err
	}
	defer db.Close()

	err = db.View(func(tx *bolt.Tx) error {
		for _, path := range c.Args.Paths {
			err := c.catPath(tx, path)
			if err != nil {
				c.PrintError([]byte(err.Error()))
			}
		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (c *Command) catPath(tx *bolt.Tx, path string) error {
	buckets := strings.Split(path, c.Delimiter)
	if len(buckets) < 2 {
		return errors.New("bad key path")
	}
	bucket := tx.Bucket([]byte(buckets[0]))
	if bucket == nil {
		return errors.New("bucket not exists: " + buckets[0])
	}
	for i:=1; i<len(buckets)-1; i++ {
		bucket = bucket.Bucket([]byte(buckets[i]))
		if bucket == nil {
			return errors.New("bucket not exists: " + strings.Join(buckets[:i+1], c.Delimiter))
		}
	}
	lastBucket := bucket.Bucket([]byte(buckets[len(buckets)-1]))
	if lastBucket != nil {
		return errors.New(os.Args[2] + " is bucket")
	}
	v := bucket.Get([]byte(buckets[len(buckets)-1]))
	c.PrintOut(v)

	return nil
}
