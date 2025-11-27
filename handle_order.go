package main

import (
	"net/http"

	"github.com/yihune21/e-commerce-api/internal/database"
)

func NewOrder(w http.ResponseWriter , r *http.Request, user database.User)  {
	   
	// id UUID PRIMARY KEY ,
    // user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    // status VARCHAR(20) NOT NULL DEFAULT 'pending',
    // total NUMERIC(10,2) NOT NULL,    
    // created_at TIMESTAMP DEFAULT NOW(),
    // updated_at TIMESTAMP DEFAULT NOW()
}