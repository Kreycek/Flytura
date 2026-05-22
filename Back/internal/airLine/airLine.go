package airLine

import (
	flytura "Flytura"
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
Alteração : 11/03/2026 17:56
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
			return nil, fmt.Errorf("Arquivo não encontrado")
		}
		return nil, fmt.Errorf("erro ao buscar plano de contas: %v", err)
	}

	// Converter o _id para string

	// Retornar o usuário como um mapa
	airLineReturn := map[string]any{
		"ID":                          airLineData.ID.Hex(), // Agora o campo ID é uma string
		"Name":                        airLineData.Name,
		"Code":                        airLineData.Code,
		"FileName":                    airLineData.FileName,
		"DtImportacao":                airLineData.CreatedAt,
		"Active":                      airLineData.Active,
		"ImportSheetOnlyVerifyNumber": airLineData.ImportSheetOnlyVerifyNumber,
		"QtdMinCharKey":               airLineData.QtdMinCharKey,
	}

	return airLineReturn, nil
}

/*
Função criada por Ricardo Silva Ferreira
Inicio da criação 26/11/2025 15:08
Data Final da criação : 26/11/2025 15:10
Verifica se o o arquivo importado realmente é o da companhia correta
*/
func ReturnFileDataAirline(client *mongo.Client, dbName, collectionName, code string, fileName string) (map[string]any, error) {

	collection := client.Database(dbName).Collection(collectionName)

	// objectID, erroId := primitive.ObjectIDFromHex(excelId)
	// if erroId != nil {
	// 	log.Fatalf("Erro ao converter string para ObjectID: %v", erroId)
	// }

	filter := bson.M{
		"code":     code,
		"fileName": fileName + "-" + code,
	}

	// Variável para armazenar o usuário retornado
	var airLineData models.AirLine

	// Usar FindOne para pegar apenas um único registro
	err := collection.FindOne(context.Background(), filter).Decode(&airLineData)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("O nome da planilha ou código estão inválidos")
		}
		return nil, fmt.Errorf("erro ao buscar plano de contas: %v", err)
	}

	// Converter o _id para string

	// Retornar o usuário como um mapa
	airLineReturn := map[string]any{
		"ID":                          airLineData.ID.Hex(), // Agora o campo ID é uma string
		"Name":                        airLineData.Name,
		"Code":                        airLineData.Code,
		"FileName":                    airLineData.FileName,
		"DtImportacao":                airLineData.CreatedAt,
		"Active":                      airLineData.Active,
		"ImportSheetOnlyVerifyNumber": airLineData.ImportSheetOnlyVerifyNumber,
	}

	return airLineReturn, nil
}

/*
	Função criada por Ricardo Silva Ferreira
	Inicio da criação 09/09/2025 22:34
	Data Final da criação : 09/09/2025 21:36
	Data modificação: 25/11/2025-> acrescentado verificação para active =true
*/
// Função para obter todos os diários para carregar o drop de buscar

func GetAirLines(client *mongo.Client, dbName, collectionName string) ([]any, error) {
	collection := db.GetCollection(client, dbName, collectionName)

	// Consultar apenas documentos com Active = true
	cursor, err := collection.Find(context.Background(), bson.M{"active": true})
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar companhias: %v", err)
	}
	defer cursor.Close(context.Background())

	var dadosBanco []any
	for cursor.Next(context.Background()) {
		var data models.AirLine
		if err := cursor.Decode(&data); err != nil {
			return nil, fmt.Errorf("erro ao decodificar companhias: %v", err)
		}

		dadosBanco = append(dadosBanco, map[string]any{
			"ID":           data.ID,
			"name":         data.Name,
			"code":         data.Code,
			"FileName":     data.FileName,
			"DtImportacao": data.CreatedAt,
			"Active":       data.Active,
		})
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("erro ao iterar no cursor: %v", err)
	}

	return dadosBanco, nil
}

/*
Função criada por Ricardo Silva Ferreira
Inicio da criação 22/05/2026 14:14
Data Final da criação :22/05/2026 14:14
*/
func SearchAirlineByName(items []any, name string) (string, string) {
	target := flytura.Normalize(name)

	for _, item := range items {

		// cast para map
		m, ok := item.(map[string]any)
		if !ok {
			continue
		}

		// extrair name
		nameVal, ok := m["name"].(string)
		if !ok {
			continue
		}

		// fmt.Println("target ", target, " ", flytura.Normalize(nameVal))
		// comparar normalizado
		if flytura.Normalize(nameVal) == target {

			// extrair code
			codeVal, _ := m["code"].(string)

			return codeVal, nameVal
		}
	}

	return "", ""
}
