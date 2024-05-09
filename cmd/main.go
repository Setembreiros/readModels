package main

import (
	"log"
	"os"
	"strings"
)

func main() {
	env := strings.TrimSpace(os.Getenv("ENVIRONMENT"))
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	infoLog.Printf("Starting ReadModels service in [%s] enviroment...\n", env)

	provider := NewProvider(infoLog, errorLog, env)
	database, err := provider.ProvideDb()
	if err != nil {
		os.Exit(1)
	}

	infoLog.Println("Readmodels service started")

	infoLog.Println("Applying migrations...")
	err = database.ApplyMigrations()
	if err != nil {
		errorLog.Println("Migrations failed")
		os.Exit(1)
	}
	infoLog.Println("Migrations finished")
}
