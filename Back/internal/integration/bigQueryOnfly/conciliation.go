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
	"strings"
	"time"

	"cloud.google.com/go/bigquery"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/api/iterator"
)

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

/*
Função criada por Ricardo Silva Ferreira
Inicio da criação 09/05/2026 10:07
Data Final da criação 09/05/2026 114:30
*/

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
			return_cancelled_reason,
			flight_origin,
			flight_destination,
			origin_country_code,
			origin_city,
			destination_country_code,
			destination_city,
			origin_airline_commercial,
			return_airline_commercial
        FROM 
			conciliation.gold_flytura 
			--where emission_date='2026-06-01'
			--where (emission_date>='2026-04-01' and emission_date<='2026-05-26') 
			where (emission_date>='2026-04-30' and emission_date<=@today) 
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
		data.Protocol = flytura.ConvertToString(row, "protocol")
		data.EmissionDate = flytura.ConvertToTimeDate(row, "emission_date")
		data.OriginDate = flytura.ConvertToTimeDate(row, "origin_date")
		data.ReturnDate = flytura.ConvertToTimeDate(row, "return_date")

		data.OriginLocator = strings.ReplaceAll(strings.ReplaceAll(flytura.ConvertToString(row, "origin_locator"), "-", ""), " ", "")
		data.ReturnLocator = strings.ReplaceAll(strings.ReplaceAll(flytura.ConvertToString(row, "return_locator"), "-", ""), " ", "")
		data.OriginETicket = strings.ReplaceAll(strings.ReplaceAll(flytura.ConvertToString(row, "origin_e_ticket"), "-", ""), " ", "")
		data.ReturnETicket = strings.ReplaceAll(strings.ReplaceAll(flytura.ConvertToString(row, "return_e_ticket"), "-", ""), " ", "")

		// fmt.Println("Tamanho", len(strings.ReplaceAll("9572288797433 ", " ", "")))
		data.OriginAirline = flytura.ConvertToString(row, "origin_airline")
		data.ReturnAirline = flytura.ConvertToString(row, "return_airline")
		data.TravelerName = flytura.ConvertToString(row, "traveler_name")
		data.TravelerFirstName = flytura.ConvertToString(row, "traveler_first_name")
		data.TravelerLastName = flytura.ConvertToString(row, "traveler_last_name")
		data.AmountOrigin = flytura.ConvertToFloat64(row, "onfly_amount_origin")
		data.AmountReturn = flytura.ConvertToFloat64(row, "onfly_amount_return")
		data.BookingStatus = strings.ReplaceAll(strings.ReplaceAll(flytura.ConvertToString(row, "origin_status"), "-", ""), " ", "")
		data.BookingStatusOld = flytura.ConvertToString(row, "origin_status_old")
		data.CurrencyCode = flytura.ConvertToString(row, "currency_code")
		data.OriginCancelledReason = flytura.ConvertToString(row, "origin_cancelled_reason")
		data.ReturnCancelledReason = flytura.ConvertToString(row, "return_cancelled_reason")
		data.Active = true

		//ESSA PARTE FOI ACRESCENTADA DIA 09/06/2026 16:46
		data.FlightOrigin = flytura.ConvertToString(row, "flight_origin")
		data.FlightDestination = flytura.ConvertToString(row, "flight_destination")
		data.OriginCountryCode = flytura.ConvertToString(row, "origin_country_code")
		data.OriginCity = flytura.ConvertToString(row, "origin_city")
		data.DestinationCountryCode = flytura.ConvertToString(row, "destination_country_code")
		data.DestinationCity = flytura.ConvertToString(row, "destination_city")
		data.OriginAirlineCommercial = flytura.ConvertToString(row, "origin_airline_commercial")
		data.ReturnAirlineCommercial = flytura.ConvertToString(row, "return_airline_commercial")

		nowUTC := time.Now().UTC()
		dh, err := flytura.DiffHours(nowUTC, flytura.Fuso1, flytura.Fuso2)
		if err != nil {
			log.Printf("erro ao calcular DiffHours: %v; usando nowUTC", err)
			data.CreatedAtContractedCountry = nowUTC
		} else {
			data.CreatedAtContractedCountry = nowUTC.Add(-time.Duration(dh) * time.Hour)
		}
		// fmt.Println("FlightOrigin", data.Protocol)
		if strings.ReplaceAll(data.Protocol, " ", "") == "040S2Z" {
			fmt.Println("TravelerName", data.TravelerName)
			fmt.Println("OriginCountryCode ", data.OriginCountryCode)
		}

		data.CreatedAtLocalCountry = nowUTC

		// fmt.Println("Total ", convertToDecimal128(row, "onfly_amount_origin"))

		// fmt.Println("Protocol:", row["protocol"])
		// fmt.Println("Emission Data:", row["emission_date"])
		// fmt.Println("traveler_first_name:", row["traveler_first_name"])
		// fmt.Println("traveler_last_name:", row["traveler_last_name"])
		// fmt.Println("Traveler:", data.TravelerName)

		// if strings.Contains(data.Protocol, "04300M") {
		// 	fmt.Println(" ")
		// 	fmt.Println("TravelerName:", data.TravelerName)
		// 	fmt.Println("OriginLocator:", data.OriginLocator)
		// 	fmt.Println("ReturnLocator:", data.ReturnLocator)
		// 	fmt.Println("OriginETicket:", data.OriginETicket)
		// 	fmt.Println("OriginETicket:", data.ReturnETicket)
		// 	fmt.Println("OriginAirline:", data.OriginAirline)
		// 	fmt.Println("ReturnAirline:", data.ReturnAirline)
		// 	fmt.Println(" ")
		// 	fmt.Println(" ")

		// }

		// fmt.Println("Amount Origin:", row["onfly_amount_origin"])
		// fmt.Println("Amount Return:", row["onfly_amount_return"])
		// fmt.Println("Currency:", row["currency_code"])
		// fmt.Println("Emission Data:", row["emission_date"])
		// fmt.Println("Data:", data)
		// fmt.Println("--------------------------------")

		// fmt.Println("data.OriginAirline ", data.OriginAirline)

		if flytura.Normalize(data.BookingStatus) == "emitted" {

			originCodAirline, _ := airLine.SearchAirlineByName(airlines, data.OriginAirline)
			returnCodAirline, _ := airLine.SearchAirlineByName(airlines, data.ReturnAirline)

			processOriginAirline := (originCodAirline == "0001" || originCodAirline == "0002" || originCodAirline == "0003")
			processReturnCodAirline := (returnCodAirline == "0001" || returnCodAirline == "0002" || returnCodAirline == "0003")

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

			origin := flytura.Normalize(data.OriginAirline)
			ret := flytura.Normalize(data.ReturnAirline)

			isOriginAM := origin == "aeromexico"
			isReturnAM := ret == "aeromexico"
			isOriginEmpty := origin == ""
			isReturnEmpty := ret == ""

			switch {
			case isOriginAM && isReturnAM:
				if data.OriginETicket == data.ReturnETicket {
					insertRegister(data.OriginETicket, pr, data.OriginAirline, airlines, "GO")
				} else {
					insertRegister(data.OriginETicket, pr, data.OriginAirline, airlines, "GO")
					insertRegister(data.ReturnETicket, pr, data.ReturnAirline, airlines, "BACK")
				}

			case isOriginAM && !isReturnAM:
				insertRegister(data.OriginETicket, pr, data.OriginAirline, airlines, "GO")
				if processReturnCodAirline {
					insertRegister(data.OriginLocator, pr, data.ReturnAirline, airlines, "BACK")
				}

			case !isOriginAM && isReturnAM:
				if processOriginAirline {
					insertRegister(data.OriginLocator, pr, data.OriginAirline, airlines, "GO")
				}
				insertRegister(data.ReturnETicket, pr, data.ReturnAirline, airlines, "BACK")

			case !isOriginAM && !isReturnAM:
				if data.OriginLocator == data.ReturnLocator {
					if processOriginAirline {
						insertRegister(data.OriginLocator, pr, data.OriginAirline, airlines, "GO")
					}
				} else {
					if processOriginAirline {
						insertRegister(data.OriginLocator, pr, data.OriginAirline, airlines, "GO")
					}
					if processReturnCodAirline {
						insertRegister(data.ReturnLocator, pr, data.ReturnAirline, airlines, "BACK")
					}
				}

			case isOriginAM && isReturnEmpty:
				if processOriginAirline {
					insertRegister(data.OriginETicket, pr, data.OriginAirline, airlines, "GO")
				}

			case isOriginEmpty && isReturnAM:
				if processReturnCodAirline {
					insertRegister(data.ReturnETicket, pr, data.ReturnAirline, airlines, "BACK")
				}

			case !isOriginAM && isReturnEmpty:
				if processOriginAirline {
					insertRegister(data.OriginLocator, pr, data.OriginAirline, airlines, "GO")
				}

			case isOriginEmpty && !isReturnAM:
				if processReturnCodAirline {
					insertRegister(data.ReturnLocator, pr, data.ReturnAirline, airlines, "BACK")
				}
			}

			// if flytura.Normalize(data.OriginAirline) == "aeromexico" &&
			// 	flytura.Normalize(data.ReturnAirline) == "aeromexico" &&
			// 	data.OriginETicket == data.ReturnETicket {
			// 	insertRegister(data.OriginETicket, pr, data.OriginAirline, airlines, "GO")

			// } else if flytura.Normalize(data.OriginAirline) == "aeromexico" &&
			// 	flytura.Normalize(data.ReturnAirline) == "aeromexico" &&
			// 	data.OriginETicket != data.ReturnETicket {
			// 	insertRegister(data.OriginETicket, pr, data.OriginAirline, airlines, "GO")
			// 	insertRegister(data.ReturnETicket, pr, data.ReturnAirline, airlines, "BACK")

			// } else if flytura.Normalize(data.OriginAirline) == "aeromexico" &&
			// 	flytura.Normalize(data.ReturnAirline) != "aeromexico" {
			// 	insertRegister(data.OriginETicket, pr, data.OriginAirline, airlines, "GO")
			// 	if processReturnCodAirline {
			// 		insertRegister(data.OriginLocator, pr, data.ReturnAirline, airlines, "BACK")
			// 	}

			// } else if flytura.Normalize(data.OriginAirline) != "aeromexico" &&
			// 	flytura.Normalize(data.ReturnAirline) == "aeromexico" {
			// 	if processOriginAirline {
			// 		insertRegister(data.OriginLocator, pr, data.OriginAirline, airlines, "GO")
			// 	}

			// 	insertRegister(data.ReturnETicket, pr, data.ReturnAirline, airlines, "BACK")

			// } else if flytura.Normalize(data.OriginAirline) != "aeromexico" &&
			// 	flytura.Normalize(data.ReturnAirline) != "aeromexico" &&
			// 	data.OriginLocator == data.ReturnLocator {
			// 	if processOriginAirline {
			// 		insertRegister(data.OriginLocator, pr, data.OriginAirline, airlines, "GO")
			// 	}

			// } else if flytura.Normalize(data.OriginAirline) != "aeromexico" &&
			// 	flytura.Normalize(data.ReturnAirline) != "aeromexico" &&
			// 	data.OriginLocator != data.ReturnLocator {
			// 	if processOriginAirline {
			// 		insertRegister(data.OriginLocator, pr, data.OriginAirline, airlines, "GO")
			// 	}
			// 	if processReturnCodAirline {
			// 		insertRegister(data.ReturnLocator, pr, data.ReturnAirline, airlines, "BACK")
			// 	}

			// } else if flytura.Normalize(data.OriginAirline) == "aeromexico" &&
			// 	flytura.Normalize(data.ReturnAirline) == "" {
			// 	if processOriginAirline {
			// 		insertRegister(data.OriginETicket, pr, data.OriginAirline, airlines, "GO")
			// 	}

			// } else if flytura.Normalize(data.OriginAirline) == "" &&
			// 	flytura.Normalize(data.ReturnAirline) == "aeromexico" {
			// 	if processReturnCodAirline {
			// 		insertRegister(data.ReturnETicket, pr, data.ReturnAirline, airlines, "BACK")
			// 	}

			// } else if flytura.Normalize(data.OriginAirline) != "aeromexico" &&
			// 	flytura.Normalize(data.ReturnAirline) == "" {
			// 	if processOriginAirline {
			// 		insertRegister(data.OriginLocator, pr, data.OriginAirline, airlines, "GO")
			// 	}

			// } else if flytura.Normalize(data.OriginAirline) == "" &&
			// 	flytura.Normalize(data.ReturnAirline) != "aeromexico" {
			// 	if processReturnCodAirline {
			// 		insertRegister(data.ReturnLocator, pr, data.ReturnAirline, airlines, "BACK")
			// 	}
			// }

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

func insertRegister(key string, pr models.PurcharseRecord, airLineName string, airLines []any, direction string) {
	codAirlineReturn, nameAirlineReturn := airLine.SearchAirlineByName(airLines, airLineName)
	pr.CompanyCode = codAirlineReturn
	pr.CompanyName = nameAirlineReturn
	pr.Key = key
	pr.DirectionOfDestination = direction

	if pr.CompanyCode == "0001" && pr.Key != "" {
		insertPurchardRecordByBigQuery(key, pr)
	} else if pr.CompanyCode == "0002" && pr.Key != "" && pr.LastName != "" {
		insertPurchardRecordByBigQuery(key, pr)
	} else if pr.CompanyCode == "0003" && pr.Key != "" && pr.Name != "" && pr.LastName != "" {
		insertPurchardRecordByBigQuery(key, pr)
	}
}

/*
Função criada por Ricardo Silva Ferreira
Inicio da criação 10/05/2026 14:56
Data Final da criação 10/05/2026 14:57
*/
func VeryExistKey(
	client *mongo.Client,
	dbName, collectionName string,
	originLocator string,
	returnLocator string,
	originETicket string,
	returnETicket string,
) (bool, error) {

	collection := client.Database(dbName).Collection(collectionName)

	filter := bson.M{
		"$expr": bson.M{
			"$and": []bson.M{
				{
					"$eq": []interface{}{
						bson.M{"$trim": bson.M{"input": "$originLocator"}},
						originLocator,
					},
				},
				{
					"$eq": []interface{}{
						bson.M{"$trim": bson.M{"input": "$returnLocator"}},
						returnLocator,
					},
				},
				{
					"$eq": []interface{}{
						bson.M{"$trim": bson.M{"input": "$originETicket"}},
						originETicket,
					},
				},
				{
					"$eq": []interface{}{
						bson.M{"$trim": bson.M{"input": "$returnETicket"}},
						returnETicket,
					},
				},
			},
		},
	}

	var excelData models.Conciliation

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
			SetSort(bson.D{
				{Key: "emissionDate", Value: -1},
				{Key: "_id", Value: -1},
			}),
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
			"FlightOrigin":               data.FlightOrigin,
			"FlightDestination":          data.FlightDestination,
			"OriginCountryCode":          data.OriginCountryCode,
			"OriginCity":                 data.OriginCity,
			"DestinationCountryCode":     data.DestinationCountryCode,
			"DestinationCity":            data.DestinationCity,
			"OriginAirlineCommercial":    data.OriginAirlineCommercial,
			"ReturnAirlineCommercial":    data.ReturnAirlineCommercial,
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
			"FlightOrigin":               data.FlightOrigin,
			"FlightDestination":          data.FlightDestination,
			"OriginCountryCode":          data.OriginCountryCode,
			"OriginCity":                 data.OriginCity,
			"DestinationCountryCode":     data.DestinationCountryCode,
			"DestinationCity":            data.DestinationCity,
			"OriginAirlineCommercial":    data.OriginAirlineCommercial,
			"ReturnAirlineCommercial":    data.ReturnAirlineCommercial,
		})
	}

	// Retorna usuários e total de registros
	return excelData, total, nil
}
