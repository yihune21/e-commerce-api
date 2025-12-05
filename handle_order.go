package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/yihune21/e-commerce-api/internal/database"
)

func(apiCfg apiConfig) CreateOrder(w http.ResponseWriter , r *http.Request, user database.User)  {
	  
    type parameters struct {
         Items []struct {
            ProductID uuid.UUID `json:"product_id"`
            Qty       int       `json:"qty"`
        } `json:"items"`
    }

    decode := json.NewDecoder(r.Body)
    params := parameters{}
    err := decode.Decode(&params)

    if err != nil {
        respondWithError(w , 400,fmt.Sprintf("Error with parsing a json %v",err))
        return
    }

    total := 0.0
    for _,item :=range params.Items {

        product ,err := apiCfg.db.GetProductById(r.Context(), item.ProductID)
        if err != nil {
            respondWithError(w , 400,fmt.Sprintf("product not found %v",err))
            return
        }
        if product.Stock < int32(item.Qty){
           respondWithError(w , 400,"Out of stock sorry!")
           return
        }
        priceToInt , err := strconv.ParseFloat(product.Price,64)
        total  += priceToInt * float64(item.Qty)

    }
    
    created_at := sql.NullTime{}
	created_at.Time = time.Now().UTC()
	created_at.Valid = true
	updated_at := sql.NullTime{}
	updated_at.Time = time.Now().UTC()
	updated_at.Valid = true
    payment_status := sql.NullString{}
    payment_status.String = "pending"
    payment_status.Valid = true
	
    order , err := apiCfg.db.CreateOrder(r.Context() , database.CreateOrderParams{
        ID: uuid.New(),
        UserID: user.ID,
        OrderStatus: "pending",
        Total: strconv.FormatFloat(total,'g',-1,64),
        PaymentStatus: payment_status,
        CreatedAt: created_at,
        UpdatedAt: updated_at,        
    })

    if err != nil {
        respondWithError(w , 400,fmt.Sprintf("Couldn't create an order %v",err))
        return
    }

    
    for _,item :=range params.Items {
        
        db_product ,_:= apiCfg.db.GetProductById(r.Context(), item.ProductID)
        
        apiCfg.db.CreateOrderItem(r.Context() , database.CreateOrderItemParams{
            ID: uuid.New(),
            OrderID: order.ID,
            ProductID: item.ProductID,
        })
       
        apiCfg.db.UpdateProductStock(r.Context() , database.UpdateProductStockParams{
            Stock: db_product.Stock - int32(item.Qty),
            ID: item.ProductID,
        })
        

    }



    
    
}