// config/config.go
package flytura

import (
	"Flytura/internal/models"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// Variável global que contém a chave secreta para JWT
var SecretKey = []byte("my_secret_key")
var UrlSiteLocalHost = "localhost:4200"
var UrlSiteProduction = "54.156.244.197"
var UrlSiteHomol = "18.210.18.180"
var BucketProd = "flytura-bucket"
var BucketRegion = "us-east-1"
var ImagesInvoices = "invoices"
var AKA = os.Getenv("AKID")
var SKA = os.Getenv("ASK")
var FileAwsS3URL = "https://flytura-bucket.s3.us-east-1.amazonaws.com"
var ConectionString = "mongodb://admin:secret@localhost:27017"
var DBName = "flytura"

var UserDBTableName = "user"
var PurcharseRecordTableName = "purcharseRecord"
var TokenAccessTableName = "tokenAccess"
var StatusImportTableName = "statusImport"
var PerfilTableName = "perfil"
var ImagesDBTableName = "imagesDB"
var Airline = "airline"
var OutPutInvoices = "outPutInvoices"

// var ConectionString = "mongodb://localhost:27017"

func TokenValido(w http.ResponseWriter, r *http.Request) (bool, string) {

	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		// http.Error(w, "Token não fornecido", httjp.StatusUnauthorized)
		return false, "Token não fornecido"
	}

	// O formato esperado do cabeçalho é "Bearer <token>"
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	if tokenString == authHeader {
		// http.Error(w, "Token malformado", http.StatusUnauthorized)
		return false, "Token malformado"
	}

	// Validar o token JWT
	_, err := ValidateToken(tokenString)
	if err != nil {
		// http.Error(w, fmt.Sprintf("Token inválido: %v", err), http.StatusUnauthorized)
		return false, "Token inválido"
	}

	return true, "Token válido"
}

// var secretKey = []byte("my_secret_key") // Chave secreta para validar o JWT

// Função para verificar o token JWT
func ValidateToken(tokenString string) (*jwt.Token, error) {
	// Parsing e validação do token JWT
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Verifique se o método de assinatura é o correto (HMAC)
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("método de assinatura inválido")
		}
		return SecretKey, nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}

func FormataRetornoHTTP(w http.ResponseWriter, obj any, codHttp int) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(codHttp) // Código 200 OK
	return json.NewEncoder(w).Encode(map[string]any{"message": obj})
}

func FormataRetornoHTTPGeneric(w http.ResponseWriter, bodyName string, body any, codHttp int) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(codHttp) // Código 200 OK
	return json.NewEncoder(w).Encode(map[any]any{"users": body})
}

/*
Função criada por Ricardo Silva Ferreira
Inicio da criação 14/10/2025 22:24
Data Final da criação :  14/10/2025 22:32
*/
func VerifyAccessValidTokenListSheet(client *mongo.Client, dbName, collectionName, token string) (map[string]any, error) {

	collection := client.Database(dbName).Collection(collectionName)

	// objectID, erroId := primitive.ObjectIDFromHex(excelId)
	// if erroId != nil {
	// 	log.Fatalf("Erro ao converter string para ObjectID: %v", erroId)
	// }

	filter := bson.M{
		"token":  token,
		"active": true,
	}

	// Variável para armazenar o usuário retornado
	var tokenAccessSheet models.TokenAccessSheet

	// Usar FindOne para pegar apenas um único registro
	err := collection.FindOne(context.Background(), filter).Decode(&tokenAccessSheet)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("token não encontrado")
		}
		return nil, fmt.Errorf("erro ao buscar plano de contas: %v", err)
	}

	// Converter o _id para string

	// Retornar o usuário como um mapa
	_return := map[string]any{
		"ID":   tokenAccessSheet.ID.Hex(), // Agora o campo ID é uma string
		"Name": tokenAccessSheet.Name,
		"Code": tokenAccessSheet.Token,
	}

	return _return, nil
}
