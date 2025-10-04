package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)
func main()  {
    godotenv.Load(".env")
	port := os.Getenv("PORT")
    db_url := os.Getenv("DB_URL")
    
    fmt.Println(db_url)

	fmt.Printf("E-commerce API Listen on %s\n" , port)
}