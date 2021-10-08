package cli

import (
	"github.com/ppal31/disq/cli/server"
	"gopkg.in/alecthomas/kingpin.v2"
	"os"
)

var application = "Disq"
var description = "This is a clone of of chukcha a distributed queuing system"

func Command() {
	app := kingpin.New(application, description)
	server.Register(app)
	kingpin.MustParse(app.Parse(os.Args[1:]))
}
