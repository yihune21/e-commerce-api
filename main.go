package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	"github.com/yihune21/e-commerce-api/internal/database"
	"github.com/yihune21/e-commerce-api/utils"
)

type apiConfig struct{
	db *database.Queries
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
	fmt.Println("Database connected succefully!")


    //server router
	router := chi.NewRouter()
	router.Use(cors.Handler(
		cors.Options{
		AllowedOrigins: []string{"https://*","http://*"},
		AllowedMethods: []string{"GET","POST","DELETE","OPTIONS","PUT","PATCH"},
		AllowedHeaders: []string{"*"},
		ExposedHeaders: []string{"Link"},
		AllowCredentials: false,
		MaxAge:             300,
	}))

	v1Router := chi.NewRouter()

	v1Router.Get("/health",handlerHealthy)
	v1Router.Get("/err" , handleError)
	v1Router.Post("/user",apiCfg.New)
	v1Router.Post("/admin",apiCfg.AdminMiddlewareAuth(apiCfg.NewAdmin))
	v1Router.Get("/user",apiCfg.middlewareAuth(apiCfg.handlerGetUserByUserId))
	v1Router.Get("/login",apiCfg.Login)
	v1Router.Post("/logout",apiCfg.middlewareAuth(apiCfg.LogOut))
	v1Router.Post("/refreshToken",apiCfg.RefreshToken)
	v1Router.Patch("/update-password",apiCfg.middlewareAuth(apiCfg.UpdateUserPassword))
	v1Router.Post("/send-otp",apiCfg.RequestForgotPassword)
    v1Router.Post("/verify-otp",apiCfg.ForgotPassword)
	v1Router.Post("/delete-user",apiCfg.AdminMiddlewareAuth(apiCfg.DeleteUser))
	v1Router.Post("/product",apiCfg.AdminMiddlewareAuth(apiCfg.CreateProduct))
    v1Router.Get("/product",apiCfg.GetProductByName)
	v1Router.Patch("/product-price",apiCfg.AdminMiddlewareAuth(apiCfg.UpdateProductPrice))
	v1Router.Patch("/product-image",apiCfg.AdminMiddlewareAuth(apiCfg.UpdateProductImage))
  	v1Router.Delete("/product",apiCfg.AdminMiddlewareAuth(apiCfg.DeleteProduct))
    v1Router.Post("/category",apiCfg.AdminMiddlewareAuth(apiCfg.NewCategory))
	v1Router.Patch("/category",apiCfg.AdminMiddlewareAuth(apiCfg.UpdateCategoryName))
	v1Router.Post("/cart",apiCfg.middlewareAuth(apiCfg.AddToCart))

	
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