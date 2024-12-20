package config

import (
	"log"
	"os"
	"path/filepath"
)

const (
	stdPort    = "7540"
	varEnvPort = "TODO_PORT"

	stdPass    = "1234"
	varEnvPass = "TODO_PASSWORD"

	varEnvDBFile = "TODO_DBFILE"
	stdDbPath    = "./db"
	stdDbName    = "scheduler.db"
)

func Port() string {
	port, exists := os.LookupEnv(varEnvPort)
	if !exists || port == "" {
		port = stdPort
	}

	log.Printf(`Retrieved port %s from env variable "%s"`, port, varEnvPort)

	return ":" + port
}

func Password() string {
	password, exists := os.LookupEnv(varEnvPass)
	if !exists || password == "" {
		password = stdPass
	}

	return password
}

func DbPath() string {
	storagePath, exists := os.LookupEnv(varEnvDBFile)

	if !exists || storagePath == "" {
		storagePath = filepath.Join(stdDbPath, stdDbName)
		log.Printf(`Database storage address: %s`, storagePath)
	} else {
		log.Printf(`Database storage address %s retrieved from env variable "%s" `,
			storagePath,
			varEnvDBFile)
	}

	return storagePath
}
