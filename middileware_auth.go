package main

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/yihune21/e-commerce-api/internal/auth"
	"github.com/yihune21/e-commerce-api/internal/database"
)


type authHandler func(http.ResponseWriter , *http.Request , database.User) 

func (apiConf *apiConfig) middlewareAuth(handler authHandler) http.HandlerFunc{
       
	return func(w http.ResponseWriter , r *http.Request){
          acces_token ,  err := auth.GetToken(r.Header)
		  if err != nil{
			respondWithError(w , 401 , fmt.Sprintf("Auth Error %s" , err))
			return
		  }
		  user , err := apiConf.db.GetUserById(r.Context() , uuid.MustParse(acces_token))
          if err != nil{
			 respondWithError(w , 404 , fmt.Sprintf("Couldn't found user %s" , err))
			 return
		  }
		  handler(w, r , user)
	}

}