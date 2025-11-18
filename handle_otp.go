package main

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"net/http"
)

func generateSecureOTP(length int) string {
	otp := ""
	for i := 0; i < length; i++ {
		num, _ := rand.Int(rand.Reader, big.NewInt(10))
		otp += fmt.Sprintf("%d", num.Int64())
	}
	return otp
}

func VerifyOTP(sentOTP , recievedOTP string) bool {
	
	if sentOTP == recievedOTP {
		return true
	}
    return false
}

func sendOTPByEmail(w http.ResponseWriter , r *http.Request)  {
     //TODO
	  
}
