package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"net/http"

	_ "github.com/jackc/pgx/v4/stdlib"

	ipfsApi "github.com/ipfs/go-ipfs-api"

	"github.com/zhoreeq/ipfd/internal/app/ipfd"
	"github.com/zhoreeq/ipfd/internal/app/ipfs"
	"github.com/zhoreeq/ipfd/internal/app/store/sqlstore"
)

var configPath string
var printVersion bool

func init() {
	flag.StringVar(&configPath, "config", ".env", "config file path")
	flag.BoolVar(&printVersion, "version", false, "print build version")
}

func newDB(dbURL string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dbURL)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func newIpfsShell(ipfsAPI string) (*ipfsApi.Shell, error) {
	shell := ipfsApi.NewShell(ipfsAPI)
	if _, _, err := shell.Version(); err != nil {
		return nil, err
	}

	return shell, nil
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

	db, err := newDB(conf.DatabaseURL)
	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()
	dbStore := sqlstore.New(db)

	var ipfsShells []ipfs.Shell
	for _, uri := range conf.IpfsAPI {
		shell, err := newIpfsShell(uri)
		if err != nil {
			logger.Println(uri, err)
			continue
		}
		ipfsShells = append(ipfsShells, shell)
	}
	if len(ipfsShells) == 0 {
		logger.Println("Failed to connect to IPFS API")
		return
	}

	srv := ipfd.New(conf, logger, dbStore, ipfsShells)

	if err = http.ListenAndServe(conf.BindAddress, srv); err != nil  {
		logger.Println(err)
		return
	}
}
