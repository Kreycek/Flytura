package db

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectMongoDB(uri string) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clienteOptions := options.Client().ApplyURI(uri)

	client, err := mongo.Connect(ctx, clienteOptions)

	if err != nil {
		return nil, fmt.Errorf("erro ao conectar com o MongoDB %v", err)
	}

	return client, nil
}

// Função para obter uma coleção do MongoDB
func GetCollection(client *mongo.Client, dbName, collectionName string) *mongo.Collection {
	return client.Database(dbName).Collection(collectionName)
}

// Função para fechar a conexão com o MongoDB
func CloseMongoDB(client *mongo.Client) error {
	return client.Disconnect(context.Background())
}
