package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/zhoreeq/ipfd/internal/app/ipfd"
)

var configPath string
var printVersion bool

func init() {
	flag.StringVar(&configPath, "config", ".env", "config file path")
	flag.BoolVar(&printVersion, "version", false, "print build version")
}

func main() {
	flag.Parse()
	if printVersion {
		fmt.Println(ipfd.Version)
		return
	}
	logger := log.New(os.Stdout, "", log.Flags())

	conf, err := ipfd.LoadConfig(configPath)
	if err != nil {
		logger.Fatal(err)
	}

	err = ipfd.Start(conf, logger)
	if err != nil {
		logger.Fatal(err)
	}
}
