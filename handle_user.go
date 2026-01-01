package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/yihune21/e-commerce-api/internal/auth"
	"github.com/yihune21/e-commerce-api/internal/database"
	jwtAuth "github.com/yihune21/e-commerce-api/jwt"
	passwordhashing "github.com/yihune21/e-commerce-api/password_hashing"
)

func (apiConf apiConfig) New(w http.ResponseWriter , r *http.Request){
	type  parameters struct{
        Name string `json:"name"`
		Email string `json:"email"`
		Phone string `json:"phone"`
		Password string `json:"password"`

	}

	decode := json.NewDecoder(r.Body)
	params := parameters{}
    
	err := decode.Decode(&params)

	if err !=  nil {
		respondWithError(w , 400 , fmt.Sprintf("Error with parsing json %v " ,err))
		return 
	}
	
	err = PasswordChecker(params.Password)

    if err != nil{
		respondWithError(w , 400 , err.Error())
		return 
	}
    
    hashed_password , err := passwordhashing.HashPassword(params.Password)
	if err != nil {
		respondWithError(w,400 ,fmt.Sprintf("Error with password hashing %v",err))
		return
	}
    
    is_valid_email , err := IsValidEmail(params.Email)
	if !is_valid_email{
		respondWithError(w , 400 , err.Error())
        return
	}

	phone := sql.NullString{}
	if params.Phone != "" {
		phone.String = params.Phone
		phone.Valid = true
	}
	
	user , err := apiConf.db.CreateUser(r.Context() , database.CreateUserParams{
		ID:uuid.New(),
		Name: params.Name,
		Email: params.Email,
		Phone: phone,
		Password:hashed_password,
		IsAdmin: false,
        CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	})
	if err != nil {
		respondWithError(w , 400 , fmt.Sprintf("Couldn't able to create user %v",err))
		return
	}

    
    fmt.Printf("Dear user %s,You've successfully created an account!\n",user.Name)
    respondWithJSON(w , 200,databaseUserToUser(user))
}
func (apiConf apiConfig) NewAdmin(w http.ResponseWriter , r *http.Request ,Admin database.User){
	type  parameters struct{
        Name string `json:"name"`
		Email string `json:"email"`
		Phone string `json:"phone"`
		Password string `json:"password"`
      
	}

	decode := json.NewDecoder(r.Body)
	params := parameters{}
    
	err := decode.Decode(&params)

	if err !=  nil {
		respondWithError(w , 400 , fmt.Sprintf("Error with parsing json %v " ,err))
		return 
	}
    err = PasswordChecker(params.Password)

    if err != nil{
		respondWithError(w , 400 , err.Error())
		return 
	}
    
    hashed_password , err := passwordhashing.HashPassword(params.Password)
	if err != nil {
		respondWithError(w,400 ,fmt.Sprintf("Error with password hashing %v",err))
		return
	}
    
    is_valid_email , err := IsValidEmail(params.Email)
	if !is_valid_email{
		respondWithError(w , 400 , err.Error())
        return
	}


	phone := sql.NullString{}
	if params.Phone != "" {
		phone.String = params.Phone
		phone.Valid = true
	}
	
	user , err := apiConf.db.CreateUser(r.Context() , database.CreateUserParams{
		ID:uuid.New(),
		Name: params.Name,
		Email: params.Email,
		Phone: phone,
		Password:hashed_password,
		IsAdmin: true,
        CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	})
	if err != nil {
		respondWithError(w , 400 , fmt.Sprintf("Couldn't able to create user %v",err))
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
    
	// fmt.Printf("Dear user %s,You're logged in successfully!\n",user.Name)

	access_token := jwtAuth.GenerateAccessToken(user)
	refresh_token := jwtAuth.GenerateRefreshToken(user)
    
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
    err = PasswordChecker(params.NewPassword)

    if err != nil{
		respondWithError(w , 400 , err.Error())
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
    apiConf.db.RevokeAllUserRefreshTokens(r.Context(),user.ID)
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

	
	userID, err := jwtAuth.VerifyRefreshToken(params.RefreshToken)
	if err != nil {
		respondWithError(w, 401, fmt.Sprintf("Invalid refresh token: %v", err))
		return
	}

	dbToken, err := apiConf.db.GetRefreshTokenByToken(r.Context(), params.RefreshToken)
	if err != nil {
		respondWithError(w, 401, "Refresh token not found or expired")
		return
	}

	if dbToken.UserID != userID {
		respondWithError(w, 401, "Token user mismatch")
		return
	}

	user, err := apiConf.db.GetUserById(r.Context(), userID)
	if err != nil {
		respondWithError(w, 404, "User not found")
		return
	}

	newAccessToken := jwtAuth.GenerateAccessToken(user)

	newRefreshToken := jwtAuth.GenerateRefreshToken(user)
    
	err = apiConf.db.RevokeRefreshToken(r.Context(), params.RefreshToken)
	if err != nil {
		respondWithError(w, 500, "Failed to revoke old refresh token")
		return
	}
    
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

	// fmt.Printf("Tokens refreshed for user %s\n", user.Name)
	respondWithJSON(w, 200, ResponseToken(newAccessToken, newRefreshToken))
}

func (apiConf apiConfig)RequestForgotPassword(w http.ResponseWriter , r *http.Request)  {
	type parameters struct{
		Email string `json:"email"`
	}
	
	decode := json.NewDecoder(r.Body)
	params := parameters{}
	
	err := decode.Decode(&params)
	
	if err != nil{
		respondWithError(w , 400 , fmt.Sprintf("Error with decoding parameters %v",err))
		return
	}
    
	user, err := apiConf.db.GetUserByEmail(r.Context(),params.Email)
    
	if err != nil {
		respondWithError(w , 404 , fmt.Sprintf("User not found %v",err))
        return
	}

	otp := generateSecureOTP(6)
	
    db_otp,err := apiConf.db.CreateOtp(r.Context(), database.CreateOtpParams{
		ID: uuid.New(),
		Otp: otp,
		UserID: user.ID,
		ExpAt: time.Now().Add(10 * time.Minute),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	},
   )
  

   respondWithJSON(w,200 , OtpRes("Sent" ,db_otp.Otp))


}
func (apiConf apiConfig)ForgotPassword(w http.ResponseWriter , r *http.Request)  {
	type parameters struct{
		Email string `json:"email"`
		Otp   string `json:"otp"`
		NewPassword string `json:"new_password"`
	}
	
	decode := json.NewDecoder(r.Body)
	params := parameters{}
	
	err := decode.Decode(&params)
	
	if err != nil{
		respondWithError(w , 400 , fmt.Sprintf("Error with decoding parameters %v",err))
		return
	}
    
	user, err := apiConf.db.GetUserByEmail(r.Context(),params.Email)
    
	if err != nil {
		respondWithError(w , 404 , fmt.Sprintf("User not found %v",err))
        return
	}

	 otp ,err := apiConf.db.GetOtpByUserId(r.Context(),user.ID)
	 if err != nil{
		respondWithError(w , 400 , fmt.Sprintf("Error with fetching user otp %v",err))
		return
	 }
    if time.Now().After(otp.ExpAt ){
		respondWithError(w ,400 ,"OTP is expired!")
		return
	}
	is_matched := VerifyOTP(otp.Otp,params.Otp)
	
	if !is_matched{
		respondWithError(w,401,"Incorrect OTP!")
		return
	}
	err = PasswordChecker(params.NewPassword)

    if err != nil{
		respondWithError(w , 400 , err.Error())
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

   apiConf.db.DeleteOtpByUserId(r.Context(),user.ID)
   apiConf.db.RevokeAllUserRefreshTokens(r.Context(),user.ID)
   respondWithJSON(w,200 , ResponseHealth("Your password updated successfully"))

}

func (apiConf apiConfig)LogOut(w http.ResponseWriter , r *http.Request , user database.User)  {
	access_token ,  err := auth.GetToken(r.Header)
	if err != nil{
		respondWithError(w , 401 , fmt.Sprintf("Auth Error %s" , err))
		return
	}

	_ , err  = apiConf.db.CreateTokenBlacklist(r.Context(),
	 database.CreateTokenBlacklistParams{
		ID: uuid.New(),
		UserID: user.ID,
		Token: access_token,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
		CreatedAt: time.Now().UTC(),
	 },
	)
	if err != nil {
	   respondWithError(w , 401 , fmt.Sprintf("Error with creating token blacklist %s" , err))
	   return
	}
	apiConf.db.RevokeAllUserRefreshTokens(r.Context() , user.ID)
}

func (apiConf apiConfig)DeleteUser(w http.ResponseWriter , r *http.Request , user database.User)  {
	idStr := chi.URLParam(r, "id")
	if idStr == "" {
		respondWithError(w, 400, "missing user id")
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		respondWithError(w, 400, "invalid UUID format")
		return
	}
	userToBeDeleted , err := apiConf.db.GetUserById(r.Context(),id)
    if err != nil {
		respondWithError(w,400,fmt.Sprintf("User not found %s" , err))
		return
	}

	err = apiConf.db.DeleteUserByUserId(r.Context() ,userToBeDeleted.ID)
    
	if err != nil {
		respondWithError(w , 400 , fmt.Sprintf("Error with deleting user %v" , err))
	    return
	}
}