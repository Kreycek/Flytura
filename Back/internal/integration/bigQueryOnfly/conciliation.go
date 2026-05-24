package bigQueryOnfly

import (
	flytura "Flytura"
	"Flytura/internal/airLine"
	"Flytura/internal/db"
	"Flytura/internal/models"
	"Flytura/internal/purcharseRecord"
	"context"
	"fmt"
	"log"
	"math/big"
	"strconv"
	"strings"
	"time"

	"cloud.google.com/go/bigquery"
	"cloud.google.com/go/civil"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/api/iterator"
)

/*
Função criada por Ricardo Silva Ferreira
Inicio da criação 10/05/2026 18:00
Data Final da criação 10/05/2026 18:00
*/
func convertToString(row map[string]bigquery.Value, fieldName string) string {

	result := ""
	if v, ok := row[fieldName].(string); ok {
		result = v
	}

	return result
}

/*
Função criada por Ricardo Silva Ferreira
Inicio da criação 10/05/2026 18:00
Data Final da criação 10/05/2026 18:00
*/
func convertToTimeDate(row map[string]bigquery.Value, fieldName string) time.Time {

	if d, ok := row[fieldName].(civil.Date); ok {
		t := time.Date(
			d.Year,
			time.Month(d.Month),
			d.Day,
			0, 0, 0, 0,
			time.UTC,
		)

		return t
	} else {
		return time.Time{}
	}

}

/*
Função criada por Ricardo Silva Ferreira
Inicio da criação 10/05/2026 19:40
Data Final da criação 10/05/2026 19:40
*/
func convertToDecimal128(
	row map[string]bigquery.Value,
	fieldName string,
) primitive.Decimal128 {

	v, exists := row[fieldName]
	if !exists || v == nil {
		return primitive.NewDecimal128(0, 0)
	}

	switch value := v.(type) {

	// ✅ BigQuery NUMERIC / BIGNUMERIC
	case *big.Rat:
		d128, err := primitive.ParseDecimal128(
			value.FloatString(2),
		)
		if err != nil {
			return primitive.NewDecimal128(0, 0)
		}
		return d128

	// ✅ JSON / Excel como string ("3402.66" ou "3402,66")
	case string:
		value = strings.TrimSpace(value)
		if value == "" {
			return primitive.NewDecimal128(0, 0)
		}

		value = strings.Replace(value, ",", ".", 1)

		d128, err := primitive.ParseDecimal128(value)
		if err != nil {
			return primitive.NewDecimal128(0, 0)
		}
		return d128

	// ✅ Caso venha como número (menos comum, mas acontece)
	case float64:
		d128, err := primitive.ParseDecimal128(
			strconv.FormatFloat(value, 'f', -1, 64),
		)
		if err != nil {
			return primitive.NewDecimal128(0, 0)
		}
		return d128

	default:
		// tipo inesperado → não quebra o fluxo
		return primitive.NewDecimal128(0, 0)
	}
}

/*
Função criada por Ricardo Silva Ferreira
Inicio da criação 12/05/2026 18:47
Data Final da criação 12/05/2026 18:47
*/
func convertToFloat64(row map[string]bigquery.Value, fieldName string) float64 {
	if v, ok := row[fieldName]; ok {
		switch value := v.(type) {

		case float64:
			return value

		case int64:
			return float64(value)

		case string:
			// remove possíveis separadores de milhar
			clean := strings.ReplaceAll(value, ",", "")
			f, err := strconv.ParseFloat(clean, 64)
			if err == nil {
				return f
			}
		}
	}

	return 0
}

/*
Função criada por Ricardo Silva Ferreira
Inicio da criação 11/05/2026 14:45
Data Final da criação 11/05/2026 14:50
*/

func insertPurchardRecordByBigQuery(key string, pr models.PurcharseRecord) {

	exist, errVerify := purcharseRecord.VeryExistKey(db.MongoClient, flytura.DBName, flytura.PurcharseRecordTableName, key)
	if errVerify != nil {
		log.Println("Erro ao inserir purcharse record by bigquery:", errVerify)
	} else {
		// fmt.Println("Existe key ", key, exist)
		if !exist {
			err := purcharseRecord.InsertPurcharseRecord(db.MongoClient, flytura.DBName, flytura.PurcharseRecordTableName, pr)
			if err != nil {

			}
		} else {

		}
	}
}

func ImportConciliationDataOnflys() {

	today := time.Now().Format("2006-01-02")

	ctx := context.Background()

	projectID := "dw-onfly-prd"

	// credPath := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
	// fmt.Println("Arquivo de credenciais:", credPath)

	client, err := bigquery.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("Erro ao criar client: %v", err)
	}
	defer client.Close()

	query := `
        SELECT 
			protocol,
			traveler_name,	
			traveler_first_name,
			traveler_last_name,		
			emission_date,
			origin_date,
			return_date,
			origin_locator,
			return_locator,
			origin_e_ticket,
			return_e_ticket,
			origin_airline,
			return_airline,
			onfly_amount_origin,
			onfly_amount_return,
			origin_status,
			origin_status_old,
			currency_code,
			origin_cancelled_reason,
			return_cancelled_reason
        FROM 
			conciliation.gold_flytura 
			--where emission_date='2026-05-14'
			where (emission_date>='2026-05-01' and emission_date<=@today) 
			--where emission_date=@today
			--AND origin_airline IN UNNEST(@companies)
        
    `

	// companies := []string{
	// 	"Aeromexico",
	// 	"Volaris",
	// 	"VivaAerobus",
	// }

	q := client.Query(query)

	q.Parameters = []bigquery.QueryParameter{
		{
			Name:  "today",
			Value: today,
		},
		// {
		// 	Name:  "companies",
		// 	Value: companies, // ARRAY de STRING
		// },
	}

	it, err := q.Read(ctx)
	if err != nil {
		log.Fatalf("Erro ao executar query: %v", err)
	}

	airlines, error := airLine.GetAirLines(db.MongoClient, flytura.DBName, flytura.AirlineTableName)
	if error != nil {
		log.Fatal("erro ao carregar as companhias aereas")
	}

	// fmt.Println("airlines", airlines)

	cont := 0
	for {
		var row map[string]bigquery.Value

		err := it.Next(&row)
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		var data models.Conciliation
		data.Protocol = convertToString(row, "protocol")
		data.EmissionDate = convertToTimeDate(row, "emission_date")
		data.OriginDate = convertToTimeDate(row, "origin_date")
		data.ReturnDate = convertToTimeDate(row, "return_date")
		data.OriginLocator = convertToString(row, "origin_locator")
		data.ReturnLocator = convertToString(row, "return_locator")
		data.OriginETicket = strings.ReplaceAll(convertToString(row, "origin_e_ticket"), "-", "")
		data.ReturnETicket = strings.ReplaceAll(convertToString(row, "return_e_ticket"), "-", "")
		data.OriginAirline = convertToString(row, "origin_airline")
		data.ReturnAirline = convertToString(row, "return_airline")
		data.TravelerName = convertToString(row, "traveler_name")
		data.TravelerFirstName = convertToString(row, "traveler_first_name")
		data.TravelerLastName = convertToString(row, "traveler_last_name")
		data.AmountOrigin = convertToFloat64(row, "onfly_amount_origin")
		data.AmountReturn = convertToFloat64(row, "onfly_amount_return")
		data.BookingStatus = convertToString(row, "origin_status")
		data.BookingStatusOld = convertToString(row, "origin_status_old")
		data.CurrencyCode = convertToString(row, "currency_code")
		data.OriginCancelledReason = convertToString(row, "origin_cancelled_reason")
		data.ReturnCancelledReason = convertToString(row, "return_cancelled_reason")
		data.Active = true

		nowUTC := time.Now().UTC()
		dh, err := flytura.DiffHours(nowUTC, flytura.Fuso1, flytura.Fuso2)
		if err != nil {
			log.Printf("erro ao calcular DiffHours: %v; usando nowUTC", err)
			data.CreatedAtContractedCountry = nowUTC
		} else {
			data.CreatedAtContractedCountry = nowUTC.Add(-time.Duration(dh) * time.Hour)
		}

		data.CreatedAtLocalCountry = nowUTC

		// fmt.Println("Total ", convertToDecimal128(row, "onfly_amount_origin"))

		// fmt.Println("Protocol:", row["protocol"])
		// fmt.Println("Emission Data:", row["emission_date"])
		// fmt.Println("traveler_first_name:", row["traveler_first_name"])
		// fmt.Println("traveler_last_name:", row["traveler_last_name"])
		// fmt.Println("Traveler:", data.TravelerName)
		// fmt.Println("OriginLocator:", data.OriginLocator)
		// fmt.Println("ReturnLocator:", data.ReturnLocator)
		// fmt.Println("OriginETicket:", data.OriginETicket)
		// fmt.Println("ReturnETicket:", data.ReturnETicket)
		// fmt.Println("Amount Origin:", row["onfly_amount_origin"])
		// fmt.Println("Amount Return:", row["onfly_amount_return"])
		// fmt.Println("Currency:", row["currency_code"])
		// fmt.Println("Emission Data:", row["emission_date"])
		// fmt.Println("Data:", data)
		// fmt.Println("--------------------------------")

		// fmt.Println("data.OriginAirline ", data.OriginAirline)

		codAirline, nameAirline := airLine.SearchAirlineByName(airlines, data.OriginAirline)

		// if data.OriginCancelledReason != "" {
		// 	fmt.Println("Achou ", data.OriginLocator)
		// 	fmt.Println("OriginStatus ", data.BookingStatus)
		// 	fmt.Println("codAirline ", codAirline)
		// 	fmt.Println("nameAirline ", nameAirline)
		// 	fmt.Println("OriginAirline", data.OriginAirline)
		// }
		if flytura.Normalize(data.BookingStatus) == "emitted" {

			if codAirline == "0001" || codAirline == "0002" || codAirline == "0003" {

				var pr models.PurcharseRecord

				// name, lastName := flytura.SplitNameLastName(data.TravelerName)

				pr.Name = data.TravelerFirstName
				pr.LastName = data.TravelerLastName
				pr.Status = "Fila"

				pr.Active = true
				pr.OriginData = "BigQuery Integration"
				pr.EmissionDate = data.EmissionDate

				if pr.CreatedAt.IsZero() {
					nowUTC := time.Now().UTC()
					dh, err := flytura.DiffHours(nowUTC, flytura.Fuso1, flytura.Fuso2)
					if err != nil {
						panic(err)
					}
					pr.CreatedAt = nowUTC.Add(-time.Duration(dh) * time.Hour)
				}

				if flytura.Normalize(data.OriginAirline) != "aeromexico" {

					pr.CompanyCode = codAirline
					pr.CompanyName = nameAirline
					pr.DirectionOfDestination = "GO"
					pr.Key = data.OriginLocator
					insertPurchardRecordByBigQuery(data.OriginLocator, pr)

					if data.OriginLocator != data.ReturnLocator {
						codAirlineReturn, nameAirlineReturn := airLine.SearchAirlineByName(airlines, data.ReturnAirline)
						pr.CompanyCode = codAirlineReturn
						pr.CompanyName = nameAirlineReturn
						if data.ReturnLocator != "" {
							pr.DirectionOfDestination = "BACK"
							pr.Key = data.ReturnLocator
							insertPurchardRecordByBigQuery(data.ReturnLocator, pr)
						}
					}

				} else {

					// fmt.Println("Airline ", data.OriginAirline, data.OriginETicket)
					pr.DirectionOfDestination = "GO"
					pr.Key = data.OriginETicket
					pr.CompanyCode = codAirline
					pr.CompanyName = nameAirline
					insertPurchardRecordByBigQuery(data.OriginETicket, pr)

					if data.OriginETicket != data.ReturnETicket {
						if data.ReturnETicket != "" {
							codAirlineReturn, nameAirlineReturn := airLine.SearchAirlineByName(airlines, data.ReturnAirline)
							pr.CompanyCode = codAirlineReturn
							pr.CompanyName = nameAirlineReturn
							pr.Key = data.ReturnETicket
							pr.DirectionOfDestination = "BACK"
							insertPurchardRecordByBigQuery(data.ReturnETicket, pr)
						}
					}
				}
			}

		}

		exist, erro := VeryExistKey(
			db.MongoClient,
			flytura.DBName,
			flytura.ConciliationTableName,
			data.OriginLocator,
			data.ReturnLocator,
			data.OriginETicket,
			data.ReturnETicket)

		if erro == nil && !exist {
			InsertOnflyConciliation(db.MongoClient, flytura.DBName, flytura.ConciliationTableName, data)
		}

		cont++

	}

	// fmt.Println(" ")
	// fmt.Println(" ")
	// fmt.Println(" ")
	// fmt.Println("Total registros ", cont)
	// fmt.Println(" ")
	// fmt.Println(" ")
	// fmt.Println(" ")

}

/*
Função criada por Ricardo Silva Ferreira
Inicio da criação 10/05/2026 14:56
Data Final da criação 10/05/2026 14:57
*/
func VeryExistKey(
	client *mongo.Client,
	dbName, collectionName,
	originLocator string,
	returnLocator string,
	originETicket string,
	returnETicket string,
) (bool, error) {

	collection := client.Database(dbName).Collection(collectionName)
	filter := bson.M{
		"originLocator": originLocator,
		"returnLocator": returnLocator,
		"originETicket": originETicket,
		"returnETicket": returnETicket,
	}
	// Variável para armazenar o usuário retornado
	var excelData models.Conciliation
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

/*
Função criada por Ricardo Silva Ferreira
Inicio da criação 10/05/2026 14:54
Data Final da criação 10/05/2026 14:55
*/
// Função para inserir um usuário na coleção "user"
func InsertOnflyConciliation(client *mongo.Client, dbName, collectionName string, data models.Conciliation) error {

	if client == nil {
		fmt.Println("mongo client é nil (não inicializado)")
	}
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
Inicio da criação 11/05/2026 09:27
Data Final da criação 11/05/2026 09:35
*/
func SearchConciliationPagination(
	client *mongo.Client,
	dbName, collectionName string,
	originLocator string,
	returnLocator string,
	originETicket string,
	returnETicket string,
	startDate *time.Time,
	endDate *time.Time,
	page,
	limit int) ([]any, int64, error) {

	collection := client.Database(dbName).Collection(collectionName)

	// Criando o filtro dinâmico
	filter := bson.M{}

	if originLocator != "" {
		filter["originLocator"] = bson.M{"$regex": originLocator, "$options": "i"}
	}

	if returnLocator != "" {
		filter["returnLocator"] = bson.M{"$regex": returnLocator, "$options": "i"}
	}

	// fmt.Println("Origem e ticket ", originETicket)
	if originETicket != "" {
		filter["originETicket"] = bson.M{"$regex": originETicket, "$options": "i"}
	}

	if returnETicket != "" {
		filter["returnETicket"] = bson.M{"$regex": returnETicket, "$options": "i"}
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

		filter["emissionDate"] = dateFilter
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
			SetSort(bson.D{{Key: "emissionDate", Value: -1}}),
	)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(context.Background())

	// Processa os resultados
	var excelData []any
	for cursor.Next(context.Background()) {
		var data models.Conciliation
		if err := cursor.Decode(&data); err != nil {
			return nil, 0, fmt.Errorf("erro ao decodificar Conciliation: %v", err)
		}

		excelData = append(excelData, map[string]any{
			"ID":                         data.ID.Hex(), // Convertendo para string
			"Protocol":                   data.Protocol,
			"EmissionDate":               data.EmissionDate,
			"OriginDate":                 data.OriginDate,
			"ReturnDate":                 data.ReturnDate,
			"OriginLocator":              data.OriginLocator,
			"ReturnLocator":              data.ReturnLocator,
			"OriginETicket":              data.OriginETicket,
			"ReturnETicket":              data.ReturnETicket,
			"OriginAirline":              data.OriginAirline,
			"ReturnAirline":              data.ReturnAirline,
			"TravelerName":               data.TravelerName,
			"TravelerFirstName":          data.TravelerFirstName,
			"TravelerLastName":           data.TravelerLastName,
			"OriginCancelledReason":      data.OriginCancelledReason,
			"ReturnCancelledReason":      data.ReturnCancelledReason,
			"CurrencyCode":               data.CurrencyCode,
			"AmountOrigin":               data.AmountOrigin,
			"AmountReturn":               data.AmountReturn,
			"Active":                     data.Active,
			"CreatedAtContractedCountry": data.CreatedAtContractedCountry,
			"CreatedAtLocalCountry":      data.CreatedAtLocalCountry,
			"BookingStatus":              data.BookingStatus,
			"BookingStatusOld":           data.BookingStatusOld,
			"CancelledReason":            data.CancelledReason,
		})
	}

	// Retorna usuários e total de registros
	return excelData, total, nil
}

/*
Função criada por Ricardo Silva Ferreira
Inicio da criação 11/05/2026 09:57
Data Final da criação 11/05/2026 09:35
*/
func SearchConciliationExcel(
	client *mongo.Client,
	dbName, collectionName string,
	originLocator string,
	returnLocator string,
	originETicket string,
	returnETicket string,
	startDate *time.Time,
	endDate *time.Time) ([]any, int64, error) {

	collection := client.Database(dbName).Collection(collectionName)

	// Criando o filtro dinâmico
	filter := bson.M{}

	if originLocator != "" {
		filter["originLocator"] = bson.M{"$regex": originLocator, "$options": "i"}
	}

	if returnLocator != "" {
		filter["returnLocator"] = bson.M{"$regex": returnLocator, "$options": "i"}
	}

	if originETicket != "" {
		filter["originETicket"] = bson.M{"$regex": originETicket, "$options": "i"}
	}

	if returnETicket != "" {
		filter["returnETicket"] = bson.M{"$regex": returnETicket, "$options": "i"}
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

		filter["emissionDate"] = dateFilter
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
			SetSort(bson.D{{Key: "emissionDate", Value: -1}}),
	)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(context.Background())

	// Processa os resultados
	var excelData []any
	for cursor.Next(context.Background()) {
		var data models.Conciliation
		if err := cursor.Decode(&data); err != nil {
			return nil, 0, fmt.Errorf("erro ao decodificar conciliation: %v", err)
		}

		excelData = append(excelData, map[string]any{
			"ID":                         data.ID.Hex(), // Convertendo para string
			"Protocol":                   data.Protocol,
			"EmissionDate":               data.EmissionDate,
			"OriginDate":                 data.OriginDate,
			"ReturnDate":                 data.ReturnDate,
			"OriginLocator":              data.OriginLocator,
			"ReturnLocator":              data.ReturnLocator,
			"OriginETicket":              data.OriginETicket,
			"ReturnETicket":              data.ReturnETicket,
			"OriginAirline":              data.OriginAirline,
			"ReturnAirline":              data.ReturnAirline,
			"TravelerName":               data.TravelerName,
			"TravelerFirstName":          data.TravelerFirstName,
			"TravelerLastName":           data.TravelerLastName,
			"OriginCancelledReason":      data.OriginCancelledReason,
			"ReturnCancelledReason":      data.ReturnCancelledReason,
			"CurrencyCode":               data.CurrencyCode,
			"AmountOrigin":               data.AmountOrigin,
			"AmountReturn":               data.AmountReturn,
			"Active":                     data.Active,
			"CreatedAtContractedCountry": data.CreatedAtContractedCountry,
			"CreatedAtLocalCountry":      data.CreatedAtLocalCountry,
			"BookingStatus":              data.BookingStatus,
			"BookingStatusOld":           data.BookingStatusOld,
		})
	}

	// Retorna usuários e total de registros
	return excelData, total, nil
}
