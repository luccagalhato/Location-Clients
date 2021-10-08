package api

import (
	maps "GoogleMAPS/googlemaps"
	"GoogleMAPS/models"
	"GoogleMAPS/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// ClientNew ...
func ClientNew(s *sql.SQLStr) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) { //w http.ResponseWriter, r *http.Request
		clientNew := models.Street{}
		if err := json.NewDecoder(r.Body).Decode(&clientNew); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(models.Response{
				Status: "Bad Request",
				Error:  "",
				Data:   err.Error(),
			})
			return
		}
		lat, long, err := maps.RequestMapsNewclient(clientNew)
		if err != nil {
			log.Println(err)
			return
		}
		data, err := s.CompareRegion(lat, long)
		if err != nil {
			fmt.Println(err)
			return
		}

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		w.Header().Set("access-control-expose-headers", "*")
		w.Header().Set("Content-Type", "application/octet-stream")
		json.NewEncoder(w).Encode(models.Response{
			Status: "OK",
			Error:  "",
			Data:   data,
		})
	}
}
