package perfil

import (
	"Flytura/internal/db"
	"Flytura/internal/models"
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// Função para obter todos os usuários do banco de dados
func GetAllPerfil(client *mongo.Client, dbName, collectionName string) ([]any, error) {
	collection := db.GetCollection(client, dbName, collectionName)

	// Consultar todos os documentos
	cursor, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar usuários: %v", err)
	}
	defer cursor.Close(context.Background())

	var dadosBanco []any
	for cursor.Next(context.Background()) {
		var perfil models.Perfil
		if err := cursor.Decode(&perfil); err != nil {
			return nil, fmt.Errorf("erro ao decodificar usuário: %v", err)
		}

		// Converter o _id do MongoDB para string para retorno
		perfilID := perfil.ID
		// Preenche o usuário com o ID convertido em string
		dadosBanco = append(dadosBanco, map[string]any{
			"ID":        perfilID, // Agora o campo ID é uma string
			"name":      perfil.Name,
			"shortName": perfil.ShortName,
		})

	}

	// Verifica se houve algum erro durante a iteração do cursor
	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("erro ao iterar no cursor: %v", err)
	}

	// Retorna os usuários
	return dadosBanco, nil
}
