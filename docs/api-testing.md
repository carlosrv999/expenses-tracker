# API Testing with cURL

Base URL: `http://localhost:8080`

## Health Check

```bash
curl -s http://localhost:8080/health | jq
```

---

## Categories

All category endpoints are under `/api/v1/categories`.

### 1. List all categories

```bash
curl -s http://localhost:8080/api/v1/categories | jq
```

### 2. Create a top-level category

```bash
curl -s -X POST http://localhost:8080/api/v1/categories \
  -H "Content-Type: application/json" \
  -d '{
    "category_name": "Comida",
    "icon": "utensils",
    "color": "#FF5733"
  }' | jq
```

### 3. Create a subcategory (with parent)

```bash
curl -s -X POST http://localhost:8080/api/v1/categories \
  -H "Content-Type: application/json" \
  -d '{
    "parent_category_id": 1,
    "category_name": "Restaurantes",
    "icon": "restaurant",
    "color": "#C70039"
  }' | jq
```

### 4. Get a single category by ID

```bash
# Replace 1 with the actual category_id returned from create/list
curl -s http://localhost:8080/api/v1/categories/1 | jq
```

### 5. Update a category

```bash
# Replace 1 with the actual category_id
curl -s -X PUT http://localhost:8080/api/v1/categories/1 \
  -H "Content-Type: application/json" \
  -d '{
    "category_name": "Comida y Restaurantes",
    "icon": "utensils",
    "color": "#FF9800"
  }' | jq
```

### 6. Delete a category

```bash
# Replace 1 with the actual category_id
curl -s -X DELETE http://localhost:8080/api/v1/categories/1
```

## Expenses

All expense endpoints are under `/api/v1/expenses`.

**Important notes before testing:**

- You must have at least one **category** and one **payment method** created first.
- `amount` is stored as an integer in the smallest currency unit (e.g. **centavos** for PEN → 2500 = 25.00 PEN).
- `expense_date` must be a valid ISO 8601 timestamp (with timezone recommended).
- `tag_ids` is optional.

### 1. List all expenses (with optional filters)

```bash
# Basic list
curl -s http://localhost:8080/api/v1/expenses | jq

# Filter by category + limit
curl -s "http://localhost:8080/api/v1/expenses?category_id=1&limit=10" | jq

# Filter by payment method, limit and offset
curl -s "http://localhost:8080/api/v1/expenses?payment_method_id=1&offset=0&limit=5" | jq

# Filter start date and end date
curl -s "http://localhost:8080/api/v1/expenses?start_date=2025-04-01&end_date=2025-04-30" | jq
```

### 2. Create a new expense

```bash
curl -s -X POST http://localhost:8080/api/v1/expenses \
  -H "Content-Type: application/json" \
  -d '{
    "category_id": 1,
    "payment_method_id": 1,
    "currency": "PEN",
    "amount": 2500,
    "expense_date": "2026-04-26T19:30:00Z",
    "merchant_name": "Supermercado Plaza Vea",
    "description": "Compras semanales de la quincena",
    "tag_ids": [1, 3]
  }' | jq
```

### 2.1 Bulk upload expenses

```bash
curl -X POST http://localhost:8080/api/v1/expenses/upload \
  -F "file=@expenses.csv" \
  -H "Accept: application/json"
```

### 3. Get a single expense by ID

```bash
# Replace 1 with the actual expense_id from the create response
curl -s http://localhost:8080/api/v1/expenses/1 | jq
```

### 4. Update an expense

```bash
# Replace 1 with the actual expense_id
curl -s -X PUT http://localhost:8080/api/v1/expenses/1 \
  -H "Content-Type: application/json" \
  -d '{
    "category_id": 2,
    "payment_method_id": 1,
    "currency": "PEN",
    "amount": 3200,
    "expense_date": "2026-04-26T20:00:00Z",
    "merchant_name": "Restaurante La Bistecca",
    "description": "Almuerzo con cliente",
    "tag_ids": [2]
  }' | jq
```

### 5. Delete an expense

```bash
# Replace 1 with the actual expense_id
curl -s -X DELETE http://localhost:8080/api/v1/expenses/1
```

## Tags

All tag endpoints are under `/api/v1/tags`.

### 1. List all tags

```bash
curl -s http://localhost:8080/api/v1/tags | jq
```

### 2. Create a new tag

```bash
curl -s -X POST http://localhost:8080/api/v1/tags \
  -H "Content-Type: application/json" \
  -d '{
    "tag_name": "negocios",
    "color": "#4CAF50",
    "icon": "briefcase"
  }' | jq
```

### 3. Create another example tag

```bash
curl -s -X POST http://localhost:8080/api/v1/tags \
  -H "Content-Type: application/json" \
  -d '{
    "tag_name": "viajes",
    "color": "#2196F3",
    "icon": "plane"
  }' | jq
```

### 4. Get a single tag by ID

```bash
# Replace 1 with the actual tag_id returned from create/list
curl -s http://localhost:8080/api/v1/tags/1 | jq
```

### 5. Update a tag

```bash
# Replace 1 with the actual tag_id
curl -s -X PUT http://localhost:8080/api/v1/tags/1 \
  -H "Content-Type: application/json" \
  -d '{
    "tag_name": "negocios_importantes",
    "color": "#FF9800",
    "icon": "briefcase-business"
  }' | jq
```

### 6. Delete a tag

```bash
# Replace 1 with the actual tag_id
curl -s -X DELETE http://localhost:8080/api/v1/tags/1
```

## Payment Methods

All payment method endpoints are under `/api/v1/payment-methods`.

### 1. List all payment methods

```bash
curl -s http://localhost:8080/api/v1/payment-methods | jq
```

### 2. Create common payment methods (recommended first ones)

**Cash**

```bash
curl -s -X POST http://localhost:8080/api/v1/payment-methods \
  -H "Content-Type: application/json" \
  -d '{
    "method_name": "Efectivo",
    "method_type": "cash",
    "icon": "money-bill"
  }' | jq
```

**Yape**

```bash
curl -s -X POST http://localhost:8080/api/v1/payment-methods \
  -H "Content-Type: application/json" \
  -d '{
    "method_name": "Yape",
    "method_type": "yape",
    "icon": "mobile"
  }' | jq
```

**Plin**

```bash
curl -s -X POST http://localhost:8080/api/v1/payment-methods \
  -H "Content-Type: application/json" \
  -d '{
    "method_name": "Plin",
    "method_type": "plin",
    "icon": "mobile"
  }' | jq
```

**Credit Card**

```bash
curl -s -X POST http://localhost:8080/api/v1/payment-methods \
  -H "Content-Type: application/json" \
  -d '{
    "method_name": "Visa BCP",
    "method_type": "credit_card",
    "icon": "credit-card"
  }' | jq
```

### 3. Get a single payment method by ID

```bash
# Replace 1 with the actual payment_method_id from create/list
curl -s http://localhost:8080/api/v1/payment-methods/1 | jq
```

### 4. Update a payment method

```bash
# Replace 1 with the actual payment_method_id
curl -s -X PUT http://localhost:8080/api/v1/payment-methods/1 \
  -H "Content-Type: application/json" \
  -d '{
    "method_name": "Yape Personal",
    "method_type": "yape",
    "icon": "phone"
  }' | jq
```

### 5. Delete a payment method

```bash
# Replace 1 with the actual payment_method_id
curl -s -X DELETE http://localhost:8080/api/v1/payment-methods/1
```
