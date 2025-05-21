package routes

import (
	"github.com/gorilla/mux"
	"github.com/umerwaheed/backend_golang/controllers"
)

func Router() *mux.Router {
	router := mux.NewRouter()

	//user routes
	router.HandleFunc("/users/{id}", controllers.GetById).Methods("GET")

	//address routes
	router.HandleFunc("/address", controllers.Create).Methods("POST")
	router.HandleFunc("/address/user/{id}", controllers.GetByUserId).Methods("GET")
	router.HandleFunc("/address/{id}", controllers.UpdateById).Methods("PATCH")
	router.HandleFunc("/address/{id}", controllers.DeleteById).Methods("DELETE")

	//brand & category routes
	router.HandleFunc("/brands", controllers.GetAllBrand).Methods("GET")
	router.HandleFunc("/categories", controllers.GetAllCategory).Methods("GET")

	// products routes
	router.HandleFunc("/products", controllers.GetAllProducts).Methods("GET")
	router.HandleFunc("/products/{id}", controllers.GetProductById).Methods("GET")

	// wishlist routes
	router.HandleFunc("/wishlist", controllers.UpdateWishlistById).Methods("POST")
	router.HandleFunc("/wishlist/user/{id}", controllers.GetWishlistByUserId).Methods("GET")
	router.HandleFunc("/wishlist/{id}", controllers.DeleteWishlistById).Methods("DELETE")

	// review route
	router.HandleFunc("/reviews", controllers.CreateReview).Methods("POST")
	router.HandleFunc("/reviews/product/{id}", controllers.GetReviewByProductId).Methods("GET")
	router.HandleFunc("/reviews/{id}", controllers.UpdateReviewById).Methods("PATCH")
	router.HandleFunc("/reviews/{id}", controllers.DeleteReviewById).Methods("DELETE")

	// order route
	router.HandleFunc("/orders", controllers.CreateOrder).Methods("POST")
	router.HandleFunc("/orders/user/{id}", controllers.GetOrdersByUserId).Methods("GET")

	// cart route
	router.HandleFunc("/cart", controllers.CreateCartItem).Methods("POST")
	router.HandleFunc("/cart/user/{id}", controllers.GetCartByUserId).Methods("GET")
	router.HandleFunc("/cart/user/{id}", controllers.ClearCartByUserId).Methods("DELETE")
	router.HandleFunc("/cart/{id}", controllers.UpdateCartItem).Methods("PATCH")
	router.HandleFunc("/cart/{id}", controllers.DeleteCartItem).Methods("DELETE")

	return router
}
