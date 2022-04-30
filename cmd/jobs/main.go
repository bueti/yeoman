package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

func Run() error {
	log.Println("starting up yeoman job")
	os.Mkdir("output", 0755)

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

		tmp := strings.TrimPrefix(url, "https://")
		file := strings.ReplaceAll(tmp, "/", "_")
		f, err := os.Create("output/" + file)
		if err != nil {
			log.Fatal(err)
			return err
		}

		defer f.Close()

		_, err = f.ReadFrom(resp.Body)

		if err != nil {
			log.Fatal(err)
		}

	}
	return nil
}

func main() {
	log.Println("Yeoman Job")
	if err := Run(); err != nil {
		fmt.Println(err)
	}
}
