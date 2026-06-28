package awsS3

import (
	"Flytura/internal/db"
	"Flytura/internal/models"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"path"
	"strings"
	"time"

	flytura "Flytura"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	bucketName = flytura.ReturnBucketEnvironment()
	region     = flytura.BucketRegion // ajuste conforme necessário
)

// func UploadToS3(file io.Reader, filename, companyCode, key string) error {
// 	accessKey := flytura.AKA
// 	secretKey := flytura.SKA

// 	region := region
// 	bucketName := bucketName
// 	directory := flytura.ImagesInvoices + "/" + filename

// 	cfg, errLdc := config.LoadDefaultConfig(context.TODO(),
// 		config.WithRegion(region),
// 		config.WithCredentialsProvider(
// 			credentials.NewStaticCredentialsProvider(accessKey, secretKey, ""),
// 		),
// 	)
// 	if errLdc != nil {
// 		return fmt.Errorf("erro ao carregar config: %w", errLdc)
// 	}

// 	client := s3.NewFromConfig(cfg)

// 	_, errLdc = client.PutObject(context.TODO(), &s3.PutObjectInput{
// 		Bucket: &bucketName,
// 		Key:    &directory,
// 		Body:   file,
// 		// ACL removido porque o bucket não permite ACLs
// 	})
// 	if errLdc != nil {
// 		fmt.Println("Erro detalhado ao enviar para S3:", errLdc)
// 		return fmt.Errorf("erro ao enviar para S3: %w", errLdc)
// 	}

// // 	clientDb, errConnectDB1 := db.ConnectMongoDB(flytura.ConectionString)
// // 	if errConnectDB1 != nil {
// // 		log.Println("Erro ao obter nome do arquivo:", errConnectDB1)
// // 		return fmt.Errorf("erro ao enviar para S3: %w", errConnectDB1)

// // 	}
// // 	defer db.CloseMongoDB(clientDb)

// 	airLineData, errAirLineName := airLine.GetAirLineFileName(clientDb, flytura.DBName, "airline", companyCode)
// 	if errAirLineName != nil {
// 		log.Println("Erro ao obter nome do arquivo:", errAirLineName)
// 		return fmt.Errorf("Erro ao pesquisa compania aérea: %w", errAirLineName)
// 	}

// 	fmt.Println("ddd", airLineData["FileName"])

// 	companyName := airLineData["Name"].(string)

// 	image := models.ImagesDB{
// 		ID:           primitive.NewObjectID(),
// 		FileName:     filename,
// 		DtImport:     time.Now().In(flytura.FusoMexic),
// 		CompanyCode:  companyCode,
// 		CompanyName:  companyName,
// 		DownloadDone: false,
// 		Key:          key,
// 		FileURL:      flytura.FileAwsS3URL + "/" + flytura.ImagesInvoices + "/" + filename, // ou uma URL pública se estiver usando S3, etc.
// 	}

// 	InsertIMGS3(clientDb, flytura.DBName, "imagesDB", image)

// 	// fmt.Println("bucketName ", bucketName)
// 	// fmt.Println("region ", region)
// 	// fmt.Println("key ", key)

// 	fmt.Printf("Arquivo enviado para: https://%s.s3.%s.amazonaws.com/%s\n", bucketName, region, key)
// 	return nil
// }

/*
Função criada por Ricardo Silva Ferreira
Inicio da criação 23/10/2025 00:50
Data Final da criação : 23/10/2025 01:15
Essa função apenas faz upload para o S3
*/
func UploadToS3Only(file io.Reader, filename string) error {
	// --- Configurações ---
	accessKey := flytura.AKA
	secretKey := flytura.SKA
	region := region
	bucketName := bucketName
	objectKey := path.Join(flytura.ImagesInvoices, filename)

	// --- Lê o conteúdo do arquivo ---
	var buf bytes.Buffer
	size, err := io.Copy(&buf, file)
	if err != nil {
		return fmt.Errorf("erro ao ler conteúdo do arquivo: %w", err)
	}

	// --- Detecta o tipo MIME ---
	header := buf.Bytes()
	sniffLen := min(len(header), 512)
	contentType := http.DetectContentType(header[:sniffLen])

	// Força tipo PDF se o nome do arquivo terminar com .pdf
	if strings.HasSuffix(strings.ToLower(filename), ".pdf") {
		contentType = "application/pdf"
	}

	// --- Cria configuração AWS ---
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(region),
		config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(accessKey, secretKey, ""),
		),
	)
	if err != nil {
		return fmt.Errorf("erro ao carregar configuração da AWS: %w", err)
	}

	client := s3.NewFromConfig(cfg)

	// --- Envia para o S3 ---
	_, err = client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:             aws.String(bucketName),
		Key:                aws.String(objectKey),
		Body:               bytes.NewReader(buf.Bytes()),
		ContentLength:      aws.Int64(size),
		ContentType:        aws.String(contentType),
		ContentDisposition: aws.String(fmt.Sprintf("inline; filename=\"%s\"", filename)),
	})
	if err != nil {
		return fmt.Errorf("erro ao enviar para o S3: %w", err)
	}

	return nil
}

/*
Função criada por Ricardo Silva Ferreira
Inicio da criação 17/10/2025 17:26
Data Final da criação : 17/10/2025 17:26
*/
// Função para inserir imagem
func InsertIMGS3(client *mongo.Client, dbName, collectionName string, data models.ImagesDB) error {
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
Inicio da criação 19/10/2025 19:07
Data Final da criação : 19/10/2025 19:09
Obs: 26/03/2026 22:50 adicionado verificação para retornar vazio caso não tenha nome do arquivo xml
Obs: 26/03/2026 23:12  adicionado ordenação descrecente por data import
*/
func SearchImagesDBPagination(
	client *mongo.Client,
	dbName,
	collectionName string,
	companyCode *string,
	key *string,
	billedFlytura *bool,
	doDonwload *bool,
	startDate *time.Time,
	endDate *time.Time,
	page int64,
	limit int64) ([]any, int64, error) {

	collection := client.Database(dbName).Collection(collectionName)

	// Criando o filtro dinâmico
	filter := bson.M{}

	if companyCode != nil && *companyCode != "" {
		filter["companyCode"] = *companyCode
	}

	if key != nil && *key != "" {
		filter["key"] = *key
	}

	if billedFlytura != nil {
		// fmt.Println("billedFlytura:", *billedFlytura)
		filter["billedFlytura"] = *billedFlytura
	}

	if doDonwload != nil {

		filter["downloadDone"] = *doDonwload
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

		filter["dtImport"] = dateFilter
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
			SetSort(bson.D{{Key: "dtImport", Value: -1}}),
	)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(context.Background())

	// Processa os resultados
	var dataImg []any
	for cursor.Next(context.Background()) {
		var data models.ImagesDB
		if err := cursor.Decode(&data); err != nil {
			return nil, 0, fmt.Errorf("erro ao decodificar usuário: %v", err)
		}

		xmlFileName := ""
		if data.XMLFileName != "" {
			xmlFileName = flytura.FileAwsS3URL + "/" + flytura.ImagesInvoices + "/" + data.XMLFileName
		}

		dataImg = append(dataImg, map[string]any{
			"ID":              data.ID.Hex(), // Convertendo para string
			"FileName":        data.FileName,
			"CompanyCode":     data.CompanyCode,
			"CompanyName":     data.CompanyName,
			"DtImport":        data.DtImport,
			"ZipFileName":     flytura.FileAwsS3URL + "/" + flytura.ImagesInvoices + "/" + data.ZipFileName,
			"PDFFileName":     flytura.FileAwsS3URL + "/" + flytura.ImagesInvoices + "/" + data.PDFFileName,
			"DownloadPDFDone": data.DownloadPDFDone,
			"XMLFileName":     xmlFileName,
			"DownloadXMLDone": data.DownloadXMLDone,
			"Active":          data.Active,
			"Key":             data.Key,
			"DownloadDone":    data.DownloadDone,
			"BilledFlytura":   data.BilledFlytura,
			"UserNameImport":  data.UserNameImport,
			"OriginData":      data.OriginData,
		})
	}

	// Retorna usuários e total de registros
	return dataImg, total, nil
}

/*
Função criada por Ricardo Silva Ferreira
Inicio da criação 20/10/2025 13:29
Data Final da criação : 20/10/2025 13:30
*/

func SearchImagesDBFull(
	client *mongo.Client,
	dbName, collectionName string,
	companyCode *string,
	startDate *time.Time,
	endDate *time.Time) ([]any, int64, error) {

	collection := client.Database(dbName).Collection(collectionName)

	// Criando o filtro dinâmico
	filter := bson.M{}

	if companyCode != nil && *companyCode != "" {
		filter["companyCode"] = *companyCode
	}

	if startDate != nil || endDate != nil {
		dateFilter := bson.M{}

		if startDate != nil {
			start := time.Date(startDate.Year(), startDate.Month(), startDate.Day(), 0, 0, 0, 0, startDate.Location())
			dateFilter["$gte"] = start
		}

		if endDate != nil {
			end := time.Date(endDate.Year(), endDate.Month(), endDate.Day(), 23, 59, 59, int(time.Second-time.Nanosecond), endDate.Location())
			dateFilter["$lte"] = end
		}

		filter["dtImport"] = dateFilter
	}

	// Contar total de documentos
	total, err := collection.CountDocuments(context.Background(), filter)
	if err != nil {
		return nil, 0, err
	}

	// Executa a consulta sem paginação
	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(context.Background())

	// Processa os resultados
	var dataImg []any
	for cursor.Next(context.Background()) {
		var data models.ImagesDB
		if err := cursor.Decode(&data); err != nil {
			return nil, 0, fmt.Errorf("erro ao decodificar imagem: %v", err)
		}

		dataImg = append(dataImg, map[string]any{
			"ID":              data.ID.Hex(),
			"FileName":        data.FileName,
			"CompanyCode":     data.CompanyCode,
			"CompanyName":     data.CompanyName,
			"DtImport":        data.DtImport,
			"ZipFileName":     flytura.FileAwsS3URL + "/" + flytura.ImagesInvoices + "/" + data.ZipFileName,
			"PDFFileName":     flytura.FileAwsS3URL + "/" + flytura.ImagesInvoices + "/" + data.PDFFileName,
			"DownloadPDFDone": data.DownloadPDFDone,
			"XMLFileName":     flytura.FileAwsS3URL + "/" + flytura.ImagesInvoices + "/" + data.XMLFileName,
			"DownloadXMLDone": data.DownloadXMLDone,
			"Active":          data.Active,
			"Key":             data.Key,
			"DownloadDone":    data.DownloadDone,
			"UserNameImport":  data.UserNameImport,
			"OriginData":      data.OriginData,
		})
	}

	return dataImg, total, nil
}

/*
Função criada por Ricardo Silva Ferreira
Inicio da criação 25/11/2025 16:15
Data Final da criação : 25/11/2025 16:15
*/
func DeleteImagesDBByID(client *mongo.Client, dbName, collectionName, id string) error {
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
Inicio da criação 25/11/2025 17:15
Data Final da criação : 25/11/2025 17:20
*/
func DeleteFromS3(filename string) error {
	accessKey := flytura.AKA
	secretKey := flytura.SKA

	// ✅ limpeza do filename
	filename = strings.TrimSpace(filename)
	filename = strings.ReplaceAll(filename, "\\", "/")
	filename = path.Base(filename)

	if filename == "" {
		return nil
	}

	region := region
	bucketName := bucketName
	objectKey := path.Join(flytura.ImagesInvoices, filename)

	// fmt.Println("Deleting S3 key:", objectKey)

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(region),
		config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(accessKey, secretKey, ""),
		),
	)
	if err != nil {
		return fmt.Errorf("erro ao carregar configuração da AWS: %w", err)
	}

	client := s3.NewFromConfig(cfg)

	// ✅ delete
	_, err = client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
	})
	if err != nil {
		return fmt.Errorf("erro ao apagar do S3: %w", err)
	}

	// ✅ confirmação
	for i := 0; i < 5; i++ {

		_, err = client.HeadObject(context.TODO(), &s3.HeadObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(objectKey),
		})

		if err != nil {
			var notFound *types.NotFound
			if errors.As(err, &notFound) {
				return nil // ✅ já apagado
			}

			// ⚠️ importante: se erro não for NotFound
			return fmt.Errorf("erro ao verificar remoção: %w", err)
		}

		time.Sleep(2 * time.Second)
	}

	return fmt.Errorf("não foi possível confirmar remoção do objeto")
}

/*
Função criada por Ricardo Silva Ferreira
Inicio da criação 11/06/2026 12:53
Data Final da criação : 11/06/2026 13:00
OBS ve na lista de arquivos enviados quais os nomes de pdf's existem na base
*/

func checkExistingPDFs(client *mongo.Client, dbName, collectionName string, files []string) ([]string, error) {
	collection := client.Database(dbName).Collection(collectionName)
	ctx := context.Background()

	// 1. Criar filtro para procurar qualquer arquivo
	var orConditions []bson.M
	for _, f := range files {
		orConditions = append(orConditions, bson.M{
			"pdfFileName": bson.M{
				"$regex": f,
			},
		})
	}

	filter := bson.M{
		"$or": orConditions,
	}

	// 2. Buscar no banco
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("erro ao consultar: %v", err)
	}
	defer cursor.Close(ctx)

	var results []bson.M
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	// 3. Map para evitar repetição
	duplicatesMap := make(map[string]bool)

	for _, doc := range results {
		if val, ok := doc["pdfFileName"].(string); ok {
			for _, f := range files {
				if strings.Contains(val, f) {
					duplicatesMap[f] = true
				}
			}
		}
	}

	// 4. Converter para slice final
	var duplicates []string
	for k := range duplicatesMap {
		duplicates = append(duplicates, k)
	}

	return duplicates, nil
}

/*
Função criada por Ricardo Silva Ferreira
Inicio da criação 12/06/2026 14:13
Data Final da criação : 12/06/2026 14:13
*/
func CountLast30DaysByDtImport(
	client *mongo.Client,
	dbName, collectionName string,
) (int64, error) {

	collection := db.GetCollection(client, dbName, collectionName)

	// 🔹 Define o intervalo: hoje - 30 dias
	endDate := time.Now()
	startDate := endDate.AddDate(0, 0, -30)

	// 🔹 Pipeline
	pipeline := mongo.Pipeline{
		{
			{Key: "$match", Value: bson.M{
				"dtImport": bson.M{
					"$gte": startDate,
					"$lt":  endDate,
				},
			}},
		},
		{
			{Key: "$count", Value: "total"},
		},
	}

	cursor, err := collection.Aggregate(context.Background(), pipeline)
	if err != nil {
		return 0, fmt.Errorf("erro ao agregar: %v", err)
	}
	defer cursor.Close(context.Background())

	// Estrutura para receber resultado
	var result []struct {
		Total int64 `bson:"total"`
	}

	if err := cursor.All(context.Background(), &result); err != nil {
		return 0, fmt.Errorf("erro ao ler resultado: %v", err)
	}

	// 🔹 Se não houver registros
	if len(result) == 0 {
		return 0, nil
	}

	return result[0].Total, nil
}

/*
Função criada por Ricardo Silva Ferreira
Inicio da criação 12/06/2026 14:30
Data Final da criação : 12/06/2026 14:36
*/func CountLast30DaysByAmountRange(
	client *mongo.Client,
	dbName, collectionName string,
) (bson.M, error) {

	collection := db.GetCollection(client, dbName, collectionName)

	endDate := time.Now()
	startDate := endDate.AddDate(0, 0, -30)

	pipeline := mongo.Pipeline{

		// ✅ Filtro
		{
			{Key: "$match", Value: bson.M{
				"createdAtContractedCountry": bson.M{
					"$gte": startDate,
					"$lt":  endDate,
				},
			}},
		},

		// ✅ Null safety
		{
			{Key: "$addFields", Value: bson.M{
				"amountOrigin": bson.M{
					"$ifNull": []interface{}{"$amountOrigin", 0},
				},
			}},
		},

		// ✅ Agrupamento por faixa
		{
			{Key: "$group", Value: bson.M{
				"_id": bson.M{
					"$switch": bson.M{
						"branches": []bson.M{
							{"case": bson.M{"$lte": []interface{}{"$amountOrigin", 1999}}, "then": "até 1999€"},
							{"case": bson.M{
								"$and": []interface{}{
									bson.M{"$gte": []interface{}{"$amountOrigin", 2000}},
									bson.M{"$lte": []interface{}{"$amountOrigin", 5000}},
								},
							}, "then": "2000€ - 5000€"},
							{"case": bson.M{"$gt": []interface{}{"$amountOrigin", 5000}}, "then": "> 5000€"},
						},
						"default": "Sem valor",
					},
				},
				"total": bson.M{"$sum": 1},
				"sum":   bson.M{"$sum": "$amountOrigin"},
			}},
		},

		// ✅ Formato base
		{
			{Key: "$project", Value: bson.M{
				"_id":   0,
				"range": "$_id",
				"total": 1,
				"sum":   1,
			}},
		},

		// ✅ ADD ORDER (IMPORTANTE)
		{
			{Key: "$addFields", Value: bson.M{
				"order": bson.M{
					"$switch": bson.M{
						"branches": []bson.M{
							{"case": bson.M{"$eq": []interface{}{"$range", "até 1999€"}}, "then": 1},
							{"case": bson.M{"$eq": []interface{}{"$range", "2000€ - 5000€"}}, "then": 2},
							{"case": bson.M{"$eq": []interface{}{"$range", "> 5000€"}}, "then": 3},
							{"case": bson.M{"$eq": []interface{}{"$range", "Sem valor"}}, "then": 4},
						},
						"default": 99,
					},
				},
			}},
		},

		// ✅ Ordena antes de agrupar
		{
			{Key: "$sort", Value: bson.M{
				"order": 1,
			}},
		},

		// ✅ Agrupamento final (AGORA levando order junto)
		{
			{Key: "$group", Value: bson.M{
				"_id": nil,
				"data": bson.M{
					"$push": bson.M{
						"range": "$range",
						"total": "$total",
						"sum":   "$sum",
						"order": "$order", // ✅ aqui!!
					},
				},
				"totalGeral": bson.M{"$sum": "$total"},
				"sumGeral":   bson.M{"$sum": "$sum"},
			}},
		},

		{
			{Key: "$project", Value: bson.M{
				"_id":        0,
				"data":       1,
				"totalGeral": 1,
				"sumGeral":   1,
			}},
		},
	}

	cursor, err := collection.Aggregate(context.Background(), pipeline)
	if err != nil {
		return nil, fmt.Errorf("erro ao agregar: %v", err)
	}
	defer cursor.Close(context.Background())

	var result []bson.M
	if err := cursor.All(context.Background(), &result); err != nil {
		return nil, fmt.Errorf("erro ao ler resultado: %v", err)
	}

	// ✅ fallback
	if len(result) == 0 {
		return bson.M{
			"data": []bson.M{
				{"range": "até 1999€", "total": 0, "sum": 0, "order": 1},
				{"range": "2000€ - 5000€", "total": 0, "sum": 0, "order": 2},
				{"range": "> 5000€", "total": 0, "sum": 0, "order": 3},
				{"range": "Sem valor", "total": 0, "sum": 0, "order": 4},
			},
			"totalGeral": 0,
			"sumGeral":   0,
		}, nil
	}

	return result[0], nil
}
