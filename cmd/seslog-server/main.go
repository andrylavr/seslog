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
)

func init() {
	flag.StringVar(&options.CHDSN, "dsn", "native://127.0.0.1:9000?compress=1", "ClickHouse DSN")
	flag.StringVar(&options.Address, "addr", ":5514", "listen address")
	flag.StringVar(&options.FlushInterval, "flush-interval", "60s", "Interval between ClickHouse flushes")
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

	accessLogServer, err := seslog.NewAccessLogServer(options)
	if err != nil {
		log.Fatal(err)
	}
	if err := accessLogServer.RunServer(); err != nil {
		log.Fatal(err)
	}
}
