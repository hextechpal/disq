package server

import (
	"errors"
	"github.com/ppal31/disq/internal/storage"
	"gopkg.in/alecthomas/kingpin.v2"
	"io"
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
		var cls io.Closer
		st, cls, err = initOnDisk(c.file)
		if err != nil {
			return err
		}
		defer cls.Close()
	default:
		return errors.New("provide appropriate backend")
	}
	s := &StorageServer{s: st}
	return s.Start()
}

func initOnDisk(filePath string) (storage.Storage, io.Closer, error) {
	if _, err := os.Stat(filePath); err != nil {
		if os.IsNotExist(err) {
			return createFile(filePath)
		}
		return nil, nil, err
	}
	if err := os.Remove(filePath); err != nil{
		return nil, nil, err
	}
	return createFile(filePath)
}

func createFile(filePath string) (storage.Storage, io.Closer, error) {
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		return nil, nil, err
	}
	return storage.NewOnDisk(file), file, nil
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
