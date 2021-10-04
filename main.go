package main

import (
	"GoogleMAPS/api"
	c "GoogleMAPS/config"
	"GoogleMAPS/sql"
	"embed"
	"flag"
	"io/fs"
	"log"
	"net/http"
	"time"
)

//go:embed html
var content embed.FS

//const apiKey = "AIzaSyDybcJ7PHZPP2es7YN_hd0D5OWjIR3kue0"
var createConfig bool

func main() {
	flag.BoolVar(&createConfig, "c", false, "create config.yaml file")
	flag.Parse()

	if createConfig {
		c.CreateConfigFile()
		return
	}

	log.Print("loading config file")
	if err := c.LoadConfig(); err != nil {
		log.Fatal(err)
	}

	log.Print("connecting sql ...")
	connection, err := sql.MakeSQL(c.Config.SQL.Host, c.Config.SQL.Port, c.Config.SQL.User, c.Config.SQL.Password)
	if err != nil {
		log.Println(err)
		return
	}

	htmlFS, err := fs.Sub(content, "html")
	if err != nil {
		log.Println(err)
		return
	}

	log.Printf("starting server '%s' at port: %s", c.Config.API.Host, c.Config.API.Port)

	go func() {
		ticker := time.NewTicker(time.Second * time.Duration(c.Config.SQL.Interval))
		tickerPing := time.NewTicker(time.Second * 10)
		for {
			select {
			case <-ticker.C:
				connection.SearchClient()
			case <-tickerPing.C:
				connection.Ping()
			}
		}
	}()

	http.HandleFunc("/request", api.ClientNew(connection))
	fs := http.FileServer(http.FS(htmlFS))
	http.Handle("/html/", http.StripPrefix("/html/", fs))
	log.Fatal(http.ListenAndServe(":8082", nil))
	// connection.disconnect()ClientNew

}
