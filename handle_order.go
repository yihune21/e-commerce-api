package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/yihune21/e-commerce-api/internal/auth"
	"github.com/yihune21/e-commerce-api/internal/database"
	jwtAuth "github.com/yihune21/e-commerce-api/jwt"
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
    
    access_token , err := auth.GetToken(r.Header)
    if err != nil{
       respondWithError(w , 400,fmt.Sprintf("Error with getting access token %v",err))
       return
    }
    user_id,err := jwtAuth.ExtractUserIDFromToken(access_token)
    if err != nil{
        respondWithError(w , 400 , fmt.Sprintf("Error with extracting user id %v",err))
        return
    }
    cart , err := apiCfg.db.GetCartByUserId(r.Context(),user_id)

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
    for _,item :=range params.Items {

        apiCfg.db.DeleteCartItem(r.Context(), database.DeleteCartItemParams{
            CartID: cart.ID,
            ProductID: item.ProductID,
        })
    }
    respondWithJSON(w,200 , DatabaseOrderToOrder(order))
}

func (apiCfg apiConfig)GetAllOrders(w http.ResponseWriter , r *http.Request, admin database.User)  {
    orders ,  err := apiCfg.db.GetAllOrders(r.Context())
    if  err != nil {
        respondWithError(w,400 , fmt.Sprintf("Couldn't found orders %v",err))
        return
    }

    respondWithJSON(w , 200 , DatabaseOrdersToOrders(orders))
}

func (apiCfg apiConfig)GetAllOrdersByUserId(w http.ResponseWriter , r *http.Request, user database.User)  {
    orders ,  err := apiCfg.db.GetOrdersByUserId(r.Context(),user.ID)
    if  err != nil {
        respondWithError(w,400 , fmt.Sprintf("Couldn't found orders %v",err))
        return
    }

    respondWithJSON(w , 200 , DatabaseOrdersToOrders(orders))
}

func (apiCfg apiConfig)GetOrderDetail(w http.ResponseWriter , r *http.Request, user database.User)  {
    idStr := chi.URLParam(r,"id")
	if idStr == "" {
		respondWithError(w ,400 , "Missing order id")
		return
	}

	id ,err := uuid.Parse(idStr)
    if err != nil {
		respondWithError(w ,400 , fmt.Sprintf("Error with parsing order id %v" ,err))
		return
	}

    order ,  err := apiCfg.db.GetOrder(r.Context() ,id )
    if  err != nil {
        respondWithError(w,400 , fmt.Sprintf("Couldn't found order %v",err))
        return
    }

    respondWithJSON(w , 200 , DatabaseOrderToOrder(order))
}

func (apiCfg apiConfig)UpdateOrderStatus(w http.ResponseWriter , r *http.Request, admin database.User)  {
    idStr := chi.URLParam(r,"id")
	if idStr == "" {
		respondWithError(w ,400 , "Missing order id")
		return
	}

	id ,err := uuid.Parse(idStr)
    if err != nil {
		respondWithError(w ,400 , fmt.Sprintf("Error with parsing order id %v" ,err))
		return
	}
    type parameters struct{
        Status string  `json:"status"`
    }

    decode := json.NewDecoder(r.Body)
    params := parameters{}
    err = decode.Decode(&params)
    if err != nil {
        respondWithError(w , 400,fmt.Sprintf("Error with parsing a json %v",err))
        return
    }

    order ,  err := apiCfg.db.UpdateOrderStatus(r.Context() ,database.UpdateOrderStatusParams{
        OrderStatus: params.Status,
        ID: id,
    } )
    if  err != nil {
        respondWithError(w,400 , fmt.Sprintf("Couldn't found order %v",err))
        return
    }

    respondWithJSON(w , 200 , DatabaseOrderToOrder(order))
}

func (apiConf apiConfig)OrderPagination(w http.ResponseWriter , r *http.Request ,user database.User)  {
     type parameters struct{
          OrdersPerPage int32 `json:"orders_per_page"`
	 }

	decode := json.NewDecoder(r.Body)
	params := parameters{}

	err := decode.Decode(&params)
    
	if err != nil {
		respondWithError(w ,400 , fmt.Sprintf("Error with parsing json %v" ,err))
		return
	}
	orders , err := apiConf.db.GetOrdersPerPage(r.Context() , params.OrdersPerPage)
    
	if err != nil {
		respondWithError(w ,400 , fmt.Sprintf("Couldn't found order %v" ,err))
		return
	}

	respondWithJSON(w , 200 , DatabaseOrdersToOrders(orders))

}