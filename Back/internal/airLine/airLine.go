package airLine

import (
	"Flytura/internal/db"
	"Flytura/internal/models"
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

/*
Função criada por Ricardo Silva Ferreira
Inicio da criação 09/09/2025 21:37
Data Final da criação : 09/09/2025 21:53
*/
func GetAirLineFileName(client *mongo.Client, dbName, collectionName, code string) (map[string]any, error) {

	collection := client.Database(dbName).Collection(collectionName)

	// objectID, erroId := primitive.ObjectIDFromHex(excelId)
	// if erroId != nil {
	// 	log.Fatalf("Erro ao converter string para ObjectID: %v", erroId)
	// }

	filter := bson.M{"code": code}

	// Variável para armazenar o usuário retornado
	var airLineData models.AirLine

	// Usar FindOne para pegar apenas um único registro
	err := collection.FindOne(context.Background(), filter).Decode(&airLineData)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("plano de contas não encontrado")
		}
		return nil, fmt.Errorf("erro ao buscar plano de contas: %v", err)
	}

	// Converter o _id para string

	// Retornar o usuário como um mapa
	airLineReturn := map[string]any{
		"ID":           airLineData.ID.Hex(), // Agora o campo ID é uma string
		"Name":         airLineData.Name,
		"Code":         airLineData.Code,
		"FileName":     airLineData.FileName,
		"DtImportacao": airLineData.CreatedAt,
		"Active":       airLineData.Active,
	}

	return airLineReturn, nil
}

/*
Função criada por Ricardo Silva Ferreira
Inicio da criação 09/09/2025 22:34
Data Final da criação : 09/09/2025 21:36
*/
// Função para obter todos os diários para carregar o drop de buscar
func GetAirLines(client *mongo.Client, dbName, collectionName string) ([]any, error) {
	collection := db.GetCollection(client, dbName, collectionName)

	// Consultar todos os documentos
	cursor, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar usuários: %v", err)
	}
	defer cursor.Close(context.Background())

	var dadosBanco []any
	for cursor.Next(context.Background()) {
		var data models.AirLine
		if err := cursor.Decode(&data); err != nil {
			return nil, fmt.Errorf("erro ao decodificar ,comapnias: %v", err)
		}

		// Converter o _id do MongoDB para string para retorno
		Id := data.ID
		// Preenche o usuário com o ID convertido em string
		dadosBanco = append(dadosBanco, map[string]any{
			"ID":           Id, // Agora o campo ID é uma string
			"name":         data.Name,
			"code":         data.Code,
			"FileName":     data.FileName,
			"DtImportacao": data.CreatedAt,
			"Active":       data.Active,
		})

	}

	// Verifica se houve algum erro durante a iteração do cursor
	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("erro ao iterar no cursor: %v", err)
	}

	// Retorna os usuários
	return dadosBanco, nil
}
