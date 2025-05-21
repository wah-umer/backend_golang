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
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var wishlistColl = database.Collection("wishlists")

func GetWishlistByUserId(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Methods", "GET")

	params := mux.Vars(r)
	userID := params["id"]

	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	filter := bson.M{"user": id}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cur, err := wishlistColl.Find(ctx, filter)

	if err != nil {
		log.Println("DB error on Find:", err)
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}
	defer cur.Close(ctx)

	var entries []models.Wishlist
	if err := cur.All(ctx, &entries); err != nil {
		log.Println("Cursor.All error:", err)
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	//populate brand and category
	var resp []struct {
		models.Wishlist `json:",inline"`
		Product         models.Product `json:"product"`
		Brand           models.Brand   `json:"brand"`
	}
	prodColl := database.Collection("products")
	brandColl := database.Collection("brands")

	for _, e := range entries {
		// Fetch product
		var p models.Product
		if err := prodColl.FindOne(ctx, bson.M{"_id": e.Product}).Decode(&p); err != nil {
			log.Println("Product lookup error:", err)
		}
		// Fetch brand
		var b models.Brand
		if err := brandColl.FindOne(ctx, bson.M{"_id": p.Brand}).Decode(&b); err != nil {
			log.Println("Brand lookup error:", err)
		}
		resp = append(resp, struct {
			models.Wishlist `json:",inline"`
			Product         models.Product `json:"product"`
			Brand           models.Brand   `json:"brand"`
		}{e, p, b})
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

func UpdateWishlistById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Methods", "GET")

	params := mux.Vars(r)
	wishId := params["id"]

	id, err := primitive.ObjectIDFromHex(wishId)
	if err != nil {
		http.Error(w, "Invalid wishlist ID", http.StatusBadRequest)
		return
	}

	var updates map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	delete(updates, "_id")
	delete(updates, "user")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	update := bson.M{"$set": updates}

	var updated models.Wishlist
	err = wishlistColl.FindOneAndUpdate(ctx, bson.M{"_id": id}, update, opts).Decode(&updated)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			http.Error(w, "Wishlist entry not found", http.StatusNotFound)
		} else {
			log.Println("DB error on update:", err)
			http.Error(w, "Server error", http.StatusInternalServerError)
		}
		return
	}

	var product models.Product
	if err := database.Collection("products").
		FindOne(ctx, bson.M{"_id": updated.Product}).
		Decode(&product); err != nil {
		log.Println("Product lookup error:", err)
	}

	resp := struct {
		models.Wishlist `json:",inline"`
		Product         models.Product `json:"product"`
	}{
		Wishlist: updated,
		Product:  product,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

func DeleteWishlistById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Methods", "GET")

	params := mux.Vars(r)
	wishId := params["id"]

	id, err := primitive.ObjectIDFromHex(wishId)
	if err != nil {
		http.Error(w, "Invalid wishlist ID", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var deleted models.Wishlist
	err = wishlistColl.
		FindOneAndDelete(ctx, bson.M{"_id": id}).
		Decode(&deleted)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			http.Error(w, "Wishlist entry not found", http.StatusNotFound)
		} else {
			log.Println("DB error on delete:", err)
			http.Error(w, "Server error", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(deleted)
}
