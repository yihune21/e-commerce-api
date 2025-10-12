package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func respondWithError(w http.ResponseWriter , code int , msg string)  {
	
	if code > 499 {
       log.Println("Responding with 5xx error" , msg)
	} 
	type errorResponse struct{
		Error string `json:"error"`
	}

	respondWithJSON(w , code , errorResponse{
		Error: msg,
	})
}

func respondWithJSON(w http.ResponseWriter ,code int ,  payload interface{} ){
	dat , err := json.Marshal(payload)
	if err != nil{
		log.Printf("Failed to marshal %v" , payload)
		w.WriteHeader(code)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(dat)
}