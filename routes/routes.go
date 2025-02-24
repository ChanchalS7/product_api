package routes

import (
	"github.com/ChanchalS7/product_api/controllers"
	"github.com/ChanchalS7/product_api/middleware"
	"github.com/gorilla/mux"
)

func RegisterAuthRoutes(router *mux.Router) {
	router.HandleFunc("/register", controllers.Register).Methods("POST")
	router.HandleFunc("/login", controllers.Login).Methods("POST")
	router.HandleFunc("/logout", controllers.Logout).Methods("POST")

}

func RegisterProductRoutes(router *mux.Router) {
	productRouter := router.PathPrefix("/products").Subrouter()
	productRouter.Use(middleware.JWTAuth)

	productRouter.HandleFunc("", controllers.CreateProduct).Methods("POST")

	productRouter.HandleFunc("", controllers.GetAllProducts).Methods("GET")

	productRouter.HandleFunc("/{id}", controllers.GetProduct).Methods("GET")

	productRouter.HandleFunc("/{id}", controllers.UpdateProduct).Methods("PUT")

	productRouter.HandleFunc("/{id}", controllers.DeleteProduct).Methods("DELETE")

}
