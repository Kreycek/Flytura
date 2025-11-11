package outPutInvoices

import (
	"Flytura/internal/models"
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

/*
Função criada por Ricardo Silva Ferreira
Inicio da criação 11/11/2025 11:00
Data Final da criação : 11/11/2025 11:13
*/
func SearchOutPutInvoicesPagination(
	client *mongo.Client,
	dbName, collectionName string,
	key *string,
	companyCode *string,
	startDate *time.Time,
	endDate *time.Time,
	page,
	limit int) ([]any, int64, error) {

	collection := client.Database(dbName).Collection(collectionName)

	// Criando o filtro dinâmico
	filter := bson.M{}

	if key != nil && *key != "" {
		filter["key"] = bson.M{"$regex": *key, "$options": "i"}
	}

	if companyCode != nil && *companyCode != "" {
		filter["companyCode"] = *companyCode
	}
	fmt.Println("page ", page)
	fmt.Println("key ", *key)
	fmt.Println("limit ", limit)
	fmt.Println("startDate ", startDate)
	fmt.Println("companyCode ", *companyCode)
	fmt.Println("endDate ", endDate)

	// fmt.Println("startDate", startDate)
	// fmt.Println("endDate", endDate)
	if startDate != nil || endDate != nil {
		dateFilter := bson.M{}

		if startDate != nil {
			// Zera a hora de startDate (00:00:00)
			start := time.Date(startDate.Year(), startDate.Month(), startDate.Day(), 0, 0, 0, 0, startDate.Location())
			dateFilter["$gte"] = start
		}

		if endDate != nil {
			// Ajusta endDate para o final do dia (23:59:59.999999999)
			end := time.Date(endDate.Year(), endDate.Month(), endDate.Day(), 23, 59, 59, int(time.Second-time.Nanosecond), endDate.Location())
			dateFilter["$lte"] = end
		}

		filter["createdAt"] = dateFilter
	}

	// Contar total de usuários antes da paginação
	total, err := collection.CountDocuments(context.Background(), filter)
	if err != nil {
		return nil, 0, err
	}

	// Executa a consulta com paginação
	cursor, err := collection.Find(
		context.Background(),
		filter,
		options.Find().SetSkip(int64((page-1)*limit)).SetLimit(int64(limit)),
	)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(context.Background())

	// Processa os resultados
	var excelData []any
	for cursor.Next(context.Background()) {
		var data models.OutPutInvoices
		if err := cursor.Decode(&data); err != nil {
			return nil, 0, fmt.Errorf("erro ao decodificar OutPutInvoices: %v", err)
		}

		excelData = append(excelData, map[string]any{
			"ID":                   data.ID.Hex(), // Convertendo para string
			"Key":                  data.Key,
			"Status":               data.Status,
			"DtProcess":            data.DtProcess,
			"MonthProcess":         data.MonthProcess,
			"DtFly":                data.DtFLy,
			"RFC":                  data.RFC,
			"IVAValue":             data.IVAValue,
			"Tax":                  data.Tax,
			"Rate":                 data.Rate,
			"FactorType":           data.FactorType,
			"TransferredBaseValue": data.TransferredBaseValue,
			"SubTotalValue":        data.SubTotalValue,
			"TotalValue":           data.TotalValue,
			"TUA":                  data.TUA,
			"Ruta":                 data.Ruta,
			"OtherValues":          data.OtherValues,
			"CompanyName":          data.CompanyName,
			"CreatedAt":            data.CreatedAt,
		})
	}

	// Retorna usuários e total de registros
	return excelData, total, nil
}

/*
	Função criada por Ricardo Silva Ferreira
	Inicio da criação 11/11/2025 11:32
	Data Final da criação : 11/11/2025 11:39
*/

func InsertOutPutInvoices(client *mongo.Client, dbName, collectionName string, data []models.OutPutInvoices) error {
	collection := client.Database(dbName).Collection(collectionName)

	if len(data) == 0 {
		return fmt.Errorf("nenhum dado fornecido para inserção")
	}

	// Criar um contexto para a operação de inserção
	ctx := context.Background()

	// airLines, erroAirline := airLine.GetAirLines(client, dbName)

	// if erroAirline != nil {
	// 	// http.Error(w, fmt.Sprintf("Token inválido: %v", err), http.StatusUnauthorized)
	// 	return fmt.Errorf("Erro ao carregar companias: %v", erroAirline)
	// }

	//Abaixo exemplo de  Bulk Insert
	docs := make([]interface{}, len(data))
	for i, v := range data {

		v.CreatedAt = time.Now() // Atualiza a data de criação
		v.Active = true

		docs[i] = v

	}

	// Inserir o documento
	_, err := collection.InsertMany(ctx, docs)
	if err != nil {
		return fmt.Errorf("erro ao inserir dados do excel: %v", err)
	}

	return nil
}
