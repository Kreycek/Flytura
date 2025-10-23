package awsS3

import (
	"Flytura/internal/airLine"
	"Flytura/internal/db"
	"Flytura/internal/models"
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"path"
	"strings"
	"time"

	flytura "Flytura"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	bucketName = flytura.BucketProd
	region     = flytura.BucketRegion // ajuste conforme necessário
)

func UploadToS3(file io.Reader, filename, companyCode, key string) error {
	accessKey := flytura.AKA
	secretKey := flytura.SKA

	region := region
	bucketName := bucketName
	directory := flytura.ImagesInvoices + "/" + filename

	cfg, errLdc := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(region),
		config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(accessKey, secretKey, ""),
		),
	)
	if errLdc != nil {
		return fmt.Errorf("erro ao carregar config: %w", errLdc)
	}

	client := s3.NewFromConfig(cfg)

	_, errLdc = client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: &bucketName,
		Key:    &directory,
		Body:   file,
		// ACL removido porque o bucket não permite ACLs
	})
	if errLdc != nil {
		fmt.Println("Erro detalhado ao enviar para S3:", errLdc)
		return fmt.Errorf("erro ao enviar para S3: %w", errLdc)
	}

	clientDb, errConnectDB1 := db.ConnectMongoDB(flytura.ConectionString)
	if errConnectDB1 != nil {
		log.Println("Erro ao obter nome do arquivo:", errConnectDB1)
		return fmt.Errorf("erro ao enviar para S3: %w", errConnectDB1)

	}
	defer db.CloseMongoDB(clientDb)

	airLineData, errAirLineName := airLine.GetAirLineFileName(clientDb, flytura.DBName, "airline", companyCode)
	if errAirLineName != nil {
		log.Println("Erro ao obter nome do arquivo:", errAirLineName)
		return fmt.Errorf("Erro ao pesquisa compania aérea: %w", errAirLineName)
	}

	fmt.Println("ddd", airLineData["FileName"])

	companyName := airLineData["Name"].(string)

	image := models.ImagesDB{
		ID:           primitive.NewObjectID(),
		FileName:     filename,
		DtImport:     time.Now(),
		CompanyCode:  companyCode,
		CompanyName:  companyName,
		DownloadDone: false,
		Key:          key,
		FileURL:      flytura.FileAwsS3URL + "/" + flytura.ImagesInvoices + "/" + filename, // ou uma URL pública se estiver usando S3, etc.
	}

	InsertIMGS3(clientDb, flytura.DBName, "imagesDB", image)

	// fmt.Println("bucketName ", bucketName)
	// fmt.Println("region ", region)
	// fmt.Println("key ", key)

	fmt.Printf("Arquivo enviado para: https://%s.s3.%s.amazonaws.com/%s\n", bucketName, region, key)
	return nil
}

func UploadToS3Only(file io.Reader, filename, companyCode, key string) error {
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
*/
func SearchImagesDBPagination(
	client *mongo.Client,
	dbName, collectionName string,
	companyCode *string,
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
		options.Find().SetSkip(int64((page-1)*limit)).SetLimit(int64(limit)),
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

		dataImg = append(dataImg, map[string]any{
			"ID":           data.ID.Hex(), // Convertendo para string
			"FileName":     data.FileName,
			"CompanyCode":  data.CompanyCode,
			"CompanyName":  data.CompanyName,
			"DtImport":     data.DtImport,
			"FileUrl":      data.FileURL,
			"Active":       data.Active,
			"Key":          data.Key,
			"DownloadDone": data.DownloadDone,
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
			"ID":           data.ID.Hex(),
			"FileName":     data.FileName,
			"CompanyCode":  data.CompanyCode,
			"CompanyName":  data.CompanyName,
			"DtImport":     data.DtImport,
			"FileUrl":      data.FileURL,
			"Active":       data.Active,
			"Key":          data.Key,
			"DownloadDone": data.DownloadDone,
		})
	}

	return dataImg, total, nil
}
