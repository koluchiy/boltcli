package ls

import (
	"errors"
	"github.com/koluchiy/boltcli/internal/commands"
	bolt "go.etcd.io/bbolt"
	"strings"
)

type Command struct {
	commands.Command
	Recursive bool `short:"R" long:"recursive" description:"Show all keys recursive"`
	All bool `short:"a" long:"all" description:"Show all data for keys"`
	Args struct {
		Path string `positional-arg-name:"path" description:"Path for ls"`
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

	if len(c.Args.Path) == 0 {
		err := db.View(func(tx *bolt.Tx) error {
			err := tx.ForEach(func(name []byte, b *bolt.Bucket) error {
				c.PrintOut(name)
				return nil
			})
			if err != nil {
				return err
			}
			return nil
		})
		if err != nil {
			return err
		}
		return nil
	}
	buckets := strings.Split(c.Args.Path, c.Delimiter)
	err = db.View(func(tx *bolt.Tx) error {
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
		var lastBucket *bolt.Bucket
		if len(buckets) == 1 {
			lastBucket = bucket
		} else {
			lastBucket = bucket.Bucket([]byte(buckets[len(buckets)-1]))
		}

		if lastBucket != nil {
			var showBuckets [][]byte
			var showKeys [][]byte
			err := lastBucket.ForEach(func(k, v []byte) error {
				if v == nil {
					showBuckets = append(showBuckets, k)
				} else {
					showKeys = append(showKeys, k)
				}
				return nil
			})
			for _, b := range showBuckets {
				c.PrintOut(b, []byte(c.Delimiter))
			}
			for _, k := range showKeys {
				c.PrintOut(k)
			}
			if err != nil {
				return err
			}
		} else {
			v := bucket.Get([]byte(buckets[len(buckets)-1]))
			c.PrintOut([]byte(buckets[len(buckets)-1]))
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}
