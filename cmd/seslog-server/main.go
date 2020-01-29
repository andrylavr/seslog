package main

import (
	"flag"
	"fmt"
	"log"
	"github.com/andrylavr/seslog/pkg/seslog"
)

var (
	BuildDate            string
	GitBranch, GitCommit string
)

var (
	options seslog.Options
	configFilename string
)

func init() {
	flag.StringVar(&configFilename, "config", seslog.DEFAULT_CONFIG_LOCATION, "Config file path")
}

func main() {
	flag.Usage = func() {
		fmt.Println("NAME:")
		fmt.Println("  Seslog - syslog server for storing nginx access_log in ClickHouse")
		fmt.Println("VERSION:")
		fmt.Printf("  rev[%s] %s (%s UTC).\n", GitCommit, GitBranch, BuildDate)
		fmt.Println("USAGE:")
		flag.PrintDefaults()
	}
	flag.Parse()

    options, err := seslog.LoadConfig(configFilename)
	if err != nil {
        log.Fatal(err)
        return
    }

	accessLogServer, err := seslog.NewAccessLogServer(options)
	if err != nil {
		log.Fatal(err)
	}

    fmt.Printf("Listening address is %s", options.Address)
	if err := accessLogServer.RunServer(); err != nil {
		log.Fatal(err)
	}
}
