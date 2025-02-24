package models

import (
	"context"
	"time"

	"github.com/ChanchalS7/product_api/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Product struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	Name        string             `bson:"name"`
	Description string             `bson:"description"`
	Price       string             `bson:"price"`
	CreatedAt   time.Time          `bson:"created_at"`
}

func CreateProduct(product *Product) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	product.CreatedAt = time.Now()
	_, err := config.DB.Collection("products").InsertOne(ctx, product)
	return err
}

func GetProduct(id string) (*Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	var product Product
	err = config.DB.Collection("products").FindOne(ctx, bson.M{"_id": objID}).Decode(&product)
	return &product, err
}

func GetAllProducts() ([]Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var products []Product

	cursor, err := config.DB.Collection("products").Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var product Product
		cursor.Decode(&product)
		products = append(products, product)
	}
	return products, nil
}

func UpdateProduct(id string, updateData bson.M) error {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = config.DB.Collection("products").UpdateOne(
		ctx,
		bson.M{"_id": objID},
		bson.M{"$set": updateData},
	)
	return err
}

func DeleteProduct(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = config.DB.Collection("products").DeleteOne(ctx, bson.M{"_id": objID})
	return err
}
