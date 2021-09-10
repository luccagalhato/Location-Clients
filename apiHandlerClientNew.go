package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

var connectionLinx *SQLStr

//SetSQLConnLinx ...
func SetSQLConnLinx(c *SQLStr) {
	connectionLinx = c
}
func ClientNew(w http.ResponseWriter, r *http.Request) { //w http.ResponseWriter, r *http.Request
	clientNew := Street{}
	if err := json.NewDecoder(r.Body).Decode(&clientNew); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Response{
			Status: "Bad Request",
			Error:  "",
			Data:   err.Error(),
		})
		return
	}
	lat, long := requestMapsNewclient(clientNew)
	data, err := connectionLinx.CompareRegion(lat, long)
	if err != nil {
		fmt.Println(err)
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("access-control-expose-headers", "*")
	w.Header().Set("Content-Type", "application/octet-stream")
	json.NewEncoder(w).Encode(Response{
		Status: "OK",
		Error:  "",
		Data:   data,
	})
}
