package users

import (
	"Flytura/internal/db"
	"Flytura/internal/models"
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Modelo de Usuário com campos do MongoDB

func GetUserByID(client *mongo.Client, dbName, collectionName, userID string) (map[string]any, error) {
	collection := client.Database(dbName).Collection(collectionName)

	// Criar filtro para buscar um usuário pelo ID

	objectID, erroId := primitive.ObjectIDFromHex(userID)
	if erroId != nil {
		log.Fatalf("Erro ao converter string para ObjectID: %v", erroId)
	}

	filter := bson.M{"_id": objectID}

	// Variável para armazenar o usuário retornado
	var user models.User

	// Usar FindOne para pegar apenas um único registro
	err := collection.FindOne(context.Background(), filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("usuário não encontrado")
		}
		return nil, fmt.Errorf("erro ao buscar usuário: %v", err)
	}

	// Converter o _id para string
	userID = user.ID.Hex()

	// Retornar o usuário como um mapa
	userData := map[string]any{
		"ID":             userID, // Agora o campo ID é uma string
		"Name":           user.Name,
		"LastName":       user.LastName,
		"Email":          user.Email,
		"PassportNumber": user.PassportNumber,
		"Perfil":         user.Perfil,
		"UserName":       user.Username,
		"Active":         user.Active,
		"Mobile":         user.Mobile,
	}

	return userData, nil
}

func SearchUsers(client *mongo.Client, dbName, collectionName string, name, email *string, perfis []int, page, limit int64) ([]any, int64, error) {
	collection := client.Database(dbName).Collection(collectionName)

	// Criando o filtro dinâmico
	filter := bson.M{}
	if name != nil && *name != "" {
		filter["name"] = bson.M{"$regex": *name, "$options": "i"}
	}
	if email != nil && *email != "" {
		filter["email"] = bson.M{"$regex": *email, "$options": "i"}
	}
	if len(perfis) > 0 {
		filter["perfil"] = bson.M{"$in": perfis}
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
	var users []any
	for cursor.Next(context.Background()) {
		var user models.User
		if err := cursor.Decode(&user); err != nil {
			return nil, 0, fmt.Errorf("erro ao decodificar usuário: %v", err)
		}

		users = append(users, map[string]any{
			"ID":             user.ID.Hex(), // Converte ObjectID para string
			"Name":           user.Name,
			"LastName":       user.LastName,
			"Email":          user.Email,
			"PassportNumber": user.PassportNumber,
			"Perfil":         user.Perfil,
			"UserName":       user.Username,
			"Active":         user.Active,
			"Mobile":         user.Mobile,
		})
	}

	// Retorna usuários e total de registros
	return users, total, nil
}

// Função para obter todos os usuários do banco de dados
func GetAllUsers(client *mongo.Client, dbName, collectionName string, page, limit int) ([]any, int, error) {
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

	var users []any
	for cursor.Next(context.Background()) {
		var user models.User
		if err := cursor.Decode(&user); err != nil {
			return nil, 0, fmt.Errorf("erro ao decodificar usuário: %v", err)
		}

		// Adiciona os usuários formatados
		users = append(users, map[string]any{
			"ID":             user.ID.Hex(), // Convertendo para string
			"Name":           user.Name,
			"LastName":       user.LastName,
			"Email":          user.Email,
			"PassportNumber": user.PassportNumber,
			"Perfil":         user.Perfil,
			"UserName":       user.Username,
			"Active":         user.Active,
			"Mobile":         user.Mobile,
		})
	}

	if err := cursor.Err(); err != nil {
		return nil, 0, fmt.Errorf("erro ao iterar no cursor: %v", err)
	}

	return users, int(total), nil
}

// Função para inserir um usuário na coleção "user"
func InsertUser(client *mongo.Client, dbName, collectionName string, user models.User) error {
	collection := client.Database(dbName).Collection(collectionName)

	// Criar um contexto para a operação de inserção
	ctx := context.Background()

	// Inserir o documento
	_, err := collection.InsertOne(ctx, user)
	if err != nil {
		return fmt.Errorf("erro ao inserir usuário: %v", err)
	}

	return nil
}
