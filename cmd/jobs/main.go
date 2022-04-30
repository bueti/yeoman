package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/bueti/yeoman/internal/db"
)

func Run() error {
	log.Println("starting up yeoman job")

	db, err := db.NewDatabase()
	if err != nil {
		fmt.Println("Failed to connect to database")
		return err
	}

	if err := db.MigrateDB(); err != nil {
		fmt.Println("failed to migrate database")
		return err
	}

	// GET all datasources
	// for each datasource
	// fetch the stuff
	// transform it into a list of ips
	// POST it to the database

	// temp.
	var urls = []string{
		"https://bunnycdn.com/api/system/edgeserverlist",
		"https://stripe.com/files/ips/ips_webhooks.json",
	}

	client := http.Client{}

	for _, url := range urls {
		log.Println("Processing:", url)
		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Set("Accept", "application/json")
		resp, err := client.Do(req)
		if err != nil {
			log.Fatal(err)
			return err
		}

		defer resp.Body.Close()

		fmt.Println(resp)

	}
	return nil
}

func main() {
	log.Println("~= Yeoman Job =~")
	if err := Run(); err != nil {
		fmt.Println(err)
	}
}
