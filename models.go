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

func DatabaseProductsToProducts(dbProduct []database.Product)[]Product  {
	products := []Product{}

	for _,product := range dbProduct{
		products = append(products, DatabaseProductToProduct(product))
	}
	return products
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
		CreatedAt: dbcat.CreatedAt,
		UpdatedAt: dbcat.UpdatedAt,
	}
}

func DatabaseCategorysToCategorys(dbCategory []database.Category)[]Category  {
	categories := []Category{}

	for _,category := range dbCategory{
		categories = append(categories, DatabaseCategoryToCategory(category))
	}
	return categories
}


type Cart struct{
	Id uuid.UUID `json:"id"`
	UserId uuid.UUID `json:"user_id"`
	Status string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func DatabaseCartToCart(dbcart database.Cart) Cart  {
	return Cart{
		Id: dbcart.ID,
		UserId: dbcart.UserID,
		Status: dbcart.Status,
		CreatedAt: dbcart.CreatedAt.Time,
		UpdatedAt: dbcart.UpdatedAt.Time,
	}
}

type CartItem struct{
	Id uuid.UUID `json:"id"`
	CartId uuid.UUID `json:"cart_id"`
	ProductID uuid.UUID    `json:"product_id"`
	Quantity  int       `json:"quantity"`
	PriceAtAdd string  `json:"price_at_add"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func DatabaseCartItemToCartItem(dbcartI database.CartItem) CartItem  {
	return CartItem{
		Id: dbcartI.ID,
		CartId:dbcartI.CartID,
		ProductID:dbcartI.ProductID,
		Quantity: int(dbcartI.Quantity),
		PriceAtAdd: dbcartI.PriceAtAdd,
		CreatedAt: dbcartI.CreatedAt.Time,
		UpdatedAt: dbcartI.UpdatedAt.Time,
	}
}

func DatabaseCartItemsToCartItems(dbcartIs []database.CartItem) []CartItem  {
	cart_items := []CartItem{}
	for _ , dbdbcartI := range dbcartIs{
		cart_items = append(cart_items , DatabaseCartItemToCartItem(dbdbcartI))
	}
	
	return cart_items
}

type Order struct{
	Id uuid.UUID `json:"id"`
    UserId uuid.UUID `json:"user_id"`
	OrderStatus string `json:"order_status"`
	Total   string `json:"total"`
	PaymentStatus string `json:"payment_status"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

func DatabaseOrderToOrder(dbOrder database.Order) Order {
	return Order{
		Id: dbOrder.ID,
		UserId: dbOrder.UserID,
		OrderStatus: dbOrder.OrderStatus,
		Total: dbOrder.Total,
		PaymentStatus: dbOrder.PaymentStatus.String,
		CreatedAt: dbOrder.CreatedAt.Time,
		UpdatedAt: dbOrder.UpdatedAt.Time,
	}
}


func DatabaseOrdersToOrders(dbOrders []database.Order) []Order {
	orders := []Order{}
	for _ , dbOrder := range dbOrders{
		orders = append(orders , DatabaseOrderToOrder(dbOrder))
	}
	
	return orders
}

type OrderItem struct{
	Id uuid.UUID `json:"id"`
    OrderId uuid.UUID `json:"order_id"`
	ProductId uuid.UUID `json:"product_id"`
	Quantity   int32 `json:"quantity"`
	UnitPrice   string `json:"unit_price"`
	TotalPrice   string `json:"total_price"`
	CreatedAt    time.Time  `json:"created_at"`
}

func DatabaseOrderItemToOrderItem(dbOrderItem database.OrderItem) OrderItem {
	return OrderItem{
		Id: dbOrderItem.ID,
		OrderId: dbOrderItem.OrderID,
		ProductId:dbOrderItem.ProductID,
		Quantity: dbOrderItem.Quantity,
		UnitPrice: dbOrderItem.UnitPrice,
		TotalPrice: dbOrderItem.TotalPrice,
		CreatedAt: dbOrderItem.CreatedAt.Time,
	}
}

func DatabaseOrderItemsToOrderItems(dbOrderItems []database.OrderItem) []OrderItem {
	orderItems := []OrderItem{}
	for _ , dbOrderItem := range dbOrderItems{
		orderItems = append(orderItems , DatabaseOrderItemToOrderItem(dbOrderItem))
	}
	
	return orderItems
}