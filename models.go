package main

import (
	"time"

	"github.com/google/uuid"
	"github.com/yihune21/e-commerce-api/internal/database"
)




type User struct{
	Id uuid.UUID `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Password string `json:"password"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`

}

func databaseUserToUser(dbuser database.User) User  {
	return User{
		Id: dbuser.ID,
		Name: dbuser.Name,
		Email: dbuser.Email,
		Password: dbuser.Password,
		CreatedAt: dbuser.CreatedAt,
		UpdatedAt: dbuser.UpdatedAt,
	}
}

type Token  struct{
	AccessToken string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func ResponseToken(accessToken string, refreshToken string) Token {
	return Token{
		AccessToken: accessToken,
		RefreshToken: refreshToken,
	}
}

type HealthRes  struct{
	Status string `json:"status"`
}
func ResponseHealth(msg string) HealthRes {
	return HealthRes{
		 Status: msg,
	}
}
