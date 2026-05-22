package outPutInvoices

import (
	flytura "Flytura"
	"Flytura/internal/db"
	"Flytura/internal/models"
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

/*
Função criada por Ricardo Silva Ferreira
Inicio da criação 09/04/2026 23:34
Data Final da criação : 09/04/2026 23:39
*/
func SearchOutPutInvoicesInforme(
	client *mongo.Client,
	dbName, collectionName string,
	key *string,
	companyCode *string,
	startDate *time.Time,
	endDate *time.Time) ([]any, int64, error) {

	collection := client.Database(dbName).Collection(collectionName)

	// Criando o filtro dinâmico
	filter := bson.M{}

	if key != nil && *key != "" {
		filter["key"] = bson.M{"$regex": *key, "$options": "i"}
	}

	if companyCode != nil && *companyCode != "" {
		filter["companyCode"] = *companyCode
	}

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
		options.Find().
			SetSort(bson.D{{Key: "dtProcess", Value: -1}}),
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
			"Ticket":               data.Ticket,
			"CreatedAt":            data.CreatedAt,
		})
	}

	// Retorna usuários e total de registros
	return excelData, total, nil
}

/*
Função criada por Ricardo Silva Ferreira
Inicio da criação 11/11/2025 11:00
Data Final da criação : 11/11/2025 11:13
OBS: Adicionada ordenação descrecente em: 09/04/2026 23:34
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
		options.Find().
			SetSkip(int64((page-1)*limit)).
			SetLimit(int64(limit)).
			SetSort(bson.D{{Key: "dtProcess", Value: -1}}),
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
			"Segment":              data.Segment,
			"Ticket":               data.Ticket,
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

		nowUTC := time.Now().UTC()

		dh, err := flytura.DiffHours(nowUTC, flytura.Fuso1, flytura.Fuso2)
		if err != nil {
			panic(err)
		}
		v.CreatedAt = nowUTC.Add(-time.Duration(dh) * time.Hour) // Atualiza a data de criação
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

/*
Função criada por Ricardo Silva Ferreira
Inicio da criação 20/11/2025 17:17
Data Final da criação : 20/11/2025 17:17
*/
func DeleteOutPutInvoicesByID(client *mongo.Client, dbName, collectionName, id string) error {
	collection := client.Database(dbName).Collection(collectionName)

	// Criar um contexto para a operação
	ctx := context.Background()

	// Converter o ID para ObjectID
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("ID inválido: %v", err)
	}

	// Executar a exclusão
	result, err := collection.DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		return fmt.Errorf("erro ao deletar registro: %v", err)
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("nenhum documento encontrado com o ID informado")
	}

	return nil
}

/*
Função criada por Ricardo Silva Ferreira
Inicio da criação 08/02/2026 01:43
Data Final da criação :   08/02/2026 01:43
*/
/*
Função criada por Ricardo Silva Ferreira
Inicio da criação 08/02/2026 01:43
Data Final da criação :   08/02/2026 01:43
*/
func GroupByCompanyCodeSumSection(
	client *mongo.Client,
	dbName, collectionName string,
	startDateStr, endDateStr string,
	companyCode string,
) ([]bson.M, error) {

	collection := db.GetCollection(client, dbName, collectionName)

	// Filtro dinâmico ($match)
	matchConditions := bson.D{}
	dateFilter := bson.D{}

	// START DATE (RFC3339)
	if startDateStr != "" {
		if start, err := time.Parse(time.RFC3339, startDateStr); err == nil {
			startUTC := time.Date(start.UTC().Year(), start.UTC().Month(), start.UTC().Day(), 0, 0, 0, 0, time.UTC)
			dateFilter = append(dateFilter, bson.E{Key: "$gte", Value: startUTC})
		}
	}

	// END DATE (RFC3339)
	if endDateStr != "" {
		if end, err := time.Parse(time.RFC3339, endDateStr); err == nil {
			endUTC := time.Date(end.UTC().Year(), end.UTC().Month(), end.UTC().Day(), 0, 0, 0, 0, time.UTC).Add(24 * time.Hour)
			dateFilter = append(dateFilter, bson.E{Key: "$lt", Value: endUTC})
		}
	}

	if len(dateFilter) > 0 {
		matchConditions = append(matchConditions, bson.E{
			Key:   "dtProcess",
			Value: dateFilter,
		})
	}

	// Filtro por companyCode (se fornecido)
	if companyCode != "" {
		matchConditions = append(matchConditions, bson.E{
			Key: "companyCode", Value: companyCode,
		})
	}

	pipeline := mongo.Pipeline{}

	// Aplica $match se houver filtros acumulados
	if len(matchConditions) > 0 {
		pipeline = append(pipeline, bson.D{{Key: "$match", Value: matchConditions}})
	}

	// Garante que existe companyCode (quando não foi passado e vamos agrupar por ele)
	pipeline = append(pipeline, bson.D{{Key: "$match", Value: bson.D{
		{Key: "companyCode", Value: bson.D{{Key: "$ne", Value: nil}}},
	}}})

	// Converte 'section' para número com segurança (pode vir string/int/double)
	pipeline = append(pipeline, bson.D{{Key: "$addFields", Value: bson.D{
		{Key: "sectionAsNumber", Value: bson.D{
			{Key: "$convert", Value: bson.D{
				{Key: "input", Value: "$section"},
				{Key: "to", Value: "double"},
				{Key: "onError", Value: 0},
				{Key: "onNull", Value: 0},
			}},
		}},
	}}})

	// Agrupa por companyCode somando 'section'
	pipeline = append(pipeline, bson.D{{Key: "$group", Value: bson.D{
		{Key: "_id", Value: "$companyCode"},
		{Key: "total", Value: bson.D{{Key: "$sum", Value: "$sectionAsNumber"}}},
		{Key: "companyName", Value: bson.D{{Key: "$first", Value: "$companyName"}}},
	}}})

	// Projeta no formato desejado (companyName -> company)
	pipeline = append(pipeline, bson.D{{Key: "$project", Value: bson.D{
		{Key: "_id", Value: 0},
		{Key: "companyCode", Value: "$_id"},
		{Key: "company", Value: "$companyName"},
		{Key: "total", Value: 1},
	}}})

	// Ordena por total desc
	pipeline = append(pipeline, bson.D{{Key: "$sort", Value: bson.D{
		{Key: "total", Value: -1},
	}}})

	cursor, err := collection.Aggregate(context.Background(), pipeline)
	if err != nil {
		return nil, fmt.Errorf("erro ao agregar soma de section por companyCode: %v", err)
	}
	defer cursor.Close(context.Background())

	var resultados []bson.M
	if err := cursor.All(context.Background(), &resultados); err != nil {
		return nil, fmt.Errorf("erro ao decodificar resultados: %v", err)
	}

	return resultados, nil
}
