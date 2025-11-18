package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/yihune21/e-commerce-api/internal/auth"
	"github.com/yihune21/e-commerce-api/internal/database"
	jwtAuth "github.com/yihune21/e-commerce-api/jwt"
)


func (apiCfg apiConfig) AddToCart(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameters struct {
		ProductId uuid.UUID `json:"product_id"`
		Quantity  int       `json:"quantity"`
	}

	var params parameters
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error parsing json: %v", err))
		return
	}

	accessToken, _ := auth.GetToken(r.Header)
	userID, err := jwtAuth.ExtractUserIDFromToken(accessToken)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error extracting user id: %v", err))
		return
	}

	cart, err := apiCfg.db.GetCartByUserId(r.Context(), userID)
	if err != nil {
		cart, err = apiCfg.db.CreateCart(r.Context(), database.CreateCartParams{
			ID:     uuid.New(),
			UserID: userID,
			Status: "active",
		})
		if err != nil {
			respondWithError(w, 400, fmt.Sprintf("Couldn't create cart: %v", err))
			return
		}
	}

	product, err := apiCfg.db.GetProductById(r.Context(), params.ProductId)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't find product: %v", err))
		return
	}

	cartItem, err := apiCfg.db.GetCartItemByCartIdAndProductId(
		r.Context(),
		database.GetCartItemByCartIdAndProductIdParams{
			CartID:    cart.ID,
			ProductID: params.ProductId,
		},
	)

	if err == nil {
		newQty := cartItem.Quantity + int32(params.Quantity)
		_, err := apiCfg.db.UpdateCartItemQuantity(r.Context(),
			database.UpdateCartItemQuantityParams{
				ID:       cartItem.ID,
				Quantity: newQty,
			})
		if err != nil {
			respondWithError(w, 400, fmt.Sprintf("Couldn't update cart item: %v", err))
			return
		}

	} else {
		_, err = apiCfg.db.CreateCartItem(r.Context(), database.CreateCartItemParams{
			ID:         uuid.New(),
			CartID:     cart.ID,
			ProductID:  params.ProductId,
			Quantity:   int32(params.Quantity),
			PriceAtAdd: product.Price,
		})
		if err != nil {
			respondWithError(w, 400, fmt.Sprintf("Couldn't create cart item: %v", err))
			return
		}
	}

	items, err := apiCfg.db.GetCartItemByCartId(r.Context(), cart.ID)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't load cart items: %v", err))
		return
	}

	response := struct {
		Cart  Cart       `json:"cart"`
		Items []CartItem  `json:"items"`
	}{
		Cart:  DatabaseCartToCart(cart),
		Items: DatabaseCartItemsToCartItems(items),
	}

	respondWithJSON(w, 200, response)
}
