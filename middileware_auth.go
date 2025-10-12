package main

import (
	"fmt"
	"net/http"

	"github.com/yihune21/e-commerce-api/internal/auth"
	"github.com/yihune21/e-commerce-api/internal/database"
	jwtAuth "github.com/yihune21/e-commerce-api/jwt"
)


type authHandler func(http.ResponseWriter , *http.Request , database.User) 

func (apiConf *apiConfig) middlewareAuth(handler authHandler) http.HandlerFunc{
       
	return func(w http.ResponseWriter , r *http.Request){
          access_token ,  err := auth.GetToken(r.Header)
		  if err != nil{
			respondWithError(w , 401 , fmt.Sprintf("Auth Error %s" , err))
			return
		  }
		  is_valid := jwtAuth.VerfiyToken(access_token)
          if !is_valid{
            respondWithError(w , 401 ,"ACCESS TOKEN EXPIRED!")
			return 
		  }

		  user_id,err := jwtAuth.ExtractUserIDFromToken(access_token)
		  if err != nil{
			respondWithError(w , 400 , fmt.Sprintf("Error with extracting user id %v",err))
			return
		  }

		  user , err := apiConf.db.GetUserById(r.Context() ,user_id)
          if err != nil{
			 respondWithError(w , 404 , fmt.Sprintf("Couldn't found user %s" , err))
			 return
		  }
		  handler(w, r , user)
	}

}