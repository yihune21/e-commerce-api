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
    godotenv.Load(".env")
	port := os.Getenv("PORT")
    db_url := os.Getenv("DB_URL")
	
	db_conn , err := utils.ConnectDb(db_url)
	if err != nil{
		log.Fatal(err)
	}
	apiCfg := apiConfig{
		db : db_conn,
	}
	fmt.Println("Database connected succefully!")


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
	v1Router.Post("/admin/delete-user/{id}",apiCfg.AdminMiddlewareAuth(apiCfg.DeleteUser))
	v1Router.Post("/admin/product",apiCfg.AdminMiddlewareAuth(apiCfg.CreateProduct))
    v1Router.Get("/product",apiCfg.GetProductByName)
	v1Router.Get("/products",apiCfg.GetAllProducts)
	v1Router.Patch("/admin/product-price",apiCfg.AdminMiddlewareAuth(apiCfg.UpdateProductPrice))
	v1Router.Patch("/admin/product-image",apiCfg.AdminMiddlewareAuth(apiCfg.UpdateProductImage))
  	v1Router.Delete("/admin/product/{id}",apiCfg.AdminMiddlewareAuth(apiCfg.DeleteProduct))
    v1Router.Post("/admin/category",apiCfg.AdminMiddlewareAuth(apiCfg.NewCategory))
	v1Router.Patch("/admin/category/{id}",apiCfg.AdminMiddlewareAuth(apiCfg.UpdateCategoryName))
	v1Router.Post("/cart",apiCfg.middlewareAuth(apiCfg.AddToCart))
	v1Router.Get("/cart", apiCfg.middlewareAuth(apiCfg.GetCart))
    v1Router.Delete("/cart/{productId}", apiCfg.middlewareAuth(apiCfg.RemoveFromCart))
    v1Router.Patch("/cart/{productId}", apiCfg.middlewareAuth(apiCfg.UpdateCartItem))
	v1Router.Post("/order",apiCfg.middlewareAuth(apiCfg.CreateOrder))
    v1Router.Get("/admin/orders",apiCfg.AdminMiddlewareAuth(apiCfg.GetAllOrders))
	v1Router.Get("/orders",apiCfg.middlewareAuth(apiCfg.GetAllOrdersByUserId))
    v1Router.Get("/order/{id}",apiCfg.middlewareAuth(apiCfg.GetOrderDetail))
    v1Router.Get("/admin/order/{id}",apiCfg.AdminMiddlewareAuth(apiCfg.GetOrderDetail))
    v1Router.Patch("/admin/order/{id}",apiCfg.AdminMiddlewareAuth(apiCfg.UpdateOrderStatus))




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