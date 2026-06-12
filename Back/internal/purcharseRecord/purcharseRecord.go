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
	"strconv"
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
Inicio da criação 17/11/2025 15:17
Data Final da criação : 17/11/2025 15:22
*/
func ReturnEmptyRowNumber(rows [][]string) (int64, string) {

	var cont int64 = 0
	emptyErro := ""

	for i, row := range rows {

		//A linha abaixo retira todos os espaços
		spaceRegex := regexp.MustCompile(`\s+`)
		key := spaceRegex.ReplaceAllString(row[0], "")
		name := spaceRegex.ReplaceAllString(row[1], "")
		lastName := spaceRegex.ReplaceAllString(row[2], "")

		// fmt.Print("key, name, lastName ", key, name, lastName)
		//se o key for vazio mas algum campo tiver dados é um erro de vazio
		if key == "" && (name != "" || lastName != "") {
			cont++
			emptyErro += strconv.Itoa(i+1) + ", "
			continue //Caso todos os campos estiverem vazio desconsidera a linha
		} else if key == "" && name == "" && lastName == "" {
			continue
		}
	}

	return cont, emptyErro
}

/*
Função criada por Ricardo Silva Ferreira
Inicio da criação 15/11/2025 11:10
Data Final da criação : 15/11/2025 11:40
*/
func ReturnSheetErrorLine(rows [][]string) (string, string, string) {

	var emptyError string = ""
	var cientificError string = ""
	var ivalidFormatError string = ""

	// Regex para espaços e para formato científico (aceita vírgula ou ponto)
	spaceRegex := regexp.MustCompile(`\s+`)
	scientificRegex := regexp.MustCompile(`(?i)^[0-9]+([.,][0-9]+)?E[+-]?[0-9]+$`)

	for i, row := range rows {
		if i == 0 {
			continue // cabeçalho
		}

		key := spaceRegex.ReplaceAllString(row[0], "")
		name := spaceRegex.ReplaceAllString(row[1], "")
		lastName := spaceRegex.ReplaceAllString(row[2], "")

		//se o key for vazio mas algum campo tiver dados é um erro de vazio
		if key == "" && (name != "" || lastName != "") {
			emptyError = emptyError + strconv.Itoa(i+1) + ", "
			continue
			//Caso todos os campos estiverem vazio desconsidera a linha
		} else if key == "" && name == "" && lastName == "" {
			continue
		}

		// Se for formato científico, considera erro
		if scientificRegex.MatchString(key) {

			cientificError = cientificError + strconv.Itoa(i+1) + ", "
			// alerts = append(alerts, "O valor da linha "+strconv.Itoa(i+1)+" está em notação científica e não é permitido")
			continue
		}

		_, err := strconv.ParseFloat(strings.ReplaceAll(key, ",", "."), 64) // normaliza vírgula para ponto
		if err != nil {
			msg := key
			if msg == "" {
				msg = "vazio"
			}
			ivalidFormatError = ivalidFormatError + strconv.Itoa(i+1) + ", "
			// alerts = append(alerts, "O valor da linha "+strconv.Itoa(i+1)+" está como "+msg+" formate ou preencha")
		}
	}

	return emptyError, cientificError, ivalidFormatError
}

// getXlsColSafe retorna a coluna idx de uma linha .xls com segurança.
// Se a coluna não existir, devolve "".
/*
	Função criada por Ricardo Silva Ferreira
	Inicio da criação 01/03/2023 19:05
	Data Final da criação : 01/03/2023 19:06
	Local: Brasil
*/
// func getXlsColSafe(r *xls.Row, idx int) string {
// 	if r == nil {
// 		return ""
// 	}
// 	// LastCol geralmente é a contagem de colunas válidas (0..LastCol-1)
// 	if idx >= 0 && idx < r.LastCol() {
// 		return r.Col(idx)
// 	}
// 	return ""
// }

// /*
// Função criada por Ricardo Silva Ferreira
// Inicio da criação 01/03/2023 19:05
// Data Final da criação : 01/03/2023 19:06
// Local: Brasil
// */
// // safeCell retorna a célula idx de uma []string (excelize GetRows)
// // Se a coluna não existir, devolve "".
// func safeCell(row []string, idx int) string {
// 	if idx >= 0 && idx < len(row) {
// 		return strings.TrimSpace(row[idx])
// 	}
// 	return ""
// }

//
// Função principal refatorada
//
/*
Função criada por Ricardo Silva Ferreira
Inicio da criação 03/09/2025 22:20
Data Final da criação : 04/09/2025 18:50
Obs: 21/03/2026 00:20 -> Acrescentado idUserInserted para sabermos que inseriu o registro
Obs: 26/03/2026 19:22 -> Acrescentado NumberMaxColunsSheet para verificar o número de colunas na planilha
*/

func ProcessPurcharseRecordExcel(
	filePath,
	fileName,
	companyName,
	companyCode,
	idUserInserted,
	userNameInserted string,
	importSheetOnlyVerifyNumber bool,
	client *mongo.Client,
	dbName, collectionName string,
) (int64, int64, bool, string, string, string, bool, error) {

	extensao := strings.ToLower(filepath.Ext(filePath))
	var rows [][]string

	var totalRecord int64 = 0
	var emptySheet = false
	var ivalidNumberCols = false

	switch extensao {
	case ".xls":
		excelXls, err := xls.Open(filePath, "utf-8")
		if err != nil {

			return 0, 0, emptySheet, "", "", "", ivalidNumberCols, err
		}

		sheet := excelXls.GetSheet(0)
		if sheet == nil {
			return 0, 0, emptySheet, "", "", "", ivalidNumberCols, fmt.Errorf("nenhuma aba encontrada no .xls")
		}

		firstRow := sheet.Row(0)
		numCols := int(firstRow.LastCol())

		// fmt.Println("numCols ", numCols)

		if numCols != flytura.NumberMaxColunsSheet {
			ivalidNumberCols = true
			return 0, 0, emptySheet,
				"O cabeçalho da planilha deve conter exatamente 3 colunas (A, B, C)",
				"", "", ivalidNumberCols, nil
		}

		// fmt.Println("firstRow 1 ", firstRow.Col(0))
		// fmt.Println("firstRow 2 ", firstRow.Col(1))
		// fmt.Println("firstRow 3 ", firstRow.Col(2))

		// Lê A,B,C com segurança para cada linha
		for i := 0; i <= int(sheet.MaxRow); i++ {
			r := sheet.Row(i)
			if r == nil {
				rows = append(rows, []string{"", "", ""})
				continue
			}
			c0 := strings.TrimSpace(flytura.GetXlsColSafe(r, 0))
			c1 := strings.TrimSpace(flytura.GetXlsColSafe(r, 1))
			c2 := strings.TrimSpace(flytura.GetXlsColSafe(r, 2))
			rows = append(rows, []string{c0, c1, c2})
		}

	case ".xlsx":
		excelXlsx, err := excelize.OpenFile(filePath)
		if err != nil {
			return 0, 0, emptySheet, "", "", "", ivalidNumberCols, err
		}
		defer excelXlsx.Close()

		sheetList := excelXlsx.GetSheetList()
		if len(sheetList) == 0 {
			return 0, 0, emptySheet, "", "", "", ivalidNumberCols, fmt.Errorf("nenhuma aba encontrada no .xlsx")
		}

		sheetName := sheetList[0]
		// Lê linhas com o valor cru para evitar notação científica
		rows, err = excelXlsx.GetRows(sheetName, excelize.Options{RawCellValue: true})
		if err != nil {
			return 0, 0, emptySheet, "", "", "", ivalidNumberCols, fmt.Errorf("erro ao ler linhas da aba %s: %v", sheetName, err)
		}

		// ✅ valida o cabeçalho
		if len(rows) == 0 || len(rows[0]) != flytura.NumberMaxColunsSheet {

			ivalidNumberCols = true

			return 0, 0, emptySheet,
				"O cabeçalho da planilha deve conter exatamente 3 colunas (A, B, C)",
				"", "", ivalidNumberCols, nil
		}

		// Normaliza todas as linhas para pelo menos 3 colunas (A,B,C)
		// preenchendo com "" onde faltar.
		normalized := make([][]string, 0, len(rows))
		for _, r := range rows {
			a := flytura.SafeCell(r, 0)
			b := flytura.SafeCell(r, 1)
			c := flytura.SafeCell(r, 2)
			normalized = append(normalized, []string{a, b, c})
		}
		rows = normalized

	default:
		return 0, 0, emptySheet, "", "", "", ivalidNumberCols, fmt.Errorf("extensão de arquivo não suportada: %s", extensao)
	}

	// Verifica se a planilha possui apenas cabeçalho ou está vazia
	if len(rows) <= 1 {
		emptySheet = true
		return 0, 0, emptySheet, "", "", "", ivalidNumberCols, nil
	}

	// Aqui verifica se tem registros vazios, em notação científica ou valores inválidos
	if importSheetOnlyVerifyNumber {
		emptyError, cientificError, ivalidFormatError := ReturnSheetErrorLine(rows)
		if emptyError != "" || cientificError != "" || ivalidFormatError != "" {
			return 0, 0, emptySheet, emptyError, cientificError, ivalidFormatError, ivalidNumberCols, nil
		}
	} else {
		// Retorna o total de linhas vazias e qual é o número da linha
		emptyRecords, emptyError := ReturnEmptyRowNumber(rows)
		if emptyRecords > 0 {
			return emptyRecords, 0, emptySheet, emptyError, "", "", ivalidNumberCols, nil
		}
	}

	collection := client.Database(dbName).Collection(collectionName)
	ctx := context.Background()

	for i, row := range rows {
		if i == 0 {
			continue // cabeçalho
		}

		nowUTC := time.Now().UTC()
		dh, err := flytura.DiffHours(nowUTC, flytura.Fuso1, flytura.Fuso2)
		if err != nil {
			// Propaga erro em vez de panic
			return 0, 0, emptySheet, "", "", "", ivalidNumberCols, fmt.Errorf("erro ao calcular DiffHours: %w", err)
		}

		obj := models.PurcharseRecord{
			Key:              "",
			Name:             "",
			LastName:         "",
			FileName:         fileName,
			Status:           "Fila",
			Active:           true,
			CreatedAt:        nowUTC.Add(-time.Duration(dh) * time.Hour),
			CompanyCode:      companyCode,
			CompanyName:      companyName,
			IdUserInserted:   idUserInserted,
			NameUserInserted: userNameInserted,
			//Abaixo como vem da planilha vamos considerar sempre como IDA
			DirectionOfDestination: "GO",
		}

		// Preenche com segurança e remove espaços internos extras da chave
		obj.Key = flytura.CompressSpaces(row[0])

		// fmt.Println("obj.Key ", obj.Key)
		obj.Name = strings.TrimSpace(row[1])     // pode ficar vazio sem erro
		obj.LastName = strings.TrimSpace(row[2]) // pode ficar vazio sem erro

		if obj.Key == "" {
			continue
		}

		// Verifica se a KEY já foi importada
		exist, errVerify := VeryExistKey(client, flytura.DBName, flytura.PurcharseRecordTableName, obj.Key)
		if errVerify != nil {
			log.Println("Erro ao verificar existência de chave:", errVerify)
			continue
		}

		if !exist {
			if _, err := collection.InsertOne(ctx, obj); err != nil {
				log.Println("Erro ao inserir:", err)
				continue
			}
			totalRecord++
		}
	}

	return 0, totalRecord, emptySheet, "", "", "", ivalidNumberCols, nil
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
			"ID":                     Id, // Agora o campo ID é uma string
			"key":                    cc.Key,
			"name":                   cc.Name,
			"lastName":               cc.LastName,
			"FileName":               cc.FileName,
			"CompanyCode":            cc.CompanyCode,
			"Status":                 cc.Status,
			"CompanyName":            cc.CompanyName,
			"DtImportacao":           cc.CreatedAt,
			"Active":                 cc.Active,
			"EmissionDate":           cc.EmissionDate,
			"DirectionOfDestination": cc.DirectionOfDestination,
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
Modificado em:  17/03/2026 23:24 adicionado SetSort
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
	statusText *string,
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

	//O CAMPO NA TELA ESTÁ DESATIVADO, OCULTO
	if statusText != nil && *statusText != "" {
		term := flytura.AccentAgnosticSpaceInsensitivePattern(*statusText) // minúsculas + sem acentos
		filter["messageReturn"] = bson.M{
			"$regex": primitive.Regex{
				Pattern: term,
				Options: "i",
			},
		}
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
	//Modificado bloco abaixo adicionado SetSort em 17/03/2026 23:24
	cursor, err := collection.Find(
		context.Background(),
		filter,
		options.Find().
			SetSkip(int64((page-1)*limit)).
			SetLimit(limit).
			SetSort(bson.D{{Key: "createdAt", Value: -1},
				{Key: "_id", Value: -1}}),
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
			"ID":                     data.ID.Hex(), // Convertendo para string
			"Key":                    data.Key,
			"Name":                   data.Name,
			"LastName":               data.LastName,
			"FileName":               data.FileName,
			"CompanyCode":            data.CompanyCode,
			"CompanyName":            data.CompanyName,
			"MessageReturn":          data.MessageReturn,
			"EmissionDate":           data.EmissionDate,
			"Status":                 data.Status,
			"DtImportacao":           data.CreatedAt,
			"Active":                 data.Active,
			"DirectionOfDestination": data.DirectionOfDestination,
			"NameUserInserted":       data.NameUserInserted,
		})
	}

	// Retorna usuários e total de registros
	return excelData, total, nil
}

/*
Função criada por Ricardo Silva Ferreira
Inicio da criação 09/04/2026 08:44
Data Final da criação : 09/04/2026 08:46
*/
func SearchPurcharseRecordByPeriod(
	client *mongo.Client,
	dbName, collectionName string,
	key *string,
	name *string,
	lastName *string,
	companyCode *string,
	startDate *time.Time,
	endDate *time.Time,
	status *string,
	statusText *string) ([]any, int64, error) {

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

	//O CAMPO NA TELA ESTÁ DESATIVADO, OCULTO
	if statusText != nil && *statusText != "" {
		term := flytura.AccentAgnosticSpaceInsensitivePattern(*statusText) // minúsculas + sem acentos
		filter["messageReturn"] = bson.M{
			"$regex": primitive.Regex{
				Pattern: term,
				Options: "i",
			},
		}
	}

	if status != nil && *status != "" {
		filter["status"] = *status
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
	//Modificado bloco abaixo adicionado SetSort em 09/04/2026 09:02
	cursor, err := collection.Find(
		context.Background(),
		filter,
		options.Find().
			SetSort(bson.D{{Key: "createdAt", Value: -1}}),
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
			"ID":                     data.ID.Hex(), // Convertendo para string
			"Key":                    data.Key,
			"Name":                   data.Name,
			"LastName":               data.LastName,
			"FileName":               data.FileName,
			"CompanyCode":            data.CompanyCode,
			"CompanyName":            data.CompanyName,
			"MessageReturn":          data.MessageReturn,
			"EmissionDate":           data.EmissionDate,
			"Status":                 data.Status,
			"DtImportacao":           data.CreatedAt,
			"Active":                 data.Active,
			"DirectionOfDestination": data.DirectionOfDestination,
			"NameUserInserted":       data.NameUserInserted,
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
Inicio da criação 04/09/2025 21:36
Data Final da criação : 04/09/2025 21:36
*/
// Função para inserir um usuário na coleção "user"
func InsertMultiplePurcharseRecord(client *mongo.Client, dbName, collectionName string, data []interface{}) error {
	collection := client.Database(dbName).Collection(collectionName)

	// Criar um contexto para a operação de inserção
	ctx := context.Background()

	opts := options.InsertMany().SetOrdered(false)
	// Inserir o documento
	_, err := collection.InsertMany(ctx, data, opts)
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
		"ID":                     excelData.ID.Hex(), // Agora o campo ID é uma string
		"Key":                    excelData.Key,
		"Name":                   excelData.Name,
		"LastName":               excelData.LastName,
		"FileName":               excelData.FileName,
		"Status":                 excelData.Status,
		"CompanyCode":            excelData.CompanyCode,
		"EmissionDate":           excelData.EmissionDate,
		"CompanyName":            excelData.CompanyName,
		"MessageReturn":          excelData.MessageReturn,
		"DtImportacao":           excelData.CreatedAt,
		"Active":                 excelData.Active,
		"DirectionOfDestination": excelData.DirectionOfDestination,
		"NameUserInserted":       excelData.NameUserInserted,
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
	//Adicionada a linha abaixo em  17/03/2025 23:20
	options.SetSort(bson.D{{Key: "createdAt", Value: -1}})

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
			"ID":                     cc.ID.Hex(), // Agora o campo ID é uma string
			"Key":                    cc.Key,
			"Name":                   cc.Name,
			"LastName":               cc.LastName,
			"FileName":               cc.FileName,
			"CompanyCode":            cc.CompanyCode,
			"Status":                 cc.Status,
			"EmissionDate":           cc.EmissionDate,
			"CompanyName":            cc.CompanyName,
			"MessageReturn":          cc.MessageReturn,
			"DtImportacao":           cc.CreatedAt,
			"Active":                 cc.Active,
			"DirectionOfDestination": cc.DirectionOfDestination,
			"NameUserInserted":       cc.NameUserInserted,
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

/*
Função criada por Ricardo Silva Ferreira
Inicio da criação 30/09/2025 19:07
Data Final da criação :  30/09/2025 19:45
Data Alteração : 13/02/2026 15:28 ->Melhoria na busca para deixa-la mais rápida
*/
func GroupByCompanyNameFiltered(
	client *mongo.Client,
	dbName, collectionName string,
	startDate, endDate *time.Time,
	status, companyCode string,
) ([]bson.M, error) {

	collection := db.GetCollection(client, dbName, collectionName)

	match := bson.D{}
	if startDate != nil && endDate != nil {
		match = append(match, bson.E{Key: "createdAt", Value: bson.D{
			{Key: "$gte", Value: *startDate},
			{Key: "$lt", Value: *endDate}, // fim exclusivo
		}})
	}
	if status != "" {
		match = append(match, bson.E{Key: "status", Value: status})
	}
	if companyCode != "" {
		match = append(match, bson.E{Key: "companyCode", Value: companyCode})
	}

	pipeline := mongo.Pipeline{}

	if len(match) > 0 {
		pipeline = append(pipeline, bson.D{{Key: "$match", Value: match}})
	}

	// (Opcional) projetar somente campos necessários se for empurrar documentos
	// pipeline = append(pipeline, bson.D{{Key: "$project", Value: bson.D{
	//  {Key: "companyName", Value: 1},
	//  {Key: "createdAt", Value: 1},
	//  {Key: "status", Value: 1},
	//  {Key: "companyCode", Value: 1},
	// }}})

	groupStage := bson.D{{
		Key: "$group", Value: bson.D{
			{Key: "_id", Value: "$companyName"},
			{Key: "total", Value: bson.D{{Key: "$sum", Value: 1}}},
			// Remova se não for necessário:
			// {Key: "documentos", Value: bson.D{{Key: "$push", Value: "$$ROOT"}}},
		},
	}}
	pipeline = append(pipeline, groupStage)

	// (Opcional) ordenar por total desc
	pipeline = append(pipeline, bson.D{{Key: "$sort", Value: bson.D{{Key: "total", Value: -1}}}})

	// (Opcional) limitar nº de grupos retornados
	// pipeline = append(pipeline, bson.D{{Key: "$limit", Value: 100}})

	// allowDiskUse só se necessário (melhor evitar se o pipeline couber em memória)
	opts := options.Aggregate().SetAllowDiskUse(false)

	ctx := context.Background()
	cursor, err := collection.Aggregate(ctx, pipeline, opts)
	if err != nil {
		return nil, fmt.Errorf("erro na agregação: %w", err)
	}
	defer cursor.Close(ctx)

	var resultados []bson.M
	if err := cursor.All(ctx, &resultados); err != nil {
		return nil, fmt.Errorf("erro ao decodificar resultados: %w", err)
	}
	return resultados, nil
}

/*
Função criada por Ricardo Silva Ferreira
Inicio da criação 14/10/2025 21:51
Data Final da criação : 14/10/2025 21:59
Data Alteração : 28/10/2025 20:36 ->Acrescentada busca por data a pedido do kaique
Data Alteração : 11/03/2026 17:39 ->Acrescentada busca por texto do erro no campo messageReturn
*/
func GetPurcharseRecordByStatus(client *mongo.Client, dbName, collectionName, companyCode, status string, startDateStr string, endDateStr string, messageReturn string) ([]any, error) {

	collection := client.Database(dbName).Collection(collectionName)

	filter := bson.M{}

	if companyCode != "" {
		filter["companyCode"] = bson.M{"$regex": companyCode, "$options": "i"}
	}

	if status != "" {
		filter["status"] = bson.M{"$regex": status, "$options": "i"}
	}

	if messageReturn != "" {
		term := flytura.AccentAgnosticSpaceInsensitivePattern(messageReturn) // minúsculas + sem acentos
		filter["messageReturn"] = bson.M{
			"$regex": primitive.Regex{
				Pattern: term,
				Options: "i",
			},
		}
	}

	if startDateStr != "" || endDateStr != "" {

		dateFilter := bson.M{}

		if startDateStr != "" {
			startDate, err := time.Parse("2006-01-02", startDateStr)
			if err == nil {
				start := time.Date(startDate.Year(), startDate.Month(), startDate.Day(), 0, 0, 0, 0, startDate.Location())
				dateFilter["$gte"] = start
			}
		}

		if endDateStr != "" {
			endDate, err := time.Parse("2006-01-02", endDateStr)
			if err == nil {
				end := time.Date(endDate.Year(), endDate.Month(), endDate.Day(), 23, 59, 59, int(time.Second-time.Nanosecond), endDate.Location())
				dateFilter["$lte"] = end
			}
		}

		filter["createdAt"] = dateFilter
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
			"ID":                     Id, // Agora o campo ID é uma string
			"key":                    cc.Key,
			"name":                   cc.Name,
			"lastName":               cc.LastName,
			"FileName":               cc.FileName,
			"CompanyCode":            cc.CompanyCode,
			"Status":                 cc.Status,
			"CompanyName":            cc.CompanyName,
			"EmissionDate":           cc.EmissionDate,
			"MessageReturn":          cc.MessageReturn,
			"DtImportacao":           cc.CreatedAt,
			"Active":                 cc.Active,
			"DirectionOfDestination": cc.DirectionOfDestination,
			"NameUserInserted":       cc.NameUserInserted,
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
Data: 06/06/2026 23:10 ignorar espaços
*/
func VeryExistKey(client *mongo.Client, dbName, collectionName, key string) (bool, error) {

	collection := client.Database(dbName).Collection(collectionName)

	filter := bson.M{
		"$expr": bson.M{
			"$eq": []interface{}{
				bson.M{"$trim": bson.M{"input": "$key"}},
				key,
			},
		},
	}

	var excelData models.PurcharseRecord

	err := collection.FindOne(context.Background(), filter).Decode(&excelData)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

/*
Função criada por Ricardo Silva Ferreira
Inicio da criação 27/11/2025 12:58
Data Final da criação : 27/11/2025 13:00
*/
func DeletePurcharseRecordByID(client *mongo.Client, dbName, collectionName, id string) error {
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

/*
Função criada por Ricardo Silva Ferreira
Inicio da criação 16/04/2026 14:06
Data Final da criação : 16/04/2026 14:06
*/
// Função para obter todos os diários para carregar o drop de buscar
func GetAllImportStatus(client *mongo.Client, dbName, collectionName string) ([]any, error) {
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
			return nil, fmt.Errorf("erro ao decodificar ,centro de custo: %v", err)
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

/*
Função criada por Ricardo Silva Ferreira
Inicio da criação 27/04/2026 15:52
Data Final da criação : 27/04/2026 15:52
*/
// Função para obter todos os diários para carregar o drop de buscar
func AgregateByImportDateAndStatus(
	client *mongo.Client,
	dbName, collectionName string,
	startDate, endDate *time.Time,
	status, companyCode string,
) ([]any, error) {

	collection := db.GetCollection(client, dbName, collectionName)

	// 🔹 Monta o $match dinamicamente
	match := bson.D{}

	if startDate != nil && endDate != nil {
		match = append(match, bson.E{
			Key: "createdAt",
			Value: bson.D{
				{Key: "$gte", Value: *startDate},
				{Key: "$lt", Value: *endDate},
			},
		})
	}

	fmt.Println("date 111", startDate, endDate)

	if status != "" {
		match = append(match, bson.E{Key: "status", Value: status})
	}

	if companyCode != "" {
		match = append(match, bson.E{Key: "companyCode", Value: companyCode})
	}

	// 🔹 Monta o pipeline
	pipeline := mongo.Pipeline{}

	// ✅ $match vem PRIMEIRO
	if len(match) > 0 {
		pipeline = append(pipeline, bson.D{
			{Key: "$match", Value: match},
		})
	}

	// ✅ $project: normaliza a data
	pipeline = append(pipeline,
		bson.D{
			{Key: "$project", Value: bson.M{
				"status": "$status",
				"date": bson.M{
					"$dateToString": bson.M{
						"format": "%Y-%m-%d",
						"date":   "$createdAt",
					},
				},
			}},
		},
		bson.D{
			{Key: "$group", Value: bson.M{
				"_id": bson.M{
					"date":   "$date",
					"status": "$status",
				},
				"total": bson.M{"$sum": 1},
			}},
		},
		bson.D{
			{Key: "$sort", Value: bson.M{
				"_id.date":   1,
				"_id.status": 1,
			}},
		},
		bson.D{
			{Key: "$project", Value: bson.M{
				"_id":    0,
				"date":   "$_id.date",
				"status": "$_id.status",
				"total":  1,
			}},
		},
	)

	cursor, err := collection.Aggregate(context.Background(), pipeline)
	if err != nil {
		return nil, fmt.Errorf("erro ao agregar registros: %v", err)
	}
	defer cursor.Close(context.Background())

	var result []any
	if err := cursor.All(context.Background(), &result); err != nil {
		return nil, fmt.Errorf("erro ao ler resultado: %v", err)
	}

	return result, nil
}
