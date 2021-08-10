package main

import (
	"flag"
	"fmt"
	"github.com/andrylavr/seslog/src/seslog"
	"log"
)

var (
	BuildDate            string
	GitBranch, GitCommit string
)

func main() {
	options, err := seslog.ReadOptions()
	if err != nil {
		log.Fatal(err)
	}

	flag.Usage = func() {
		fmt.Println("NAME:")
		fmt.Println(" Seslog JSON2ClickHouse - push .json.gz files to ClickHouse")
		fmt.Println("VERSION:")
		fmt.Printf("  rev[%s] %s (%s UTC).\n", GitCommit, GitBranch, BuildDate)
		fmt.Println("USAGE:")
		flag.PrintDefaults()
	}
	flag.Parse()

	chWriter := seslog.NewCHWriter(options)
	chWriter.FromFilesToCH()
}
