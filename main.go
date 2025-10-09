package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	"github.com/yihune21/e-commerce-api/utils"
)

type apiConfig struct{
	db *sql.DB
}

func main()  {
	//Load env var
    godotenv.Load(".env")
	port := os.Getenv("PORT")
    db_url := os.Getenv("DB_URL")
	
    //Begin DB connection
	db_conn , err := utils.ConnectDb(db_url)
	if err != nil{
		log.Fatal(err)
	}
	apiCfg := apiConfig{
		db : db_conn,
	}
	apiCfg.db.Ping()
	fmt.Println("Database connected succefully!")


    //server router
	router := chi.NewRouter()
	router.Use(cors.Handler(
		cors.Options{
		AllowedOrigins: []string{"https://*","http://*"},
		AllowedMethods: []string{"GET","POST","DELETE","OPTIONS"},
		AllowedHeaders: []string{"*"},
		ExposedHeaders: []string{"Link"},
		AllowCredentials: false,
		MaxAge:             300,
	}))

	v1Router := chi.NewRouter()

	v1Router.Get("/health",handlerHealthy)
	v1Router.Get("/err" , handelError)

	router.Mount("/v1",v1Router)


    srv := &http.Server{
		Handler: router,
		Addr:":" +  port,
	}	
	fmt.Printf("E-commerce server Listen on port %s \n" , port)
  
	err = srv.ListenAndServe()
	if err != nil{
		fmt.Printf("Server Listen and Serve error %s \n" , err)
	}

}