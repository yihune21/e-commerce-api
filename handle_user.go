package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/yihune21/e-commerce-api/internal/database"
	jwtAuth "github.com/yihune21/e-commerce-api/jwt"
	passwordhashing "github.com/yihune21/e-commerce-api/password_hashing"
)

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
		respondWithError(w,400 ,fmt.Sprintf("Error with password hashing %v",err))
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
		return
	}

    
    fmt.Printf("Dear user %s,You've successfully created an account!\n",user.Name)
    respondWithJSON(w , 200,databaseUserToUser(user))
}

func (apiConf *apiConfig)handlerGetUserByUserId(w http.ResponseWriter ,r *http.Request , user database.User){
	respondWithJSON(w, 200 , databaseUserToUser(user))
}

func (apiConf apiConfig)Login(w http.ResponseWriter , r *http.Request){
	type parameters struct{
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
    
	user , err := apiConf.db.GetUserByEmail(r.Context(), params.Email)
	if err != nil{
		respondWithError(w ,404 , "User not found")
		return
	}

	is_matched := passwordhashing.VerifyPassword(params.Password , user.Password)
	if !is_matched{
        respondWithError(w , 400 ,"Invalid Credential")
		return
	}
    
	fmt.Printf("Dear user %s,You're logged in successfully!\n",user.Name)

	// Generate both access and refresh tokens
	access_token := jwtAuth.GenerateAccessToken(user)
	refresh_token := jwtAuth.GenerateRefreshToken(user)

	// Store refresh token in database
	_, err = apiConf.db.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		ID: uuid.New(),
		UserID: user.ID,
		Token: refresh_token,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
		CreatedAt: time.Now().UTC(),
	})
	if err != nil {
		respondWithError(w, 500, fmt.Sprintf("Failed to store refresh token: %v", err))
		return
	}

    respondWithJSON(w , 200 ,ResponseToken(access_token, refresh_token))
}

func (apiConf apiConfig) UpdateUserPassword(w http.ResponseWriter , r *http.Request , user database.User)  {
	type parameters struct{
          Currentassword string `json:"current_password"`
          NewPassword string `json:"new_password"`
	}

	decode := json.NewDecoder(r.Body)
	params := parameters{}

	err := decode.Decode(&params)

	if err !=  nil {
		respondWithError(w , 400 , fmt.Sprintf("Error with parsing json %v " ,err))
		return 
	}

	is_matched := passwordhashing.VerifyPassword(params.Currentassword , user.Password)
	if !is_matched {
		respondWithError(w , 401 , "Incorrect current password!")
		return
	}
    
	hash_password,err := passwordhashing.HashPassword(params.NewPassword)
	if err != nil {
	   respondWithError(w,400 , "Couldn't able to hash the new password.")
	   return
	}
    user,err = apiConf.db.UpdateUserPasword(r.Context(),database.UpdateUserPaswordParams{
		Password: hash_password,
		ID: user.ID,
	})
    
	fmt.Printf("Dear user %s,password updated successfully!\n",user.Name)
	respondWithJSON(w,200,databaseUserToUser(user))

}


func (apiConf apiConfig) RefreshToken(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		RefreshToken string `json:"refresh_token"`
	}

	decode := json.NewDecoder(r.Body)
	params := parameters{}

	err := decode.Decode(&params)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}

	// Verify the refresh token and extract user ID
	userID, err := jwtAuth.VerifyRefreshToken(params.RefreshToken)
	if err != nil {
		respondWithError(w, 401, fmt.Sprintf("Invalid refresh token: %v", err))
		return
	}

	// Check if refresh token exists in database and is not revoked
	dbToken, err := apiConf.db.GetRefreshTokenByToken(r.Context(), params.RefreshToken)
	if err != nil {
		respondWithError(w, 401, "Refresh token not found or expired")
		return
	}

	// Verify that the token belongs to the correct user
	if dbToken.UserID != userID {
		respondWithError(w, 401, "Token user mismatch")
		return
	}

	// Get user from database
	user, err := apiConf.db.GetUserById(r.Context(), userID)
	if err != nil {
		respondWithError(w, 404, "User not found")
		return
	}

	// Generate new access token
	newAccessToken := jwtAuth.GenerateAccessToken(user)

	// Optionally generate new refresh token (refresh token rotation)
	newRefreshToken := jwtAuth.GenerateRefreshToken(user)

	// Revoke old refresh token
	err = apiConf.db.RevokeRefreshToken(r.Context(), params.RefreshToken)
	if err != nil {
		respondWithError(w, 500, "Failed to revoke old refresh token")
		return
	}

	// Store new refresh token
	_, err = apiConf.db.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		ID:        uuid.New(),
		UserID:    user.ID,
		Token:     newRefreshToken,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
		CreatedAt: time.Now().UTC(),
	})
	if err != nil {
		respondWithError(w, 500, "Failed to store new refresh token")
		return
	}

	fmt.Printf("Tokens refreshed for user %s\n", user.Name)
	respondWithJSON(w, 200, ResponseToken(newAccessToken, newRefreshToken))
}

func (apiConf apiConfig)ForgotPassword(w http.ResponseWriter , r *http.Request , user database.User)  {
	type parameters struct{
		// Email string `json:"email"`
		NewPassword string `json:"new_password"`
	}
	
	decode := json.NewDecoder(r.Body)
	params := parameters{}
	
	err := decode.Decode(&params)
	
	if err != nil{
		respondWithError(w , 400 , fmt.Sprintf("Error with decoding parameters %v",err))
		return
	}

	// user, err = apiConf.db.GetUserByEmail(r.Context(),params.Email)
    
	// if err != nil {
	// 	respondWithError(w , 404 , fmt.Sprintf("User not found %v",err))
    //     return
	// }

	otp := generateSecureOTP(6)
	fmt.Printf("You otp %v",otp)
    db_otp,err := apiConf.db.CreateOtp(r.Context(), database.CreateOtpParams{
		ID: uuid.New(),
		Otp: otp,
		UserID: user.ID,
		ExpAt: time.Now().Add(5 * time.Minute),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	},
   )
   is_matched := VerifyOTP(db_otp.Otp,otp)
   
   if !is_matched{
	  respondWithError(w,401,"Incorrect OTP!")
	  return
   }
   hashed_password,err := passwordhashing.HashPassword(params.NewPassword)
   if err != nil {
	 respondWithError(w, 400 ,fmt.Sprintf("Error with password hashing %v",err))
	 return
   }
   user , err = apiConf.db.UpdateUserPasword(r.Context(),database.UpdateUserPaswordParams{
		Password: hashed_password,
		ID: user.ID,
   })
    
   if err != nil {
	 respondWithError(w, 400 ,fmt.Sprintf("Error with updating user password  %v",err))
	 return
   }


   respondWithJSON(w,200 ,databaseUserToUser(user) )


}

func (apiConf apiConfig)LogOut(w http.ResponseWriter , r *http.Request , user database.User)  {
	 //TODO:
}