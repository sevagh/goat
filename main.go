package main // import "github.com/sevagh/goat"

import (
	"flag"
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
)

//Goat version substituted by the Makefile
var VERSION string

func main() {
	logLevelPtr := flag.String("logLevel", "info", "Log level")
	versionPtr := flag.Bool("version", false, "Display version and exit")
	debugPtr := flag.Bool("debug", false, "Interactive debug prompts")

	flag.Parse()

	if *versionPtr {
		fmt.Printf("goat %s", VERSION)
		os.Exit(0)
	}

	if flag.NArg() != 1 {
		log.Fatalf("Usage: goat [OPTIONS] ebs|eni")
	}
	command := flag.Args()[0]

	log.SetOutput(os.Stderr)
	if level, err := log.ParseLevel(*logLevelPtr); err != nil {
		log.Fatalf("%v", err)
	} else {
		log.SetLevel(level)
	}

	log.SetFormatter(&log.TextFormatter{})

	log.Printf("Running goat for %s", command)
	if command == "ebs" {
		GoatEbs(*debugPtr)
	} else if command == "eni" {
		GoatEni(*debugPtr)
	} else {
		log.Fatalf("Unrecognized command: %s", command)
	}
}
