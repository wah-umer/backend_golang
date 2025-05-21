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

var prodColl = database.Collection("products")

func GetProductById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Methods", "GET")

	params := mux.Vars(r)
	prodId := params["id"]

	id, err := primitive.ObjectIDFromHex(prodId)
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	filter := bson.M{"_id": id}

	var prod models.Product

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	//fetch prod
	err = prodColl.FindOne(ctx, filter).Decode(&prod)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			http.Error(w, "Product not found", http.StatusNotFound)
			return
		}
		log.Println("DB error:", err)
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	//fetch brand
	filterbrand := bson.M{"_id": prod.Brand}
	var brand models.Brand
	err = database.Collection("brands").FindOne(ctx, filterbrand).Decode(&brand)
	if err != nil {
		log.Println("Brand lookup error:", err)
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	//fetch category
	filtercategory := bson.M{"_id": prod.Category}
	var category models.Category
	err = database.Collection("categories").FindOne(ctx, filtercategory).Decode(&category)
	if err != nil {
		log.Println("Category lookup error:", err)
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	// build combines response now
	resp := struct {
		models.Product `json:",inline"`
		Brand          models.Brand    `json:"brand"`
		Category       models.Category `json:"category"`
	}{
		Product:  prod,
		Brand:    brand,
		Category: category,
	}

	//send response
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

func GetAllProducts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Methods", "GET")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	q := r.URL.Query()
	filter := bson.M{}

	// filter by brand
	if brands := q["brand"]; len(brands) > 0 {
		var ids []primitive.ObjectID
		for _, b := range brands {
			if id, err := primitive.ObjectIDFromHex(b); err == nil {
				ids = append(ids, id)
			}
		}
		if len(ids) > 0 {
			filter["brand"] = bson.M{"$in": ids}
		}
	}

	// filter by category
	if cat := q["category"]; len(cat) > 0 {
		var ids []primitive.ObjectID
		for _, c := range cat {
			if id, err := primitive.ObjectIDFromHex(c); err == nil {
				ids = append(ids, id)
			}
		}
		if len(ids) > 0 {
			filter["category"] = bson.M{"$in": ids}
		}
	}

	// remove deleted
	if _, ok := q["user"]; ok {
		filter["isDeleted"] = false
	}

	// build find options
	findOpts := options.Find()
	if sortField := q.Get("sort"); sortField != "" {
		order := int32(1)
		if q.Get("order") == "desc" {
			order = -1
		}
		findOpts.SetSort(bson.D{{Key: sortField, Value: order}})
	}

	// Execute the query
	cursor, err := prodColl.Find(ctx, filter, findOpts)
	if err != nil {
		log.Println("Find products error:", err)
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	var products []models.Product
	if err := cursor.All(ctx, &products); err != nil {
		log.Println("Cursor.All error:", err)
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	// Populate brand and category for each product
	var resp []struct {
		models.Product `json:",inline"`
		Brand          models.Brand    `json:"brand"`
		Category       models.Category `json:"category"`
	}
	for _, p := range products {
		var b models.Brand
		if err := database.Collection("brands").
			FindOne(ctx, bson.M{"_id": p.Brand}).
			Decode(&b); err != nil {
			log.Println("Brand lookup error:", err)
		}
		var c models.Category
		if err := database.Collection("categories").
			FindOne(ctx, bson.M{"_id": p.Category}).
			Decode(&c); err != nil {
			log.Println("Category lookup error:", err)
		}
		resp = append(resp, struct {
			models.Product `json:",inline"`
			Brand          models.Brand    `json:"brand"`
			Category       models.Category `json:"category"`
		}{p, b, c})
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}
