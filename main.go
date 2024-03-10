package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const uri = "mongodb+srv://aryakhochare:R5jnOIJ4zSbrG9hz@cluster0.eerc5uv.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0"

type App struct {
	Client *mongo.Client
}

type Store struct {
	ID          primitive.ObjectID   `bson:"_id,omitempty"`
	Name        string               `bson:"name"`
	Type        string               `bson:"type"`
	Location    string               `bson:"location"`
	ProducstIDs []primitive.ObjectID `bson:"productIds"`
}

type Product struct {
	ID    primitive.ObjectID `bson:"_id,omitempty"`
	Name  string             `bson:"name"`
	Price string             `bson:"price"`
}

// func (app *App) getProductsFromStore(c *gin.Context, storeID string) {
// 	var results []Product
// }

func (app *App) postStore(c *gin.Context) {
	var newStore Store

	if err := c.ShouldBindJSON(&newStore); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	collection := app.Client.Database("test").Collection("stores")
	_, err := collection.InsertOne(context.Background(), newStore)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Store Added Successfully!", "store": newStore})

}

func (app *App) getAllStores(c *gin.Context) {
	var stores []Store
	collection := app.Client.Database("test").Collection("stores")

	cursor, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching Stores"})
		return
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var store Store

		if err := cursor.Decode(&store); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error Decoding Store"})
		}

		stores = append(stores, store)
	}
	if err := cursor.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Cursor Error"})
	}

	c.JSON(http.StatusOK, stores)
}

func main() {

	router := gin.Default()

	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI)

	client, err := mongo.Connect(context.TODO(), opts)

	if err != nil {
		panic(err)
	}
	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	var result bson.M
	if err := client.Database("admin").RunCommand(context.TODO(), bson.D{{"ping", 1}}).Decode(&result); err != nil {
		panic(err)
	}
	fmt.Println("Pinged your deployment. You successfully connected to MongoDB!")

	app := App{
		Client: client,
	}

	router.POST("addStore", app.postStore)
	router.GET("getAllStores", app.getAllStores)

	router.Run("localhost:8080")
}
