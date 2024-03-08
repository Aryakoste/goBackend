package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const uri = "mongodb+srv://aryakhochare:R5jnOIJ4zSbrG9hz@cluster0.eerc5uv.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0"

type App struct {
	Client *mongo.Client
}

type Store struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	Location string `json:"location"`
}

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

	router.Run("localhost:8080")
}
