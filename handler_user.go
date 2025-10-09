package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/yihune21/e-commerce-api/internal/database"
	passwordhashing "github.com/yihune21/e-commerce-api/password_hashing"
)


type user struct{
	Id uuid.UUID `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Password string `json:"password"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`

}

func (apiConf apiConfig) New(w http.ResponseWriter , r *http.Request){
	type  parameters struct{
        Name string `json:"name"`
		Email string `json:"email"`
		Password string `json:"password"`
	}

	decode := json.NewDecoder(r.Body)
	params := parameters{}

	err := decode.Decode(&params)

	if err !=  nil {
		respondWithError(w , 400 , fmt.Sprintf("Error with parsing json %v " ,err))
		return 
	}
    
    hashed_password , err := passwordhashing.HashPassword(params.Password)
	if err != nil {
		fmt.Printf("Error with password hashing %v",err)
		return
	}


	user , err := apiConf.db.CreateUser(r.Context() , database.CreateUserParams{
		ID:uuid.New(),
		Name: params.Name,
		Email: params.Email,
		Password:hashed_password,
        CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	})
	if err != nil {
		respondWithError(w , 201 , fmt.Sprintf("Couldn't able to create user %v",err))
	}

   respondWithJSON(w , 200,user)
}