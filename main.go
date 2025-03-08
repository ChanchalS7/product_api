package main

import (
	"log"
	"net/http"
	"os"

	"github.com/ChanchalS7/product_api/config"
	"github.com/ChanchalS7/product_api/middleware"
	"github.com/ChanchalS7/product_api/routes"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

func main() {
	config.LoadEnv()
	router := mux.NewRouter()
	//middleware
	router.Use(middleware.Logging)
	router.Use(middleware.RateLimit)
	//Routes
	routes.RegisterAuthRoutes(router)
	routes.RegisterProductRoutes(router)

	port := os.Getenv("PORT")
	logrus.Info("Server starting on port:", port)
	log.Fatal(http.ListenAndServe(":"+port, router))

}
