package main

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)


type User struct {
	Name string `bson:"name"`
	Age  int    `bson:"age"`
}

func connectMongo() (*mongo.Client, error) {
	uri := "mongodb://localhost:27017"

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	if err = client.Ping(ctx, nil); err != nil {
		return nil, err
	}
	return client, nil
}

func insertUser(ctx context.Context, c *mongo.Collection, name string, age int) (*mongo.InsertOneResult, error) {
	return c.InsertOne(ctx, bson.M{"name": name, "age": age})
}

func userAgeMongo(ctx context.Context, c *mongo.Collection, name string) (int, error) {
	var user User
	err := c.FindOne(ctx, bson.M{"name": name}).Decode(&user)
	if err != nil {
		return 0, err
	}
	return user.Age, nil
}

func countMongo(ctx context.Context, c *mongo.Collection) (int64, error) {
	return c.CountDocuments(ctx, bson.M{})
}

// setAgeMongo обновляет возраст пользователя через $set
func setAgeMongo(ctx context.Context, c *mongo.Collection, name string, age int) error {
	_, err := c.UpdateOne(ctx, bson.M{"name": name}, bson.M{"$set": bson.M{"age": age}})
	return err
}

func main() {
	ctx := context.Background()

	client, err := connectMongo()
	if err != nil {
		panic(err)
	}
	defer client.Disconnect(ctx)

	c := client.Database("shop").Collection("users")
	c.Drop(ctx) // Очищаем коллекцию перед тестом

	_, e1 := insertUser(ctx, c, "Ann", 30)
	age, e2 := userAgeMongo(ctx, c, "Ann")
	n, e3 := countMongo(ctx, c)
	e4 := setAgeMongo(ctx, c, "Ann", 31)
	age2, _ := userAgeMongo(ctx, c, "Ann")

	result := e1 == nil && e2 == nil && age == 30 &&
		e3 == nil && n == 1 && e4 == nil && age2 == 31

	fmt.Println(result)
}
