package main

import (
	"net/http"
)

func handlerHealthy(w http.ResponseWriter , r *http.Request)  {
	respondWithJSON(w ,200 , struct{}{})
}