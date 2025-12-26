package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/yihune21/e-commerce-api/internal/database"
)

// UpdateDeliveryStatus - For delivery person to update delivery status
func (apiCfg apiConfig) UpdateDeliveryStatus(w http.ResponseWriter, r *http.Request, user database.User) {
	idStr := chi.URLParam(r, "id")
	if idStr == "" {
		respondWithError(w, 400, "Missing order id")
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error with parsing order id %v", err))
		return
	}

	type parameters struct {
		Status string `json:"status"` // 'out_for_delivery', 'delivered', 'failed'
	}

	decode := json.NewDecoder(r.Body)
	params := parameters{}
	err = decode.Decode(&params)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error with parsing json %v", err))
		return
	}

	// Validate status
	validStatuses := map[string]bool{
		"out_for_delivery": true,
		"delivered":        true,
		"failed":           true,
	}

	if !validStatuses[params.Status] {
		respondWithError(w, 400, "Invalid delivery status. Must be: out_for_delivery, delivered, or failed")
		return
	}

	// Get the order first to check payment method
	order, err := apiCfg.db.GetOrder(r.Context(), id)
	if err != nil {
		respondWithError(w, 404, fmt.Sprintf("Order not found %v", err))
		return
	}

	updated_at := sql.NullTime{}
	updated_at.Time = time.Now().UTC()
	updated_at.Valid = true

	// Update delivery status
	updatedOrder, err := apiCfg.db.UpdateDeliveryStatus(r.Context(), database.UpdateDeliveryStatusParams{
		DeliveryStatus: sql.NullString{String: params.Status, Valid: true},
		ID:             id,
		UpdatedAt:      updated_at,
	})

	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't update delivery status %v", err))
		return
	}

	// If delivered and it's a COD order, automatically complete the payment
	if params.Status == "delivered" && order.PaymentMethod.String == "cod" {
		updatedOrder, err = apiCfg.db.UpdateOrderStatus(r.Context(), database.UpdateOrderStatusParams{
			OrderStatus: "delivered",
			ID:          id,
			UpdatedAt:   updated_at,
		})
		if err != nil {
			respondWithError(w, 400, fmt.Sprintf("Couldn't update order status %v", err))
			return
		}
	}

	respondWithJSON(w, 200, DatabaseOrderToOrder(updatedOrder))
}

// CompletePaymentOnDelivery - For delivery person to mark payment as received
func (apiCfg apiConfig) CompletePaymentOnDelivery(w http.ResponseWriter, r *http.Request, user database.User) {
	idStr := chi.URLParam(r, "id")
	if idStr == "" {
		respondWithError(w, 400, "Missing order id")
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error with parsing order id %v", err))
		return
	}

	// Get the order first
	order, err := apiCfg.db.GetOrder(r.Context(), id)
	if err != nil {
		respondWithError(w, 404, fmt.Sprintf("Order not found %v", err))
		return
	}

	// Check if it's a COD order
	if order.PaymentMethod.String != "cod" {
		respondWithError(w, 400, "This endpoint is only for COD orders")
		return
	}

	// Check if payment is already completed
	if order.PaymentStatus.String == "completed" {
		respondWithError(w, 400, "Payment already completed")
		return
	}

	updated_at := sql.NullTime{}
	updated_at.Time = time.Now().UTC()
	updated_at.Valid = true

	// Complete the payment and mark order as delivered
	updatedOrder, err := apiCfg.db.CompletePaymentOnDelivery(r.Context(), database.CompletePaymentOnDeliveryParams{
		ID:        id,
		UpdatedAt: updated_at,
	})

	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't complete payment %v", err))
		return
	}

	respondWithJSON(w, 200, DatabaseOrderToOrder(updatedOrder))
}

// GetOrderWithShippingAddress - Get order details with shipping address
func (apiCfg apiConfig) GetOrderWithShippingAddress(w http.ResponseWriter, r *http.Request, user database.User) {
	idStr := chi.URLParam(r, "id")
	if idStr == "" {
		respondWithError(w, 400, "Missing order id")
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error with parsing order id %v", err))
		return
	}

	// Get order
	order, err := apiCfg.db.GetOrder(r.Context(), id)
	if err != nil {
		respondWithError(w, 404, fmt.Sprintf("Order not found %v", err))
		return
	}

	// Get shipping address
	shippingAddress, err := apiCfg.db.GetShippingAddressByOrderId(r.Context(), id)
	if err != nil {
		respondWithError(w, 404, fmt.Sprintf("Shipping address not found %v", err))
		return
	}

	// Get order items
	orderItems, err := apiCfg.db.GetOrderItemsByOrderID(r.Context(), id)
	if err != nil {
		orderItems = []database.OrderItem{} // Empty array if no items
	}

	// Create response structure
	type OrderDetailsResponse struct {
		Order           Order           `json:"order"`
		ShippingAddress ShippingAddress `json:"shipping_address"`
		OrderItems      []OrderItem     `json:"order_items"`
	}

	response := OrderDetailsResponse{
		Order:           DatabaseOrderToOrder(order),
		ShippingAddress: DatabaseShippingAddressToShippingAddress(shippingAddress),
		OrderItems:      DatabaseOrderItemsToOrderItems(orderItems),
	}

	respondWithJSON(w, 200, response)
}

// GetPendingDeliveries - Get all orders ready for delivery
func (apiCfg apiConfig) GetPendingDeliveries(w http.ResponseWriter, r *http.Request, user database.User) {
	// This would typically be restricted to delivery personnel or admin
	// For now, we'll get all confirmed orders with COD payment method
	
	orders, err := apiCfg.db.GetAllOrders(r.Context())
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't fetch orders %v", err))
		return
	}

	// Filter for confirmed COD orders
	var pendingDeliveries []Order
	for _, order := range orders {
		if order.DeliveryStatus.String == "confirmed" && 
		   order.PaymentMethod.String == "cod" &&
		   order.PaymentStatus.String == "pending" {
			pendingDeliveries = append(pendingDeliveries, DatabaseOrderToOrder(order))
		}
	}

	respondWithJSON(w, 200, pendingDeliveries)
}