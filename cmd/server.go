package main

import (
	"flag"
	"log"
	"os"
	"time"

	"github.com/belljustin/stamp/internal/configs"
	"github.com/belljustin/stamp/internal/db"
	"github.com/belljustin/stamp/internal/eth"
	"github.com/belljustin/stamp/internal/server"
)

var config *configs.Config

func init() {
	fname := flag.String("config", "", "config filepath")
	flag.Parse()

	if *fname == "" {
		log.Fatalf("must specify config file. -config /path/to/config.json")
	}

	cfile, err := os.Open(*fname)
	if err != nil {
		log.Fatalf("Failed to open config file: %v", err)
	}
	config, err = configs.ParseConfig(cfile)
	if err != nil {
		log.Fatalf("Failed to parse config file, %s: %v", *fname, err)
	}
}

func main() {
	conn := db.InitDB(config.Database)
	stampRequestDAO := db.NewStampRequestDAO(conn)
	stampDAO := db.NewStampDAO(conn)

	stamper := eth.InitStamper(config.Contract)
	scheduledStamper := eth.NewScheduledStamper(stamper, stampRequestDAO, stampDAO)
	s := scheduledStamper.Start(30 * time.Second)
	defer func() {
		s <- true
	}()

	server := server.NewServer(config, stampRequestDAO)
	server.Start()
}
