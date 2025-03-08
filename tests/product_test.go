package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ChanchalS7/product_api/controllers"
	"github.com/ChanchalS7/product_api/middleware"
	"github.com/ChanchalS7/product_api/models"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

var authToken string

func setup() *mux.Router {
	router := mux.NewRouter()
	router.Use(middleware.JWTAuth)

	// Product routes
	router.HandleFunc("/products", controllers.CreateProduct).Methods("POST")
	router.HandleFunc("/products", controllers.GetAllProducts).Methods("GET")
	router.HandleFunc("/products/{id}", controllers.GetProduct).Methods("GET")
	router.HandleFunc("/products/{id}", controllers.UpdateProduct).Methods("PUT")
	router.HandleFunc("/products/{id}", controllers.DeleteProduct).Methods("DELETE")

	// Get auth token
	creds := map[string]string{
		"email":    "test@example.com",
		"password": "password123",
	}
	body, _ := json.Marshal(creds)
	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(body))
	rr := httptest.NewRecorder()
	controllers.Login(rr, req)

	var response map[string]string
	json.Unmarshal(rr.Body.Bytes(), &response)
	authToken = response["token"]

	return router
}

func TestProductCRUD(t *testing.T) {
	router := setup()
	var productID string

	t.Run("Create Product", func(t *testing.T) {
		product := models.Product{
			Name:        "Test Product",
			Description: "Test Description",
			Price:       99.99,
		}
		body, _ := json.Marshal(product)
		req, _ := http.NewRequest("POST", "/products", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+authToken)

		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusCreated, rr.Code)

		var createdProduct models.Product
		json.Unmarshal(rr.Body.Bytes(), &createdProduct)
		productID = createdProduct.ID.Hex()
	})

	t.Run("Get All Products", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/products", nil)
		req.Header.Set("Authorization", "Bearer "+authToken)

		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.NotEmpty(t, rr.Body.String())
	})

	t.Run("Update Product", func(t *testing.T) {
		update := map[string]interface{}{
			"price": 129.99,
		}
		body, _ := json.Marshal(update)
		req, _ := http.NewRequest("PUT", "/products/"+productID, bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+authToken)

		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
	})

	t.Run("Delete Product", func(t *testing.T) {
		req, _ := http.NewRequest("DELETE", "/products/"+productID, nil)
		req.Header.Set("Authorization", "Bearer "+authToken)

		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
	})
}

func TestProtectedRoutes(t *testing.T) {
	router := setup()

	t.Run("Unauthorized Access", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/products", nil)
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("Invalid Token", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/products", nil)
		req.Header.Set("Authorization", "Bearer invalidtoken")
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusUnauthorized, rr.Code)
	})
}
