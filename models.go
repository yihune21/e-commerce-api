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
	IsAdmin  bool `json:"is_admin"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`

}

func databaseUserToUser(dbuser database.User) User  {
	return User{
		Id: dbuser.ID,
		Name: dbuser.Name,
		Email: dbuser.Email,
		Password: dbuser.Password,
		IsAdmin: dbuser.IsAdmin,
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
 type Otp struct{
      Status string `json:"status"`
	  Otp    string `json:"otp"`
}
func OtpRes(status , otp string) Otp {
	return Otp{
		Status: status,
		Otp: otp,
	}
}
type Product struct{
	Name string `json:"name"`
	Description string `json:"description"`
	Price string `json:"price"`
	Stock int32 `json:"stock"`
	CategoryId uuid.UUID `json:"category_id"`
	ImageUrl string `json:"image_url"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

}
func DatabaseProductToProduct(dbProduct database.Product)Product  {
	return Product{
		Name  :dbProduct.Name,
		Description :dbProduct.Description.String,
		Price :dbProduct.Price,
		Stock  :dbProduct.Stock,
		CategoryId :dbProduct.CategoryID,
		ImageUrl :dbProduct.ImageUrl.String,
		CreatedAt: dbProduct.CreatedAt.Time,
		UpdatedAt: dbProduct.UpdatedAt.Time,
	}
}

type Category struct{
    Id   uuid.UUID `json:"id"`
	Name string `json:"name"`
	Description string `json:"description"`
    CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func DatabaseCategoryToCategory(dbcat database.Category) Category  {
	return Category{
		Id: dbcat.ID,
		Name: dbcat.Name,
		Description: dbcat.Description.String,
		CreatedAt: dbcat.CreatedAt.Time,
		UpdatedAt: dbcat.UpdatedAt.Time,
	}
}

