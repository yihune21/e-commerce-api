# Cash-on-Delivery (COD) Workflow

## Overview
This e-commerce API now supports Cash-on-Delivery (COD) payment method where customers can order products and pay when the delivery person arrives at their location.

## Database Changes

### New Tables
1. **shipping_addresses** - Stores delivery address for each order
   - Full name, phone, address details
   - Delivery instructions

### Updated Tables
1. **orders**
   - `payment_method`: 'cod', 'online', 'card'
   - `delivery_status`: 'pending', 'confirmed', 'out_for_delivery', 'delivered', 'failed'

2. **users**
   - `phone`: Contact number for delivery

## COD Order Flow

### 1. Customer Places Order
```json
POST /v1/order
{
  "items": [
    {
      "product_id": "uuid",
      "qty": 2
    }
  ],
  "payment_method": "cod",
  "shipping_address": {
    "full_name": "John Doe",
    "phone": "+251912345678",
    "address_line1": "123 Main St",
    "address_line2": "Apt 4B",
    "city": "Addis Ababa",
    "state": "Addis Ababa",
    "postal_code": "1000",
    "country": "Ethiopia",
    "delivery_instructions": "Call when you arrive"
  }
}
```

**Result:**
- Order created with `order_status: 'confirmed'`
- Payment status: `'pending'`
- Delivery status: `'confirmed'`
- Shipping address saved

### 2. Delivery Person Views Pending Deliveries
```
GET /v1/delivery/pending
```
Returns all confirmed COD orders ready for delivery.

### 3. Delivery Person Picks Up Order
```json
PATCH /v1/delivery/order/{id}/status
{
  "status": "out_for_delivery"
}
```

### 4. Delivery Completion
When the delivery person arrives and collects payment:

```json
POST /v1/delivery/order/{id}/complete-payment
```

**This automatically:**
- Sets `payment_status` to `'completed'`
- Sets `order_status` to `'delivered'`
- Sets `delivery_status` to `'delivered'`

### 5. Failed Delivery
If customer is not available or refuses:

```json
PATCH /v1/delivery/order/{id}/status
{
  "status": "failed"
}
```

## Status Definitions

### Order Status
- `pending`: Initial state
- `confirmed`: Order confirmed, ready for processing
- `delivered`: Order successfully delivered
- `cancelled`: Order cancelled

### Payment Status
- `pending`: Payment not received
- `completed`: Payment received
- `failed`: Payment failed

### Delivery Status
- `pending`: Not yet assigned for delivery
- `confirmed`: Ready for delivery pickup
- `out_for_delivery`: Delivery person has the order
- `delivered`: Successfully delivered
- `failed`: Delivery attempt failed

## User Registration with Phone
```json
POST /v1/user
{
  "name": "John Doe",
  "email": "john@example.com",
  "phone": "+251912345678",
  "password": "SecurePassword123!"
}
```

## Key Features
✅ No online payment required for COD orders
✅ Automatic status updates when payment is collected
✅ Shipping address stored with each order
✅ Phone contact for delivery coordination
✅ Delivery tracking through status updates
✅ Support for both COD and online payment methods

## Security Considerations
- Only authenticated users can create orders
- Delivery endpoints require authentication
- Payment completion only works for COD orders
- Cannot complete payment twice for the same order