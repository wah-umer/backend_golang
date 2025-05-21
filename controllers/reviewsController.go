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

var reviewColl = database.Collection("reviews")

func CreateReview(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Methods", "POST")

	var rev models.Review
	if err := json.NewDecoder(r.Body).Decode(&rev); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	rev.ID = primitive.NewObjectID()
	rev.CreatedAt = time.Now()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if _, err := reviewColl.InsertOne(ctx, rev); err != nil {
		log.Println("DB error on create review:", err)
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	var usr models.User
	if err := usersColl.
		FindOne(ctx, bson.M{"_id": rev.User}).
		Decode(&usr); err != nil {
		log.Println("User lookup error:", err)
	}

	resp := struct {
		models.Review `json:",inline"`
		User          models.User `json:"user"`
	}{rev, usr}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

func GetReviewByProductId(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Methods", "GET")

	prodID := mux.Vars(r)["id"]
	pid, err := primitive.ObjectIDFromHex(prodID)
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"product": pid}

	cur, err := reviewColl.Find(ctx, filter)
	if err != nil {
		log.Println("Find reviews error:", err)
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}
	defer cur.Close(ctx)

	var reviews []models.Review
	if err := cur.All(ctx, &reviews); err != nil {
		log.Println("Cursor.All error:", err)
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	var resp []struct {
		models.Review `json:",inline"`
		User          models.User `json:"user"`
	}
	for _, rv := range reviews {
		var usr models.User
		if err := usersColl.
			FindOne(ctx, bson.M{"_id": rv.User}).
			Decode(&usr); err != nil {
			log.Println("User lookup error:", err)
		}
		resp = append(resp, struct {
			models.Review `json:",inline"`
			User          models.User `json:"user"`
		}{rv, usr})
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

func UpdateReviewById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Methods", "PATCH")

	revID := mux.Vars(r)["id"]
	rid, err := primitive.ObjectIDFromHex(revID)
	if err != nil {
		http.Error(w, "Invalid review ID", http.StatusBadRequest)
		return
	}

	var updates map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	delete(updates, "_id")
	delete(updates, "user")
	delete(updates, "product")
	delete(updates, "createdAt")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	update := bson.M{"$set": updates}

	var updated models.Review
	err = reviewColl.
		FindOneAndUpdate(ctx, bson.M{"_id": rid}, update, opts).
		Decode(&updated)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			http.Error(w, "Review not found", http.StatusNotFound)
		} else {
			log.Println("DB error on update:", err)
			http.Error(w, "Server error", http.StatusInternalServerError)
		}
		return
	}

	var usr models.User
	if err := usersColl.
		FindOne(ctx, bson.M{"_id": updated.User}).
		Decode(&usr); err != nil {
		log.Println("User lookup error:", err)
	}

	resp := struct {
		models.Review `json:",inline"`
		User          models.User `json:"user"`
	}{updated, usr}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

func DeleteReviewById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Methods", "DELETE")

	revID := mux.Vars(r)["id"]
	rid, err := primitive.ObjectIDFromHex(revID)
	if err != nil {
		http.Error(w, "Invalid review ID", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var deleted models.Review
	err = reviewColl.
		FindOneAndDelete(ctx, bson.M{"_id": rid}).
		Decode(&deleted)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			http.Error(w, "Review not found", http.StatusNotFound)
		} else {
			log.Println("DB error on delete:", err)
			http.Error(w, "Server error", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(deleted)
}
