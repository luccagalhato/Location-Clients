package main

import (
	"GoogleMAPS/api"
	c "GoogleMAPS/config"
	"GoogleMAPS/sql"
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

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
	log.Printf("starting server '%s' at port: %s", c.Config.API.Host, c.Config.API.Port)

	r := mux.NewRouter()
	r.PathPrefix("/html/").Handler(http.StripPrefix("/html/", http.FileServer(http.Dir("html"))))
	r.Path("/request").Methods(http.MethodPost).HandlerFunc(api.ClientNew)

	r.Path("/").Methods(http.MethodGet).HandlerFunc(redirect)
	server := &http.Server{
		Handler:      r,
		Addr:         ":" + c.Config.API.Port,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Fatal(server.ListenAndServe())
	// connection.disconnect()ClientNew
}
func redirect(w http.ResponseWriter, req *http.Request) {
	// remove/add not default ports from req.Host
	target := "http://" + req.Host + "/html/register.html"
	// log.Println(target)
	if len(req.URL.RawQuery) > 0 {
		target += "?" + req.URL.RawQuery
	}
	// log.Printf("redirect to: %s", target)
	http.Redirect(w, req, target,
		// see comments below and consider the codes 308, 302, or 301
		http.StatusTemporaryRedirect)
}
