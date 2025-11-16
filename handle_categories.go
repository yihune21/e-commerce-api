package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/yihune21/e-commerce-api/internal/database"
)

func (apiCfg apiConfig) NewCategory(w http.ResponseWriter , r *http.Request , admin database.User)  {
	type parameters struct{
		Name string `json:"name"`
		Description string `json:"description"`
	 }
	
	decode := json.NewDecoder(r.Body)
	params := parameters{}

	err := decode.Decode(&params)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}

	description := sql.NullString{}
	if params.Description != ""{
		description.String = params.Description
		description.Valid = true
	}

    category , err := apiCfg.db.CreateCategoty(r.Context() , database.CreateCategotyParams{
		Name: params.Name,
		Description: description,
	})

	if err != nil {
		respondWithError(w , 201 , fmt.Sprintf("Couldn't create category %v",err))
		return
	}


    respondWithJSON(w , 200 , DatabaseCategoryToCategory(category))


}