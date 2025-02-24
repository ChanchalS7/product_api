package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/ChanchalS7/product_api/models"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
)

func CreateProduct(w http.ResponseWriter, r *http.Request) {
	var product models.Product

	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		logrus.Error("Invalid request payload", err)
		http.Error(w, "Invalid request payload:", http.StatusBadRequest)
		return
	}
	err = models.CreateProduct(&product)
	if err != nil {
		logrus.Error("Error creating product:", err)
		http.Error(w, "Error creating product", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(product)
}

func GetProduct(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	product, err := models.GetProduct(params["id"])
	if err != nil {
		logrus.Warn("Product not found:", err)
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(product)
}

func GetAllProducts(w http.ResponseWriter, r *http.Request) {
	products, err := models.GetAllProducts()
	if err != nil {
		logrus.Error("Error fetching products:", err)
		http.Error(w, "Error fetching products", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(products)
}

func UpdateProduct(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var updateData bson.M
	err := json.NewDecoder(r.Body).Decode(&updateData)
	if err != nil {
		logrus.Error("Invalid request payload:", err)
		http.Error(w, "Invalid request payload:", http.StatusBadRequest)
		return
	}
	err = models.UpdateProduct(params["id"], updateData)
	if err != nil {
		logrus.Error("Error updating product:", err)
		http.Error(w, "Error updating product", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Product udpated successfully"})
}

func DeleteProduct(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	err := models.DeleteProduct(params["id"])
	if err != nil {
		logrus.Error(w, "Error deleting product", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Product deleted successfully"})
}
