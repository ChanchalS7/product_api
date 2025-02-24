package controllers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/ChanchalS7/product_api/config"
	"github.com/ChanchalS7/product_api/models"
	"github.com/dgrijalva/jwt-go"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

func Register(w http.ResponseWriter, r *http.Request) {
	var user models.User

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		logrus.Error("Invalid request payload", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	existingUser, _ := models.FindUserByEmail(user.Email)
	if existingUser != nil {
		http.Error(w, "User already exists", http.StatusConflict)
		return
	}
	err = models.CreateUser(&user)
	if err != nil {
		logrus.Error("Error creating user:", err)
		http.Error(w, "Error creating user", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "User created successfully"})
}

func Login(w http.ResponseWriter, r *http.Request) {
	var credentials struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		logrus.Error("Invalid request payload:", err)
		logrus.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	user, err := models.FindUserByEmail(credentials.Email)
	if err != nil {
		logrus.Error("User not found", err)
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credentials.Password))
	if err != nil {
		logrus.Error("Invalid password:", err)
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:    user.ID.Hex(),
		ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
	})
	token, err := claims.SignedString([]byte(config.JWTSecret))
	if err != nil {
		logrus.Error("Error generating token:", err)
		logrus.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "jwt",
		Value:    token,
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	})
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

func Logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "jwt",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HttpOnly: true,
	})
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Logged out successfully"})
}
