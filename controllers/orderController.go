package controllers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/umerwaheed/backend_golang/database"
	"github.com/umerwaheed/backend_golang/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var orderColl = database.Collection("orders")

func CreateOrder(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Methods", "POST")

	var ord models.Order
	err := json.NewDecoder(r.Body).Decode(&ord)
	if err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	ord.ID = primitive.NewObjectID()
	ord.CreatedAt = time.Now()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if _, err := orderColl.InsertOne(ctx, ord); err != nil {
		log.Println("DB error on create order:", err)
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(ord)
}

func GetOrdersByUserId(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Methods", "GET")

	params := mux.Vars(r)
	userID := params["id"]

	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"user": id}

	cur, err := orderColl.Find(ctx, filter)
	if err != nil {
		log.Println("DB error on Find orders by user:", err)
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}
	defer cur.Close(ctx)

	var orders []models.Order
	if err := cur.All(ctx, &orders); err != nil {
		log.Println("Cursor.All error:", err)
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(orders)
}
