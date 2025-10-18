package awsS3

import (
	"Flytura/internal/db"
	"Flytura/internal/models"
	"context"
	"fmt"
	"io"
	"log"
	"time"

	flytura "Flytura"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	bucketName = flytura.BucketProd
	region     = flytura.BucketRegion // ajuste conforme necessário
)

func UploadToS3(file io.Reader, filename string) error {
	accessKey := flytura.AKA
	secretKey := flytura.SKA
	region := region
	bucketName := bucketName
	key := "uploads/" + filename

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(region),
		config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(accessKey, secretKey, ""),
		),
	)
	if err != nil {
		return fmt.Errorf("erro ao carregar config: %w", err)
	}

	client := s3.NewFromConfig(cfg)

	_, err = client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: &bucketName,
		Key:    &key,
		Body:   file,
		// ACL removido porque o bucket não permite ACLs
	})
	if err != nil {
		fmt.Println("Erro detalhado ao enviar para S3:", err)
		return fmt.Errorf("erro ao enviar para S3: %w", err)
	}

	clientDb, errConnectDB1 := db.ConnectMongoDB(flytura.ConectionString)
	if errConnectDB1 != nil {
		log.Println("Erro ao obter nome do arquivo:", err)
		return fmt.Errorf("erro ao enviar para S3: %w", err)

	}
	defer db.CloseMongoDB(clientDb)

	image := models.ImagesDB{
		ID:          primitive.NewObjectID(),
		FileName:    "",
		DtImport:    time.Now(),
		CompanyCode: "",
		CompanyName: "",
		FileURL:     "", // ou uma URL pública se estiver usando S3, etc.
	}

	InsertIMGS3(clientDb, flytura.DBName, "ImagesDB", image)

	fmt.Printf("Arquivo enviado para: https://%s.s3.%s.amazonaws.com/%s\n", bucketName, region, key)
	return nil
}

/*
Função criada por Ricardo Silva Ferreira
Inicio da criação 17/10/2025 17:26
Data Final da criação : 17/10/2025 17:26
*/
// Função para inserir um usuário na coleção "user"
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
