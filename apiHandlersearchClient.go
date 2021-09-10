package main

import (
	"encoding/json"
	"net/http"
)

type Clifor struct {
	Code string `json:"code,omitempty"`
}

// SearchClient ...
func SearchClient(w http.ResponseWriter, r *http.Request) {
	codClifor := Clifor{}
	if err := json.NewDecoder(r.Body).Decode(&codClifor); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Response{
			Status: "Bad Request",
			Error:  "",
			Data:   err.Error(),
		})
		return
	}
	//data, err := connectionLinx.SearchClient(codClifor.Code)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("access-control-expose-headers", "*")
	w.Header().Set("Content-Type", "application/octet-stream")
	json.NewEncoder(w).Encode(Response{
		Status: "OK",
		Error:  "",
		Data:   nil,
	})
}
