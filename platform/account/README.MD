Account
=======

Generate Pub/Sub Keys for JWT Auth
----------------------------------

```bash
openssl genrsa -out private.pem 2048
openssl rsa -in private.pem -pubout -out public.pem

## Better naming

openssl genrsa -out cert_private.pem 2048
openssl rsa -in cert_private.pem -pubout -out cert_public.pem

```

Use Login with JWT:

```bash
# 1. Get token from auth
curl -X POST http://localhost:8081/auth/token \
  -H "Content-Type: application/json" \
  -d '{"username":"test", "password":"password", "client_id":"game_1"}'

# 2. Use token with storage
curl -H "Authorization: Bearer <PASTE_TOKEN>" http://localhost:8082/storage/secure

curl -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJjbGllbnRfaWQiOiJnYW1lXzEueC55LnoiLCJleHAiOjE3NTAxNTEyNTEsImlzcyI6ImdhbWUtaGFuZ2FyIiwicm9sZSI6InVzZXIiLCJzdWIiOiJ1c2VyXzEifQ.DD28dMwDQQvQFqt3s4l0JXq6mVI1cnTiuKr7KJqyaFE" 

curl --location 'localhost:8000/api/v1/account/user_1?full_profile=false' \
--header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJjbGllbnRfaWQiOiJnYW1lXzEueC55LnoiLCJleHAiOjE3NTAxNTEyNTEsImlzcyI6ImdhbWUtaGFuZ2FyIiwicm9sZSI6InVzZXIiLCJzdWIiOiJ1c2VyXzEifQ.DD28dMwDQQvQFqt3s4l0JXq6mVI1cnTiuKr7KJqyaFE' 


```


Related
-------

* [code-heim/go_36_gin_jwt](https://github.com/code-heim/go_36_gin_jwt)
  * [#36 Golang - Mastering JWT Authentication](https://www.youtube.com/watch?v=XFxm6NEWakQ)

Hexagonal / DDD Architecture
----------------------------

**Flow**

```bash
HTTP Request → Controller → Handler → Service → Domain → Repository
```

```bash
├── internal/
│   ├── domain/              # Business entities and interfaces
│   │   └── user.go
│   │
│   ├── application/         # Application layer
│   │   ├── handlers/        # Application handlers (use case handlers)
│   │   │   └── user_handler.go
│   │   ├── services/        # Domain services
│   │   │   └── user_service.go
│   │   └── dto/            # Data Transfer Objects
│   │       └── user_dto.go
│   │
│   └── infrastructure/      # Infrastructure layer
│       ├── http/           
│       │   ├── controllers/ # HTTP controllers
│       │   └── routes/      # Route definitions
│       └── repositories/    # Repository implementations

```
1. **Clear Separation of Concerns**:
    - Controllers handle HTTP-specific logic (request parsing, status codes, etc.)
    - Handlers contain application logic
    - Services contain domain logic

2. **Better Testability**:
    - You can test handlers without HTTP concerns
    - Controllers can be tested for HTTP-specific behavior
    - Easier to mock dependencies

3. **Clean Dependencies**:
    - Application layer doesn't know about HTTP
    - Controllers depend on application, not vice versa


### Domain Service vs. Application Service

#### Domain Service

Domain Services handle **core business rules** that:
- Operate on domain entities
- Are independent of use cases
- Would be true regardless of the application

```go
// internal/domain/services/pricing_service.go
package services

import "your-project/internal/domain"

// PricingService is a Domain Service
// It handles core business rules about how prices should be calculated
type PricingService struct{}

func (s *PricingService) CalculateOrderTotal(order *domain.Order) domain.Money {
    var total domain.Money

    for _, item := range order.Items {
        // Core business rule: Basic price calculation
        itemTotal := item.Price.Multiply(item.Quantity)
        
        // Core business rule: Bulk discount
        if item.Quantity > 10 {
            itemTotal = itemTotal.ApplyDiscount(0.10) // 10% off for bulk orders
        }
        
        total = total.Add(itemTotal)
    }

    return total
}
```


#### Application Service

Application Services orchestrate **use cases** that:
- Coordinate multiple operations
- Handle transaction boundaries
- Deal with technical concerns
- Implement specific user stories/features


```go
// internal/application/services/order_service.go
package services

import (
    "context"
    "your-project/internal/domain"
    "your-project/internal/domain/services"
)

// OrderService is an Application Service
// It orchestrates the process of placing an order
type OrderService struct {
    orderRepo       domain.OrderRepository
    inventoryRepo   domain.InventoryRepository
    pricingService  *services.PricingService    // Domain service
    notifier        domain.NotificationService
}

// PlaceOrder is a use case that coordinates multiple operations
func (s *OrderService) PlaceOrder(ctx context.Context, userID string, items []OrderItem) (*domain.Order, error) {
    // Start transaction
    tx, err := s.beginTransaction()
    if err != nil {
        return nil, err
    }
    defer tx.Rollback()

    // 1. Create order entity
    order := domain.NewOrder(userID, items)

    // 2. Use domain service to calculate price
    total := s.pricingService.CalculateOrderTotal(order)
    order.SetTotal(total)

    // 3. Check inventory availability (application-level concern)
    if err := s.checkInventory(items); err != nil {
        return nil, err
    }

    // 4. Save order
    if err := s.orderRepo.Save(order); err != nil {
        return nil, err
    }

    // 5. Update inventory
    if err := s.updateInventory(items); err != nil {
        return nil, err
    }

    // 6. Send confirmation (technical concern)
    if err := s.notifier.SendOrderConfirmation(order); err != nil {
        // Log but don't fail the transaction
        log.Printf("Failed to send confirmation: %v", err)
    }

    // 7. Commit transaction
    if err := tx.Commit(); err != nil {
        return nil, err
    }

    return order, nil
}

func (s *OrderService) checkInventory(items []OrderItem) error {
    // Application-level coordination
    return nil
}

func (s *OrderService) updateInventory(items []OrderItem) error {
    // Application-level coordination
    return nil
}
```

### Example 2

#### Domain Service

```go
// internal/domain/services/discount_service.go
package services

// DiscountService handles core business rules about discounts
type DiscountService struct{}

func (s *DiscountService) CalculateDiscount(customer *Customer, order *Order) Money {
	// Core business rules that are true regardless of use case:

	// Rule 1: VIP customers get 10% off
	if customer.IsVIP() {
		return order.Total.Multiply(0.10)
	}

	// Rule 2: Orders over $100 get 5% off
	if order.Total.IsGreaterThan(Money{100}) {
		return order.Total.Multiply(0.05)
	}

	// Rule 3: First-time customers get 15% off
	if customer.IsFirstPurchase() {
		return order.Total.Multiply(0.15)
	}

	return Money{0}
}

```


#### Application Service

```go
// internal/application/services/checkout_service.go
package services

type CheckoutService struct {
	orderRepo       OrderRepository
	discountService *domain.DiscountService
	paymentGateway  PaymentGateway
	emailService    EmailService
}

// ProcessCheckout is a use case that uses domain services
func (s *CheckoutService) ProcessCheckout(ctx context.Context, checkoutRequest CheckoutRequest) (*Order, error) {
	// 1. Load customer
	customer, err := s.customerRepo.FindByID(checkoutRequest.CustomerID)
	if err != nil {
		return nil, err
	}

	// 2. Create order
	order := NewOrder(customer, checkoutRequest.Items)

	// 3. Apply discounts using domain service
	discount := s.discountService.CalculateDiscount(customer, order)
	order.ApplyDiscount(discount)

	// 4. Process payment (technical concern)
	if err := s.paymentGateway.ProcessPayment(order); err != nil {
		return nil, err
	}

	// 5. Save order
	if err := s.orderRepo.Save(order); err != nil {
		return nil, err
	}

	// 6. Send confirmation email (technical concern)
	go s.emailService.SendOrderConfirmation(order)

	return order, nil
}
```

Key Differences:
1. **Domain Services**:
    - Contain pure business rules
    - Don't know about databases, emails, or external services
    - Could be explained to a business expert
    - Are stateless and side-effect free
    - Example rules:
        - How to calculate prices
        - How to apply discounts
        - How to validate business rules
        - How to check if an order is eligible for shipping

2. **Application Services**:
    - Orchestrate use cases
    - Handle technical concerns (transactions, emails, etc.)
    - Coordinate multiple operations
    - Know about repositories and external services
    - Example operations:
        - Processing a checkout
        - Handling user registration
        - Coordinating order fulfillment
        - Managing inventory updates
