package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/yihune21/e-commerce-api/internal/database"
)

func (apiCfg apiConfig) NewCategory(w http.ResponseWriter , r *http.Request , admin database.User)  {
	type parameters struct{
		Name string `json:"name"`
		Description string `json:"description"`
		// ParentId uuid.UUID  `json:"parent_id"`
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
		// ParentID: params.ParentId,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	})

	if err != nil {
		respondWithError(w , 201 , fmt.Sprintf("Couldn't create category %v",err))
		return
	}
    
    
    respondWithJSON(w , 200 , DatabaseCategoryToCategory(category))
}

func (apiCfg apiConfig)UpdateCategoryName(w http.ResponseWriter , r *http.Request , admin database.User)  {
	type parameters struct{
		Name string `json:"name"`
	}
	
	decode := json.NewDecoder(r.Body)
	params := parameters{}

	err := decode.Decode(&params)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}

	idStr := chi.URLParam(r,"id")

	if idStr == ""{
		respondWithError(w  , 400 , "Missing category id")
		return
	}

    id , err := uuid.Parse(idStr)

	if err!=  nil{
		respondWithError(w , 400 , fmt.Sprintf("Error with parsing categpry id %v" ,err))
		return
	}

	category ,err  := apiCfg.db.GetCategoryById(r.Context() , id)
    if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't found the category: %v", err))
		return
	}
	dbcat , err := apiCfg.db.UpdateCategoryName(r.Context(), database.UpdateCategoryNameParams{
		Name: params.Name,
		ID: category.ID,
	})
	
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't update the category: %v", err))
		return
	}

	respondWithJSON(w,200 , DatabaseCategoryToCategory(dbcat))
	
}

func (apiCfg apiConfig)GetCategory(w http.ResponseWriter , r *http.Request , admin database.User)  {
	  categories,err := apiCfg.db.GetCategory(r.Context())
	  if err != nil {
		respondWithError(w , 400 , fmt.Sprintf("Couldn't get categories %v" , err))
	    return
	  }

	  respondWithJSON(w , 200 ,DatabaseCategorysToCategorys(categories))
}

func (apiCfg apiConfig)DeleteCategory(w http.ResponseWriter , r *http.Request , admin database.User)  {
	idStr := chi.URLParam(r,"id")
	if idStr == "" {
		respondWithError(w ,400 , "Missing category id")
		return
	}

	id,err := uuid.Parse(idStr)
    if err != nil {
		respondWithError(w ,400 , fmt.Sprintf("Error with parsing category id %v" ,err))
		return
	}
	
	
	err = apiCfg.db.DeleteCategoryById(r.Context() , id)
	if err != nil {
		respondWithError(w , 400 , fmt.Sprintf("Couldn't delete category %v" , err))
		return
	}

	respondWithJSON(w , 200 ,struct{}{})
}