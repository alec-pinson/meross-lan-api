package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// writes api response to webpage and app log
func writeResponse(w http.ResponseWriter, response any, Error bool) {
	responseString, err := json.Marshal(response)
	if err != nil {
		log.Printf("Error writing API response: %v", err)
		fmt.Fprintf(w, "Error writing API response: %v", err)
		return
	}
	log.Println(string(responseString))
	if Error {
		w.WriteHeader(http.StatusBadRequest)
	}
	json.NewEncoder(w).Encode(response)
}
