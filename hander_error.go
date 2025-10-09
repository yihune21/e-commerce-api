package main

import "net/http"

func handelError(w http.ResponseWriter ,  r * http.Request)  {
	respondWithError(w , 400 , "Something went wrong!")
}