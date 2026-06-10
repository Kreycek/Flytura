package outPutInvoices

import (
	flytura "Flytura"
	"Flytura/internal/airLine"
	"Flytura/internal/db"
	"Flytura/internal/models"
	"context"
	"fmt"
	"log"
	"path/filepath"
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
			SetSort(bson.D{{Key: "serverDate", Value: -1}}),
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
			"UserNameimport":       data.UserNameImport,
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
			SetSort(bson.D{{Key: "serverDate", Value: -1}}),
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
			"UserNameimport":       data.UserNameImport,
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
		v.OriginData = "RPA"

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

func ProcessOutputInvoicesExcel(
	filePath, fileName, idUserInserted, userNameimport string,
	client *mongo.Client,
	dbName, collectionName string,
) (bool, bool, bool, int, error) {

	extensao := strings.ToLower(filepath.Ext(filePath))
	var rows [][]string

	var totalRecord int = 1
	var emptySheet bool = false
	var noSheet bool = false
	var minTotalColuns bool = false

	switch extensao {
	case ".xls":
		excelXls, err := xls.Open(filePath, "utf-8")
		if err != nil {

			return noSheet, emptySheet, minTotalColuns, totalRecord, fmt.Errorf("Erro ao tentar abrir .xls")
		}

		sheet := excelXls.GetSheet(0)
		if sheet == nil {
			noSheet = true
			return noSheet, emptySheet, minTotalColuns, totalRecord, fmt.Errorf("nenhuma aba encontrada no .xls")
		}

		firstRow := sheet.Row(0)
		numCols := int(firstRow.LastCol())

		if numCols != flytura.NumberMaxColunsSheetOutPutInvoices {
			minTotalColuns = true
			return noSheet, emptySheet, minTotalColuns, totalRecord, fmt.Errorf("Número de colunas na planilha inválido deve ter 18")
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
			c3 := strings.TrimSpace(flytura.GetXlsColSafe(r, 3))
			c4 := strings.TrimSpace(flytura.GetXlsColSafe(r, 4))
			c5 := strings.TrimSpace(flytura.GetXlsColSafe(r, 5))
			c6 := strings.TrimSpace(flytura.GetXlsColSafe(r, 6))
			c7 := strings.TrimSpace(flytura.GetXlsColSafe(r, 7))
			c8 := strings.TrimSpace(flytura.GetXlsColSafe(r, 8))
			c9 := strings.TrimSpace(flytura.GetXlsColSafe(r, 9))
			c10 := strings.TrimSpace(flytura.GetXlsColSafe(r, 10))
			c11 := strings.TrimSpace(flytura.GetXlsColSafe(r, 11))
			c12 := strings.TrimSpace(flytura.GetXlsColSafe(r, 12))
			c13 := strings.TrimSpace(flytura.GetXlsColSafe(r, 13))
			c14 := strings.TrimSpace(flytura.GetXlsColSafe(r, 14))
			c15 := strings.TrimSpace(flytura.GetXlsColSafe(r, 15))
			c16 := strings.TrimSpace(flytura.GetXlsColSafe(r, 16))
			c17 := strings.TrimSpace(flytura.GetXlsColSafe(r, 17))
			c18 := strings.TrimSpace(flytura.GetXlsColSafe(r, 18))

			rows = append(rows, []string{
				c0, c1, c2, c3, c4, c5, c6, c7, c8,
				c9, c10, c11, c12, c13, c14, c15, c16, c17, c18,
			})

		}

	case ".xlsx":
		excelXlsx, err := excelize.OpenFile(filePath)

		fmt.Println("Path ", filePath)
		if err != nil {
			return noSheet, emptySheet, minTotalColuns, totalRecord, fmt.Errorf("Erro ao tentar abrir .xlsx")
		}
		defer excelXlsx.Close()

		sheetList := excelXlsx.GetSheetList()
		if len(sheetList) == 0 {
			return noSheet, emptySheet, minTotalColuns, totalRecord, fmt.Errorf("nenhuma aba encontrada no .xls")
		}

		sheetName := sheetList[0]

		// Lê linhas com o valor cru para evitar notação científica
		rows, err = excelXlsx.GetRows(sheetName, excelize.Options{RawCellValue: true})
		if err != nil {
			return noSheet, emptySheet, minTotalColuns, totalRecord, fmt.Errorf("Erro ao ler linhas", err)
		}

		// ✅ valida o cabeçalho
		if len(rows) == 0 || len(rows[0]) != flytura.NumberMaxColunsSheetOutPutInvoices {
			minTotalColuns = true
			return noSheet, emptySheet, minTotalColuns, totalRecord, fmt.Errorf("Número de colunas na planilha inválido deve ter 18")
		}

		// Normaliza todas as linhas para pelo menos 3 colunas (A,B,C)
		// preenchendo com "" onde faltar.
		normalized := make([][]string, 0, len(rows))
		for _, row := range rows {
			a := flytura.SafeCell(row, 0)
			b := flytura.SafeCell(row, 1)
			c := flytura.SafeCell(row, 2)
			d := flytura.SafeCell(row, 3)
			e := flytura.SafeCell(row, 4)
			f := flytura.SafeCell(row, 5)
			g := flytura.SafeCell(row, 6)
			h := flytura.SafeCell(row, 7)
			i := flytura.SafeCell(row, 8)
			j := flytura.SafeCell(row, 9)
			k := flytura.SafeCell(row, 10)
			l := flytura.SafeCell(row, 11)
			m := flytura.SafeCell(row, 12)
			n := flytura.SafeCell(row, 13)
			o := flytura.SafeCell(row, 14)
			p := flytura.SafeCell(row, 15)
			q := flytura.SafeCell(row, 16)
			r := flytura.SafeCell(row, 17)
			s := flytura.SafeCell(row, 18)

			normalized = append(normalized, []string{a, b, c, d, e, f, g, h, i, j, k, l, m, n, o, p, q, r, s})
		}
		rows = normalized

	default:
		return noSheet, emptySheet, minTotalColuns, totalRecord, fmt.Errorf("Extensão do arquivo inválida")
	}

	// Verifica se a planilha possui apenas cabeçalho ou está vazia
	if len(rows) <= 1 {
		emptySheet = true
		return noSheet, emptySheet, minTotalColuns, totalRecord, fmt.Errorf("Extensão do arquivo inválida")
	}

	collection := client.Database(dbName).Collection(collectionName)
	ctx := context.Background()

	airlines, error := airLine.GetAirLines(db.MongoClient, flytura.DBName, flytura.AirlineTableName)
	if error != nil {
		log.Fatal("erro ao carregar as companhias aereas")
	}

	for i, row := range rows {

		if i == 0 {
			continue // cabeçalho
		}

		nowUTC := time.Now().UTC()
		dh, err := flytura.DiffHours(nowUTC, flytura.Fuso1, flytura.Fuso2)
		if err != nil {
			// Propaga erro em vez de panic
			return noSheet, emptySheet, minTotalColuns, totalRecord, fmt.Errorf("EErro ao calcular fuso horário")
		}

		obj := models.OutPutInvoices{
			Key:                  "",
			Status:               "",
			DtProcess:            time.Time{},
			MonthProcess:         0,
			DtFLy:                "",
			RFC:                  "",
			TransferredBaseValue: 0,
			IVAValue:             0,
			Tax:                  "",
			Rate:                 "",
			FactorType:           "",
			SubTotalValue:        0,
			TotalValue:           0,
			CreatedAt:            nowUTC.Add(-time.Duration(dh) * time.Hour),
			Active:               true,
			Ruta:                 "",
			CompanyName:          "",
			CompanyCode:          "",
			TUA:                  0,
			OtherValues:          0,
			Segment:              0,
			Ticket:               "",
			IdUserInserted:       idUserInserted,
			UserNameImport:       userNameimport,
			OriginData:           "Excel import Manual",
			ServerDate:           time.Now().UTC(),
		}

		// Preenche com segurança e remove espaços internos extras da chave
		obj.Key = flytura.CompressSpaces(row[0])

		// fmt.Println("obj.Key ", obj.Key)
		obj.Status = strings.TrimSpace(row[1]) // pode ficar vazio sem erro

		obj.DtProcess = flytura.ConvertNumberToDate(row[2])

		mp, err := flytura.ConvertMonthToNumber(flytura.CompressSpaces(row[3]))
		if err != nil {
			return noSheet, emptySheet, minTotalColuns, totalRecord, fmt.Errorf("erro ao converter o mês de processamento na linha%s", strconv.Itoa(totalRecord))
		}

		obj.MonthProcess = mp
		// Verifica se a KEY já foi importada

		obj.DtFLy = flytura.CompressSpaces(row[4])

		obj.RFC = flytura.CompressSpaces(row[5])

		tbv, err := flytura.ConvertStringToFloat64(flytura.CompressSpaces(row[6]))
		if err != nil {
			return noSheet, emptySheet, minTotalColuns, totalRecord, fmt.Errorf("erro ao converter T. Base na linha %s", strconv.Itoa(totalRecord))
		}

		obj.TransferredBaseValue = tbv

		iva, err := flytura.ConvertStringToFloat64(flytura.CompressSpaces(row[7]))
		if err != nil {
			return noSheet, emptySheet, minTotalColuns, totalRecord, fmt.Errorf("erro ao converter Iva na linha %s", strconv.Itoa(totalRecord))
		}

		obj.IVAValue = iva
		obj.Tax = flytura.CompressSpaces(row[8])
		obj.Rate = flytura.CompressSpaces(row[9])
		obj.FactorType = flytura.CompressSpaces(row[10])

		st, err := flytura.ConvertStringToFloat64(flytura.CompressSpaces(row[11]))
		if err != nil {
			return noSheet, emptySheet, minTotalColuns, totalRecord, fmt.Errorf("erro ao converter SubTotal na linha %s", strconv.Itoa(totalRecord))
		}
		obj.SubTotalValue = st

		total, err := flytura.ConvertStringToFloat64(flytura.CompressSpaces(row[12]))
		if err != nil {
			return noSheet, emptySheet, minTotalColuns, totalRecord, fmt.Errorf("erro ao converter Total na linha %s", strconv.Itoa(totalRecord))
		}
		obj.TotalValue = total

		obj.Ruta = flytura.CompressSpaces(row[13])

		obj.CompanyName = flytura.CompressSpaces(row[14])

		originCodAirline, _ := airLine.SearchAirlineByName(airlines, obj.CompanyName)

		if originCodAirline == "" {
			return noSheet, emptySheet, minTotalColuns, totalRecord, fmt.Errorf("Não existe a companhia aerea da linha %s", strconv.Itoa(totalRecord))
		}

		tua, err := flytura.ConvertStringToFloat64(flytura.CompressSpaces(row[15]))
		if err != nil {
			return noSheet, emptySheet, minTotalColuns, totalRecord, fmt.Errorf("erro ao converter TUA na linha %s", strconv.Itoa(totalRecord))
		}
		obj.TUA = tua

		otv, err := flytura.ConvertStringToFloat64(flytura.CompressSpaces(row[16]))
		if err != nil {
			return noSheet, emptySheet, minTotalColuns, totalRecord, fmt.Errorf("erro ao converter Outros Valores na linha %s", strconv.Itoa(totalRecord))
		}
		obj.OtherValues = otv

		segment, err := flytura.ConvertToInt(flytura.CompressSpaces(row[17])) // pode ficar vazio sem erro
		if err != nil {
			return noSheet, emptySheet, minTotalColuns, totalRecord, fmt.Errorf("erro ao converter segmento na linha %s", strconv.Itoa(totalRecord))
		}
		obj.Segment = segment

		obj.Ticket = flytura.CompressSpaces(row[18])

		totalRecord++
	}

	totalRecord = 0

	for i, row := range rows {

		if i == 0 {
			continue // cabeçalho
		}

		nowUTC := time.Now().UTC()
		dh, err := flytura.DiffHours(nowUTC, flytura.Fuso1, flytura.Fuso2)
		if err != nil {
			// Propaga erro em vez de panic
			return noSheet, emptySheet, minTotalColuns, totalRecord, fmt.Errorf("EErro ao calcular fuso horário")
		}

		obj := models.OutPutInvoices{
			Key:                  "",
			Status:               "",
			DtProcess:            time.Time{},
			MonthProcess:         0,
			DtFLy:                "",
			RFC:                  "",
			TransferredBaseValue: 0,
			IVAValue:             0,
			Tax:                  "",
			Rate:                 "",
			FactorType:           "",
			SubTotalValue:        0,
			TotalValue:           0,
			CreatedAt:            nowUTC.Add(-time.Duration(dh) * time.Hour),
			Active:               true,
			Ruta:                 "",
			CompanyName:          "",
			CompanyCode:          "",
			TUA:                  0,
			OtherValues:          0,
			Segment:              0,
			Ticket:               "",
			IdUserInserted:       idUserInserted,
			UserNameImport:       userNameimport,
			OriginData:           "Excel import Manual",
			ServerDate:           time.Now().UTC(),
		}

		// Preenche com segurança e remove espaços internos extras da chave
		obj.Key = flytura.CompressSpaces(row[0])

		// fmt.Println("obj.Key ", obj.Key)
		obj.Status = strings.TrimSpace(row[1]) // pode ficar vazio sem erro

		obj.DtProcess = flytura.ConvertNumberToDate(row[2])
		mp, _ := flytura.ConvertMonthToNumber(flytura.CompressSpaces(row[3]))
		obj.MonthProcess = mp
		// Verifica se a KEY já foi importada

		obj.DtFLy = flytura.CompressSpaces(row[4])
		obj.RFC = flytura.CompressSpaces(row[5])
		tbv, _ := flytura.ConvertStringToFloat64(flytura.CompressSpaces(row[6]))
		obj.TransferredBaseValue = tbv
		iva, _ := flytura.ConvertStringToFloat64(flytura.CompressSpaces(row[7]))
		obj.IVAValue = iva
		obj.Tax = flytura.CompressSpaces(row[8])
		obj.Rate = flytura.CompressSpaces(row[9])
		obj.FactorType = flytura.CompressSpaces(row[10])
		st, _ := flytura.ConvertStringToFloat64(flytura.CompressSpaces(row[11]))
		obj.SubTotalValue = st
		total, _ := flytura.ConvertStringToFloat64(flytura.CompressSpaces(row[12]))
		obj.TotalValue = total
		obj.Ruta = flytura.CompressSpaces(row[13])

		originCodAirline, originNameAirline := airLine.SearchAirlineByName(airlines, flytura.CompressSpaces(row[14]))

		obj.CompanyName = originNameAirline
		obj.CompanyCode = originCodAirline

		tua, _ := flytura.ConvertStringToFloat64(flytura.CompressSpaces(row[15]))
		obj.TUA = tua
		otv, _ := flytura.ConvertStringToFloat64(flytura.CompressSpaces(row[16]))
		obj.OtherValues = otv
		segment, _ := flytura.ConvertToInt(flytura.CompressSpaces(row[17])) // pode ficar vazio sem erro
		obj.Segment = segment
		obj.Ticket = flytura.CompressSpaces(row[18])
		if _, err := collection.InsertOne(ctx, obj); err != nil {
			log.Println("Erro ao inserir:", err)
			continue
		}
		totalRecord++
	}
	return noSheet, emptySheet, minTotalColuns, totalRecord, nil
}
