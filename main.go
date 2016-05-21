package main

import (
	"flag"
	"os"
	"path"

	"github.com/auvii/pidia/diaweb"
)

var (
	// host = flag.String("h", "localhost", "host to expose server on")
	port = flag.Uint("p", uint(80), "web service port")
	dir  = flag.String("d", path.Join(os.Getenv("PWD"), "tmp"), "directory for local mirror")
	conf = flag.String("c", path.Join(os.Getenv("HOME"), ".pidiarc"), "configuration file")
)

func main() {
	flag.Parse()
	diaweb.NewServer("localhost", *port, *dir, *conf, true).Start()
}
