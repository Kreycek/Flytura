package logs

import (
	"Flytura/internal/models"
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
)

/*
	Função criada por Ricardo Silva Ferreira
	Inicio da criação 03/12/2025 15:59
	Data Final da criação : 03/12/2025 16:00
*/

func InsertLog(client *mongo.Client, dbName, collectionName string, data models.Log) error {
	collection := client.Database(dbName).Collection(collectionName)

	// Criar um contexto para a operação de inserção
	ctx := context.Background()

	// Inserir o documento
	_, err := collection.InsertOne(ctx, data)
	if err != nil {
		return fmt.Errorf("erro ao inserir dados do excel: %v", err)
	}

	return nil
}
