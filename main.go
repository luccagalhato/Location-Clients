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
	sql.SetSQLConnLinx(connection)
	htmlFS, err := fs.Sub(content, "html")
	if err != nil {
		log.Println(err)
		return
	}

	log.Printf("starting server '%s' at port: %s", c.Config.API.Host, c.Config.API.Port)

	http.HandleFunc("/request", api.ClientNew)
	fs := http.FileServer(http.FS(htmlFS))
	http.Handle("/html/", http.StripPrefix("/html/", fs))
	log.Fatal(http.ListenAndServe(":8082", nil))
	// connection.disconnect()ClientNew
}

// func redirect(w http.ResponseWriter, req *http.Request) {
// 	// remove/add not default ports from req.Host
// 	target := "http://" + req.Host + "/html/register.html"
// 	// log.Println(target)
// 	if len(req.URL.RawQuery) > 0 {
// 		target += "?" + req.URL.RawQuery
// 	}
// 	// log.Printf("redirect to: %s", target)
// 	http.Redirect(w, req, target,
// 		// see comments below and consider the codes 308, 302, or 301
// 		http.StatusTemporaryRedirect)
// }
