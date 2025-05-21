package controllers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/umerwaheed/backend_golang/database"
	"github.com/umerwaheed/backend_golang/models"
	"go.mongodb.org/mongo-driver/bson"
)

var brandColl = database.Collection("brands")

func GetAllBrand(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Methods", "GET")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cur, err := brandColl.Find(ctx, bson.M{})
	if err != nil {
		log.Println("DB error on Find brands:", err)
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}
	defer cur.Close(ctx)

	var brands []models.Brand
	cur.All(ctx, &brands)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(brands)
}
