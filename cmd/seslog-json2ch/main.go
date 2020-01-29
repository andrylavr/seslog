package main

import (
	"flag"
	"fmt"
	"pkg/seslog"
)

var (
	BuildDate            string
	GitBranch, GitCommit string
)

var (
	options seslog.Options
)

func init() {
	flag.StringVar(&options.CHDSN, "dsn", "native://127.0.0.1:9000?compress=1", "ClickHouse DSN")
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

	chWriter := seslog.NewCHWriter(options)
	chWriter.FromFilesToCH()
}
