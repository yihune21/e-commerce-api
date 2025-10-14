package main

import "net/http"

func(apiConf apiConfig) CreateProduct(w http.ResponseWriter , r *http.Request)  {
	type parameters struct{
		Name string `json:"name"`
	}
	
}