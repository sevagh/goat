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
	tagPrefixPtr := flag.String("tagPrefix", "GOAT-IN", "Prefix for GOAT related tags")

	tagPrefixEnv := os.Getenv("GOAT_TAG_PREFIX")
	logLevelEnv := os.Getenv("GOAT_LOG_LEVEL")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: goat [OPTIONS] ebs|eni\n\nOPTIONS\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	if *versionPtr {
		fmt.Println("goat ", VERSION)
		os.Exit(0)
	}

	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(1)
	}

	command := flag.Args()[0]

	logLevel := ""
	if logLevelEnv != "" {
		logLevel = logLevelEnv // env var takes precedence
	} else {
		logLevel = *logLevelPtr
	}

	tagPrefix := ""
	if tagPrefixEnv != "" {
		tagPrefix = tagPrefixEnv
	} else {
		tagPrefix = *tagPrefixPtr
	}

	log.SetOutput(os.Stderr)
	if level, err := log.ParseLevel(logLevel); err != nil {
		log.Fatalf("%v", err)
	} else {
		log.SetLevel(level)
	}

	log.SetFormatter(&log.TextFormatter{})

	log.Printf("Running goat for %s", command)
	if command == "ebs" {
		GoatEbs(*debugPtr, tagPrefix)
	} else if command == "eni" {
		GoatEni(*debugPtr, tagPrefix)
	} else {
		log.Fatalf("Unrecognized command: %s", command)
	}
}
