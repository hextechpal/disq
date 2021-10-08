package server

import (
	"errors"
	"github.com/ppal31/disq/internal/storage"
	"gopkg.in/alecthomas/kingpin.v2"
	"os"
)

type Command struct {
	backend string
	file    string
}

func (c *Command) run(*kingpin.ParseContext) error {
	var st storage.Storage
	var err error
	switch c.backend {
	case "InMemory":
		st = initInMemory()
	case "OnDisk":
		if c.file == "" {
			return errors.New("provide filepath for ondisk backend")
		}
		st, err = initOnDisk(c.file)
		if err != nil {
			return err
		}
	default:
		return errors.New("provide appropriate backend")
	}
	s := &StorageServer{s: st}
	return s.Start()
}

func initOnDisk(filePath string) (storage.Storage, error) {
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		return nil, err
	}
	return storage.NewOnDisk(file), nil
}

func initInMemory() storage.Storage {
	return storage.NewInMemory()
}

func Register(app *kingpin.Application) {
	c := new(Command)
	cmd := app.Command("server", "Starts a disq server").Action(c.run)
	cmd.Flag("backend", "Backend to use for disq server").
		Required().
		StringVar(&c.backend)
	cmd.Flag("file", "File for Ondisk").
		StringVar(&c.file)
}
