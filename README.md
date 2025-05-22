# Go E-commerce Backend API

This repository contains a Go-based RESTful API for an e-commerce application, powered by MongoDB. It supports user management, product catalog, addresses, brands, categories, shopping cart, wishlist, reviews, and orders.

## 📦 Tech Stack

* **Language:** Go 1.20+
* **Web Framework:** net/http, [Gorilla Mux](https://github.com/gorilla/mux)
* **Database:** MongoDB (via [mongo-go-driver](https://github.com/mongodb/mongo-go-driver))
* **Environment Management:** [godotenv](https://github.com/joho/godotenv)

## 📁 Project Structure

```
├── controllers      # HTTP handler functions
├── models           # Data models (MongoDB documents)
├── routes           # Route registrations
├── database         # MongoDB connection and helpers
├── main.go          # Application entry point
└── go.mod           # Go modules
```

## ⚙️ Installation

1. **Clone the repo**

   ```bash
   git clone https://github.com/your-username/backend_golang.git
   cd backend_golang
   ```

2. **Install dependencies**

   ```bash
   go mod download
   ```

3. **Run the server**

   ```bash
   go run main.go
   ```

   The API will listen on port `5000` by default.

## 🔌 Environment Variables

| Key        | Description                  | Example                                        |
| ---------- | ---------------------------- | ---------------------------------------------- |
| MONGO\_URI | MongoDB connection string    | `mongodb+srv://user:pass@cluster0.mongodb.net` |
| DB\_NAME   | Name of the MongoDB database | `test`                                         |

## 📚 API Endpoints

> **Base URL:** `http://localhost:5000`

### Users

| Method | Endpoint      | Description        |
| ------ | ------------- | ------------------ |
| GET    | `/users/{id}` | Fetch a user by ID |

### Addresses

| Method | Endpoint             | Description                  |
| ------ | -------------------- | ---------------------------- |
| POST   | `/address`           | Create an address            |
| GET    | `/address/user/{id}` | Get all addresses for a user |
| PATCH  | `/address/{id}`      | Update an address by ID      |
| DELETE | `/address/{id}`      | Delete an address by ID      |

### Brands & Categories

| Method | Endpoint      | Description         |
| ------ | ------------- | ------------------- |
| GET    | `/brands`     | List all brands     |
| GET    | `/categories` | List all categories |

### Products

| Method | Endpoint         | Description                                |
| ------ | ---------------- | ------------------------------------------ |
| GET    | `/products`      | List products (with filtering, sorting)    |
| GET    | `/products/{id}` | Get product details, with brand & category |

### Wishlist

| Method | Endpoint               | Description                         |
| ------ | ---------------------- | ----------------------------------- |
| GET    | `/wishlists/user/{id}` | Get wishlist for a user (populated) |
| PATCH  | `/wishlists/{id}`      | Update a wishlist entry             |
| DELETE | `/wishlists/{id}`      | Delete a wishlist entry             |
| DELETE | `/wishlists/user/{id}` | Clear entire wishlist for a user    |

### Reviews

| Method | Endpoint                | Description                           |
| ------ | ----------------------- | ------------------------------------- |
| POST   | `/reviews`              | Create a product review               |
| GET    | `/reviews/product/{id}` | Get reviews for a product (populated) |
| PATCH  | `/reviews/{id}`         | Update a review                       |
| DELETE | `/reviews/{id}`         | Delete a review                       |

### Orders

| Method | Endpoint            | Description                       |
| ------ | ------------------- | --------------------------------- |
| POST   | `/orders`           | Create an order                   |
| GET    | `/orders/user/{id}` | Get orders for a user             |
| GET    | `/orders`           | List all orders (with pagination) |
| PATCH  | `/orders/{id}`      | Update an order status or details |

### Cart

| Method | Endpoint           | Description                    |
| ------ | ------------------ | ------------------------------ |
| POST   | `/carts`           | Add item to cart               |
| GET    | `/carts/user/{id}` | Get user’s cart (populated)    |
| PATCH  | `/carts/{id}`      | Update cart item (quantity)    |
| DELETE | `/carts/{id}`      | Remove single cart item        |
| DELETE | `/carts/user/{id}` | Clear all items in user’s cart |

## 🎯 Contributing

Contributions are welcome! Please open issues or submit pull requests for enhancements or bug fixes.

## 📄 License

This project is licensed under the MIT License. Feel free to use and modify.
