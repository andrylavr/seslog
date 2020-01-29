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
		fmt.Println(" Seslog JSON2ClickHouse - push .json.gz files to ClickHouse")
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

	chWriter := seslog.NewCHWriter(options)
	chWriter.FromFilesToCH()
}
