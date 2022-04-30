package main

import (
	"fmt"
	"log"

	"github.com/bueti/yeoman/internal/datasource"
	"github.com/bueti/yeoman/internal/db"
	transportHttp "github.com/bueti/yeoman/internal/transport/http"
)

func Run() error {
	log.Println("starting up Yeoman API Server")
	db, err := db.NewDatabase()
	if err != nil {
		fmt.Println("Failed to connect to database")
		return err
	}

	if err := db.MigrateDB(); err != nil {
		fmt.Println("failed to migrate database")
		return err
	}

	// open up the service
	dsService := datasource.NewService(db)

	httpHandler := transportHttp.NewHandler(dsService)
	if err := httpHandler.Serve(); err != nil {
		return err
	}

	return nil
}

func main() {
	fmt.Println("~= Yeoman API Server =~")
	if err := Run(); err != nil {
		fmt.Println(err)
	}
}
