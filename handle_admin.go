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

func (apiConf apiConfig)CreateAdmin(w http.ResponseWriter ,r *http.Request)  {
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

	if len(params.Password) < 8{
		respondWithError(w , 400 , "Password should be equal or larger than 8")
		return 
	}
    
    hashed_password , err := passwordhashing.HashPassword(params.Password)
	if err != nil {
		respondWithError(w,400 ,fmt.Sprintf("Error with password hashing %v",err))
		return
	}


	admin , err := apiConf.db.CreateAdmin(r.Context() , database.CreateAdminParams{
		ID:uuid.New(),
		Name: params.Name,
		Email: params.Email,
		Password:hashed_password,
        CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	})
	if err != nil {
		respondWithError(w , 201 , fmt.Sprintf("Couldn't able to create admi  %v",err))
		return
	}

    
    fmt.Printf("Dear user %s,You've successfully created an account!\n",admin.Name)
    respondWithJSON(w , 200,databaseAdminToAdmin(admin))

}

func (apiConf apiConfig)DeleteUser(w http.ResponseWriter ,r *http.Request,admin database.Admin)  {
	 type parameters struct{
		UserID uuid.UUID `json:"user_id"`
	 }
	 decode := json.NewDecoder(r.Body)
	 params := parameters{}

	 err := decode.Decode(&params)

	 if err != nil {
		respondWithError(w , 400 , fmt.Sprintf("Error with parsing json %v", err))
		return
	 }

     apiConf.db.DeleteUserByUserId(r.Context() , params.UserID)

     respondWithJSON(w,200 ,ResponseHealth("Successfully deleted user ") )
}