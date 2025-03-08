package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ChanchalS7/product_api/controllers"
	"github.com/ChanchalS7/product_api/models"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestRegisterHandler(t *testing.T) {
	router := mux.NewRouter()
	router.HandleFunc("/register", controllers.Register).Methods("POST")

	t.Run("Valid Registration", func(t *testing.T) {
		user := models.User{
			Email:    "test@example.com",
			Password: "test@123",
		}
		body, _ := json.Marshal(user)
		req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer((body)))
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusCreated, rr.Code)
	})

	t.Run("Duplicate Registration", func(t *testing.T) {
		user := models.User{
			Email:    "test@example.com",
			Password: "test@123",
		}
		body, _ := json.Marshal(user)
		req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer((body)))
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusCreated, rr.Code)
	})

}

func TestLoginHandler(t *testing.T) {
	router := mux.NewRouter()
	router.HandleFunc("/login", controllers.Login).Methods("POST")

	t.Run("Valid Login", func(t *testing.T) {
		creds := map[string]string{
			"email":    "test@example.com",
			"password": "test@123",
		}
		body, _ := json.Marshal(creds)
		req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
	})
}
