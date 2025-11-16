package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/yihune21/e-commerce-api/internal/database"
)

func(apiConf apiConfig) CreateProduct(w http.ResponseWriter , r *http.Request , admin database.User)  {
	type parameters struct{
		Name string `json:"name"`
		Description string `json:"description"`
		Price string `json:"price"`
		Stock int32 `json:"stock"`
		CategoryId uuid.UUID `json:"category_id"`
		ImageUrl string `json:"image_url"`
        
	}
	
	decode := json.NewDecoder(r.Body)
	params := parameters{}

	err := decode.Decode(&params)
    
	if err != nil {
		respondWithError(w ,400 , fmt.Sprintf("Error with parsing json %v" ,err))
		return
	}
    description := sql.NullString{}
	if params.Description != ""{
		description.String = params.Description
		description.Valid = true
	}
	image_url := sql.NullString{}
	if params.ImageUrl != ""{
		description.String = params.ImageUrl
		description.Valid = true
	}
	created_at := sql.NullTime{}
	created_at.Time = time.Now().UTC()
	created_at.Valid = true
	updated_at := sql.NullTime{}
	updated_at.Time = time.Now().UTC()
	updated_at.Valid = true
	
	product , err := apiConf.db.CreateProduct(r.Context() , database.CreateProductParams{
		Name: params.Name,
		Description: description,
		Price: params.Price,
		Stock: params.Stock,
		CategoryID: params.CategoryId,
		ImageUrl: image_url,
		CreatedAt: created_at,
		UpdatedAt: updated_at,
	})
	if err != nil {
		respondWithError(w , 400, fmt.Sprintf("Couldn't able to add a product %v" ,err))
		return
	}

	respondWithJSON(w , 200 , DatabaseProductToProduct(product))

}

func(apiConf apiConfig) GetProductByName(w http.ResponseWriter , r *http.Request )  {
    type parameters struct {
		Name string `json:"name"`
	}
	decode := json.NewDecoder(r.Body)
	params := parameters{}

	err :=  decode.Decode(&params)
	if err != nil {
		respondWithError(w ,400 , fmt.Sprintf("Error with parsing json %v" ,err))
		return
	}
	product,err :=  apiConf.db.GetProductByName(r.Context() , params.Name)
	if err != nil {
		respondWithError(w ,400 , fmt.Sprintf("Couldn't find product %v" ,err))
		return
	}

	respondWithJSON(w , 200 , DatabaseProductToProduct(product))
}

func (apiConf apiConfig)UpdateProductPrice(w http.ResponseWriter , r *http.Request)  {
   type parameters struct{
	    Name string `json:"name"`
		Price string `json:"price"`
	}
	decode := json.NewDecoder(r.Body)
	params := parameters{}

	err :=  decode.Decode(&params)
	if err != nil {
		respondWithError(w ,400 , fmt.Sprintf("Error with parsing json %v" ,err))
		return
	}

	product,err := apiConf.db.UpdateProductPrice(r.Context() ,database.UpdateProductPriceParams{
		Price: params.Price,
		Name: params.Name,
	})
	if err != nil {
		respondWithError(w ,400 , fmt.Sprintf("Couldn't update product %v" ,err))
		return
	}

	respondWithJSON(w , 200 , DatabaseProductToProduct(product))

}
func (apiConf apiConfig)UpdateProductImage(w http.ResponseWriter , r *http.Request)  {
	//TODO
}

