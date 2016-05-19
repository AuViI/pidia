package main

import (
	"flag"
	"fmt"
	"github.com/auvii/pidia/diaweb"
	"os"
)

var (
	// host = flag.String("h", "localhost", "host to expose server on")
	port = flag.Uint("p", uint(80), "web service port")
	dir  = flag.String("d", ".", "directory for local mirror")
	conf = flag.String("c", fmt.Sprintf("%s/.pidiarc", os.Getenv("HOME")), "configuration file")
)

func main() {
	flag.Parse()
	diaweb.NewServer("localhost", *port, *dir, *conf, false).Start()
}
