package purcharseRecord

import (
	flytura "Flytura"
	"Flytura/internal/db"
	"Flytura/internal/models"
	"context"
	"fmt"
	"log"
	"path/filepath"
	"regexp"
	"strings"

	"time"

	"github.com/extrame/xls"
	"github.com/xuri/excelize/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

/*
Função criada por Ricardo Silva Ferreira
Inicio da criação 03/09/2025 22:20
Data Final da criação : 04/09/2025 18:50
*/
func ProcessPurcharseRecordExcel(filePath, fileName, companyName, companyCode string, client *mongo.Client, dbName, collectionName string) (int64, int64, bool, error) {

	extensao := strings.ToLower(filepath.Ext(filePath))
	var sheetList []string
	var rows [][]string
	var emptyRecord int64 = 0
	var totalRecord int64 = 0
	var emptySheet = true

	// fmt.Println("Extensão ", extensao)

	if extensao == ".xls" {
		excelXls, err := xls.Open(filePath, "utf-8")

		if err != nil {
			return 0, 0, emptySheet, err
		}

		// fmt.Println("XLS ")

		sheet := excelXls.GetSheet(0)
		if sheet == nil {
			return 0, 0, emptySheet, fmt.Errorf("nenhuma aba encontrada no .xls")
		}

		for i := 0; i <= int(sheet.MaxRow); i++ {
			row := sheet.Row(i)
			if row == nil {
				continue
			}
			linha := []string{row.Col(0), row.Col(1), row.Col(2)}
			rows = append(rows, linha)
		}

	} else if extensao == ".xlsx" {
		// fmt.Println("XLSX")
		excelXlsx, err := excelize.OpenFile(filePath)
		if err != nil {
			return 0, 0, emptySheet, err
		}
		sheetList = excelXlsx.GetSheetList()
		if len(sheetList) == 0 {
			log.Fatal("Nenhuma aba encontrada")
		}

		sheetName := sheetList[0]
		// fmt.Println("sheetName ", sheetName)
		rows, err = excelXlsx.GetRows(sheetName)
		if err != nil {
			return 0, 0, emptySheet, fmt.Errorf("erro ao ler linhas da aba %s: %v", sheetName, err)
		}

	}

	collection := client.Database(dbName).Collection(collectionName)

	// Criar um contexto para a operação de inserção
	ctx := context.Background()

	for i, row := range rows {

		if i == 0 {
			continue // cabeçalho
		}

		obj := models.PurcharseRecord{

			Key:       "",
			Name:      "",
			LastName:  "",
			FileName:  "",
			Status:    "Fila",
			Active:    true,
			CreatedAt: time.Now(),
		}

		//Retira os espaços das string
		re := regexp.MustCompile(`\s+`)
		// fmt.Println(row[0])
		obj.Key = re.ReplaceAllString(row[0], "")
		obj.Name = row[1]
		obj.LastName = row[2]
		obj.FileName = fileName
		obj.CompanyCode = companyCode
		obj.CompanyName = companyName

		// preco, _ := strconv.ParseFloat(row[1], 64)
		// produto := bson.M{
		// 	"nome":  row[0],
		// 	"preco": preco,
		// }

		if obj.Key != "" {

			emptySheet = false
			//A FUNÇÃO ABAIXO VERIFICA SE A CHAVE KEY JÁ FOI IMPORTADA
			exist, errVerify := VeryExistKey(client, flytura.DBName, flytura.PurcharseRecordTableName, obj.Key)
			if errVerify != nil {
				log.Println("Erro ao inserir:", errVerify)
			} else {

				if !exist {
					totalRecord++
					_, err := collection.InsertOne(ctx, obj)
					if err != nil {
						log.Println("Erro ao inserir:", err)
					}
				}
			}
		} else {
			//SE NA PLANILHA SE A CHAVE KEY ESTIVER VAZIA CONTA
			emptyRecord++
		}

	}

	return emptyRecord, totalRecord, emptySheet, nil
}

/*
Função criada por Ricardo Silva Ferreira
Inicio da criação 04/09/2025 21:20
Data Final da criação : 04/09/2025 21:31
*/
// Função para obter todos os diários para carregar o drop de buscar
func GetPurcharseRecord(client *mongo.Client, dbName, collectionName string) ([]any, error) {
	collection := db.GetCollection(client, dbName, collectionName)

	// Consultar todos os documentos
	cursor, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar usuários: %v", err)
	}
	defer cursor.Close(context.Background())

	var dadosBanco []any
	for cursor.Next(context.Background()) {
		var cc models.PurcharseRecord
		if err := cursor.Decode(&cc); err != nil {
			return nil, fmt.Errorf("erro ao decodificar ,centro de custo: %v", err)
		}

		// Converter o _id do MongoDB para string para retorno
		Id := cc.ID
		// Preenche o usuário com o ID convertido em string
		dadosBanco = append(dadosBanco, map[string]any{
			"ID":           Id, // Agora o campo ID é uma string
			"key":          cc.Key,
			"name":         cc.Name,
			"lastName":     cc.LastName,
			"FileName":     cc.FileName,
			"CompanyCode":  cc.CompanyCode,
			"Status":       cc.Status,
			"CompanyName":  cc.CompanyName,
			"DtImportacao": cc.CreatedAt,
			"Active":       cc.Active,
		})

	}

	// Verifica se houve algum erro durante a iteração do cursor
	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("erro ao iterar no cursor: %v", err)
	}

	// Retorna os usuários
	return dadosBanco, nil
}

/*
Função criada por Ricardo Silva Ferreira
Inicio da criação 04/09/2025 21:31
Data Final da criação : 04/09/2025 21:35
*/
func SearchPurcharseRecordPagination(
	client *mongo.Client,
	dbName, collectionName string,
	key *string,
	name *string,
	lastName *string,
	companyCode *string,
	startDate *time.Time,
	endDate *time.Time,
	status *string,
	page,
	limit int64) ([]any, int64, error) {

	collection := client.Database(dbName).Collection(collectionName)

	// Criando o filtro dinâmico
	filter := bson.M{}
	if name != nil && *name != "" {
		filter["name"] = bson.M{"$regex": *name, "$options": "i"}
	}
	if key != nil && *key != "" {
		filter["key"] = bson.M{"$regex": *key, "$options": "i"}
	}

	if lastName != nil && *lastName != "" {
		filter["lastName"] = bson.M{"$regex": *lastName, "$options": "i"}
	}

	if companyCode != nil && *companyCode != "" {
		filter["companyCode"] = *companyCode
	}

	if status != nil && *status != "" {
		filter["status"] = *status
	}

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
		var data models.PurcharseRecord
		if err := cursor.Decode(&data); err != nil {
			return nil, 0, fmt.Errorf("erro ao decodificar usuário: %v", err)
		}

		excelData = append(excelData, map[string]any{
			"ID":            data.ID.Hex(), // Convertendo para string
			"Key":           data.Key,
			"Name":          data.Name,
			"LastName":      data.LastName,
			"FileName":      data.FileName,
			"CompanyCode":   data.CompanyCode,
			"CompanyName":   data.CompanyName,
			"MessageReturn": data.MessageReturn,
			"Status":        data.Status,
			"DtImportacao":  data.CreatedAt,
			"Active":        data.Active,
		})
	}

	// Retorna usuários e total de registros
	return excelData, total, nil
}

/*
Função criada por Ricardo Silva Ferreira
Inicio da criação 04/09/2025 21:36
Data Final da criação : 04/09/2025 21:36
*/
// Função para inserir um usuário na coleção "user"
func InsertPurcharseRecord(client *mongo.Client, dbName, collectionName string, data models.PurcharseRecord) error {
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

/*
Função criada por Ricardo Silva Ferreira
Inicio da criação 04/09/2025 21:37
Data Final da criação : 04/09/2025 21:38
*/
func GetPurcharseRecordByID(client *mongo.Client, dbName, collectionName, excelId string) (map[string]any, error) {

	collection := client.Database(dbName).Collection(collectionName)

	objectID, erroId := primitive.ObjectIDFromHex(excelId)
	if erroId != nil {
		log.Fatalf("Erro ao converter string para ObjectID: %v", erroId)
	}

	filter := bson.M{"_id": objectID}

	// Variável para armazenar o usuário retornado
	var excelData models.PurcharseRecord

	// Usar FindOne para pegar apenas um único registro
	err := collection.FindOne(context.Background(), filter).Decode(&excelData)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("plano de contas não encontrado")
		}
		return nil, fmt.Errorf("erro ao buscar plano de contas: %v", err)
	}

	// Converter o _id para string

	// Retornar o usuário como um mapa
	excelDatas := map[string]any{
		"ID":            excelData.ID.Hex(), // Agora o campo ID é uma string
		"Key":           excelData.Key,
		"Name":          excelData.Name,
		"LastName":      excelData.LastName,
		"FileName":      excelData.FileName,
		"Status":        excelData.Status,
		"CompanyCode":   excelData.CompanyCode,
		"CompanyName":   excelData.CompanyName,
		"MessageReturn": excelData.MessageReturn,
		"DtImportacao":  excelData.CreatedAt,
		"Active":        excelData.Active,
	}

	return excelDatas, nil
}

/*
Função criada por Ricardo Silva Ferreira
Inicio da criação 04/09/2025 21:39
Data Final da criação : 04/09/2025 21:40
*/
func GetAllPurcharseRecord(client *mongo.Client, dbName, collectionName string, page, limit int) ([]any, int, error) {
	collection := db.GetCollection(client, dbName, collectionName)

	// Criar o filtro (por enquanto vazio, pode ser expandido)
	filter := bson.M{}

	// Obter a contagem total de usuários antes da paginação
	total, err := collection.CountDocuments(context.Background(), filter)
	if err != nil {
		return nil, 0, fmt.Errorf("erro ao contar documentos: %v", err)
	}

	// Definir opções de busca com paginação
	options := options.Find()
	options.SetLimit(int64(limit))
	options.SetSkip(int64((page - 1) * limit))

	// Buscar usuários com paginação
	cursor, err := collection.Find(context.Background(), filter, options)
	if err != nil {
		return nil, 0, fmt.Errorf("erro ao buscar usuários: %v", err)
	}
	defer cursor.Close(context.Background())

	var ccs []any
	for cursor.Next(context.Background()) {
		var cc models.PurcharseRecord
		if err := cursor.Decode(&cc); err != nil {
			return nil, 0, fmt.Errorf("erro ao decodificar centro de custo: %v", err)
		}

		// Adiciona os usuários formatados
		ccs = append(ccs, map[string]any{
			"ID":            cc.ID.Hex(), // Agora o campo ID é uma string
			"Key":           cc.Key,
			"Name":          cc.Name,
			"LastName":      cc.LastName,
			"FileName":      cc.FileName,
			"CompanyCode":   cc.CompanyCode,
			"Status":        cc.Status,
			"CompanyName":   cc.CompanyName,
			"MessageReturn": cc.MessageReturn,
			"DtImportacao":  cc.CreatedAt,
			"Active":        cc.Active,
		})
	}

	if err := cursor.Err(); err != nil {
		return nil, 0, fmt.Errorf("erro ao iterar no cursor: %v", err)
	}

	return ccs, int(total), nil
}

/*
Função criada por Ricardo Silva Ferreira
Inicio da criação 30/09/2025 17:07
Data Final da criação :  30/09/2025 17:10
*/
// Função para obter todos os diários para carregar o drop de buscar
func GetImportStatus(client *mongo.Client, dbName, collectionName string) ([]any, error) {
	collection := db.GetCollection(client, dbName, collectionName)

	// Consultar todos os documentos
	cursor, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar usuários: %v", err)
	}
	defer cursor.Close(context.Background())

	var dadosBanco []any
	for cursor.Next(context.Background()) {
		var cc models.StatusImport
		if err := cursor.Decode(&cc); err != nil {
			return nil, fmt.Errorf("erro ao decodificar ,status de importação: %v", err)
		}

		// Converter o _id do MongoDB para string para retorno
		Id := cc.ID
		// Preenche o usuário com o ID convertido em string
		dadosBanco = append(dadosBanco, map[string]any{
			"ID":   Id, // Agora o campo ID é uma string
			"name": cc.Name,
			"code": cc.Code,
		})

	}

	// Verifica se houve algum erro durante a iteração do cursor
	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("erro ao iterar no cursor: %v", err)
	}

	// Retorna os usuários
	return dadosBanco, nil
}

func GroupByCompanyNameFiltered(client *mongo.Client, dbName, collectionName string, startDate, endDate *time.Time, status, companyCode string) ([]bson.M, error) {
	collection := db.GetCollection(client, dbName, collectionName)

	// Construir filtro dinamicamente
	matchConditions := bson.D{}

	if startDate != nil && endDate != nil {
		matchConditions = append(matchConditions, bson.E{
			Key: "createdAt", Value: bson.D{
				{Key: "$gte", Value: *startDate},
				{Key: "$lte", Value: *endDate},
			},
		})
	}

	if status != "" {
		matchConditions = append(matchConditions, bson.E{Key: "status", Value: status})
	}

	if companyCode != "" {
		matchConditions = append(matchConditions, bson.E{Key: "companyCode", Value: companyCode})
	}

	pipeline := mongo.Pipeline{}

	// Adiciona $match se houver condições
	if len(matchConditions) > 0 {
		pipeline = append(pipeline, bson.D{{Key: "$match", Value: matchConditions}})
	}

	// Etapa de agrupamento
	groupStage := bson.D{
		{Key: "$group", Value: bson.D{
			{Key: "_id", Value: "$companyName"},
			{Key: "total", Value: bson.D{{Key: "$sum", Value: 1}}},
			{Key: "documentos", Value: bson.D{{Key: "$push", Value: "$$ROOT"}}},
		}},
	}

	pipeline = append(pipeline, groupStage)

	cursor, err := collection.Aggregate(context.Background(), pipeline)
	if err != nil {
		return nil, fmt.Errorf("erro ao agrupar por companyName com filtros: %v", err)
	}
	defer cursor.Close(context.Background())

	var resultados []bson.M
	if err := cursor.All(context.Background(), &resultados); err != nil {
		return nil, fmt.Errorf("erro ao decodificar resultados: %v", err)
	}

	return resultados, nil
}

/*
Função criada por Ricardo Silva Ferreira
Inicio da criação 14/10/2025 21:51
Data Final da criação : 14/10/2025 21:59
*/
func GetPurcharseRecordByStatus(client *mongo.Client, dbName, collectionName, companyCode, status string) ([]any, error) {

	collection := client.Database(dbName).Collection(collectionName)

	filter := bson.M{}

	if companyCode != "" {
		filter["companyCode"] = bson.M{"$regex": companyCode, "$options": "i"}
	}
	if status != "" {
		filter["status"] = bson.M{"$regex": status, "$options": "i"}
	}

	// Consultar todos os documentos
	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar usuários: %v", err)
	}
	defer cursor.Close(context.Background())

	var dadosBanco []any
	for cursor.Next(context.Background()) {
		var cc models.PurcharseRecord
		if err := cursor.Decode(&cc); err != nil {
			return nil, fmt.Errorf("erro ao decodificar ,centro de custo: %v", err)
		}

		// Converter o _id do MongoDB para string para retorno
		Id := cc.ID
		// Preenche o usuário com o ID convertido em string
		dadosBanco = append(dadosBanco, map[string]any{
			"ID":            Id, // Agora o campo ID é uma string
			"key":           cc.Key,
			"name":          cc.Name,
			"lastName":      cc.LastName,
			"FileName":      cc.FileName,
			"CompanyCode":   cc.CompanyCode,
			"Status":        cc.Status,
			"CompanyName":   cc.CompanyName,
			"MessageReturn": cc.MessageReturn,
			"DtImportacao":  cc.CreatedAt,
			"Active":        cc.Active,
		})

	}

	// Verifica se houve algum erro durante a iteração do cursor
	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("erro ao iterar no cursor: %v", err)
	}

	// Retorna os usuários
	return dadosBanco, nil

}

/*
Função criada por Ricardo Silva Ferreira
Inicio da criação 21/10/2025 21:09
Data Final da criação : 21/10/2025 21:10
*/
func VeryExistKey(client *mongo.Client, dbName, collectionName, key string) (bool, error) {

	collection := client.Database(dbName).Collection(collectionName)
	filter := bson.M{"key": key}
	// Variável para armazenar o usuário retornado
	var excelData models.PurcharseRecord
	exist := true
	// Usar FindOne para pegar apenas um único registro
	err := collection.FindOne(context.Background(), filter).Decode(&excelData)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			exist = false
		}
	}
	// Converter o _id para string
	return exist, nil
}

// func GroupByCompanyName(client *mongo.Client, dbName, collectionName string) ([]bson.M, error) {
// 	collection := db.GetCollection(client, dbName, collectionName)

// 	pipeline := mongo.Pipeline{
// 		{{Key: "$group", Value: bson.D{
// 			{Key: "_id", Value: "$companyName"},
// 			{Key: "total", Value: bson.D{{Key: "$sum", Value: 1}}},
// 			{Key: "documentos", Value: bson.D{{Key: "$push", Value: "$$ROOT"}}},
// 		}}},
// 	}

// 	cursor, err := collection.Aggregate(context.Background(), pipeline)
// 	if err != nil {
// 		return nil, fmt.Errorf("erro ao agrupar por companyName: %v", err)
// 	}
// 	defer cursor.Close(context.Background())

// 	var resultados []bson.M
// 	if err := cursor.All(context.Background(), &resultados); err != nil {
// 		return nil, fmt.Errorf("erro ao decodificar resultados: %v", err)
// 	}

// 	return resultados, nil
// }
