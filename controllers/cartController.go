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

var cartColl = database.Collection("carts")

func CreateCartItem(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Methods", "POST")

	var item models.Cart

	err := json.NewDecoder(r.Body).Decode(&item)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	item.ID = primitive.NewObjectID()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if _, err := cartColl.InsertOne(ctx, item); err != nil {
		log.Println("Insert error:", err)
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(item)
}

func GetCartByUserId(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Methods", "GET")

	params := mux.Vars(r)
	userId := params["id"]

	id, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"user": id}

	cur, err := cartColl.Find(ctx, filter)
	if err != nil {
		log.Println("Find error:", err)
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}
	defer cur.Close(ctx)

	var items []models.Cart
	if err := cur.All(ctx, &items); err != nil {
		log.Println("Cursor error:", err)
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	var resp []struct {
		models.Cart `json:",inline"`
		Product     models.Product `json:"product"`
		Brand       models.Brand   `json:"brand"`
	}

	prodColl := database.Collection("products")
	brandColl := database.Collection("brands")

	for _, e := range items {
		var p models.Product
		if err := prodColl.FindOne(ctx, bson.M{"_id": e.Product}).Decode(&p); err != nil {
			log.Println("Product lookup error:", err)
			continue
		}

		var b models.Brand
		if err := brandColl.FindOne(ctx, bson.M{"_id": p.Brand}).Decode(&b); err != nil {
			log.Println("Brand lookup error:", err)
		}

		resp = append(resp, struct {
			models.Cart `json:",inline"`
			Product     models.Product `json:"product"`
			Brand       models.Brand   `json:"brand"`
		}{e, p, b})
	}

	json.NewEncoder(w).Encode(resp)
}

func UpdateCartItem(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Methods", "PATCH")

	params := mux.Vars(r)
	cId := params["id"]

	id, err := primitive.ObjectIDFromHex(cId)
	if err != nil {
		http.Error(w, "Invalid cart ID", http.StatusBadRequest)
		return
	}

	var upd map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&upd); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	delete(upd, "_id")
	delete(upd, "user")
	delete(upd, "product")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	var updated models.Cart

	filter1 := bson.M{"_id": id}
	filter2 := bson.M{"$set": upd}
	err = cartColl.FindOneAndUpdate(ctx, filter1, filter2, opts).Decode(&updated)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			http.Error(w, "Not found", http.StatusNotFound)
		} else {
			log.Println("Update error:", err)
			http.Error(w, "Server error", http.StatusInternalServerError)
		}
		return
	}

	json.NewEncoder(w).Encode(updated)
}

func DeleteCartItem(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Methods", "DELETE")

	params := mux.Vars(r)
	cId := params["id"]

	id, err := primitive.ObjectIDFromHex(cId)
	if err != nil {
		http.Error(w, "Invalid cart ID", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var deleted models.Cart

	err = cartColl.FindOneAndDelete(ctx, bson.M{"_id": id}).Decode(&deleted)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			http.Error(w, "Not found", http.StatusNotFound)
		} else {
			log.Println("Delete error:", err)
			http.Error(w, "Server error", http.StatusInternalServerError)
		}
		return
	}

	json.NewEncoder(w).Encode(deleted)
}

func ClearCartByUserId(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	userId := params["id"]

	id, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"user": id}

	if _, err := cartColl.DeleteMany(ctx, filter); err != nil {
		log.Println("DeleteMany error:", err)
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
