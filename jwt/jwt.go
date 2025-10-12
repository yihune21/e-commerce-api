package jwtAuth

import (
	"crypto/rsa"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/yihune21/e-commerce-api/internal/database"
)
func readPrivateKey() *rsa.PrivateKey {
	privateKeyData, err :=  os.ReadFile("keys/private.pem")
	if err != nil {
		panic(err)
	}
	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privateKeyData)
	if err != nil {
		panic(err)
	}

   return privateKey
}
func readPublicKey() *rsa.PublicKey {
	publicKeyData,  err := os.ReadFile("keys/public.pem")

	if err != nil{
		panic(err)
	}
	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(publicKeyData)
	if err != nil {
		panic(err)
	}

   return publicKey
}

func GenerateToken(user database.User) string {
	 private_key := readPrivateKey() 
	 
     claims := jwt.MapClaims{
        "sub":user.ID.String(),
		"exp": time.Now().Add(300 * time.Second).Unix(),
	 }
	 token := jwt.NewWithClaims(jwt.SigningMethodRS256 ,claims )
	 access_token ,err:= token.SignedString(private_key)
	 if err != nil {
		panic(err)
	 }
	 fmt.Printf("Access token , %v\n" ,access_token)
	 return access_token
}

func VerfiyToken(tokenString string) bool {
	public_key := readPublicKey()
	 
	parsedToken, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return public_key, nil
	})
	if err != nil {
		panic(err)
	}

	return  parsedToken.Valid
}
func ExtractUserIDFromToken(tokenString string) (uuid.UUID, error) {
	publicKey := readPublicKey()
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return publicKey, nil
	})
	if err != nil {
		return uuid.Nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if sub, ok := claims["sub"].(string); ok {
			userID, err := uuid.Parse(sub)
			if err != nil {
				return uuid.Nil, err
			}
			return userID, nil
		}
		return uuid.Nil, fmt.Errorf("sub claim missing or invalid")
	}
	return uuid.Nil, fmt.Errorf("invalid token")
}