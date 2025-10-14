package users

import (
	flytura "Flytura"
	"Flytura/internal/db"
	"Flytura/internal/models"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

// Função de handler para a rota GET /users
func GetAllUsersHandler(w http.ResponseWriter, r *http.Request) {
	status, msg := flytura.TokenValido(w, r)
	if !status {
		http.Error(w, fmt.Sprintf("erro ao validar token: %v", msg), http.StatusUnauthorized)
		return
	}

	client, err := db.ConnectMongoDB(flytura.ConectionString)
	if err != nil {
		http.Error(w, fmt.Sprintf("erro ao conectar ao MongoDB: %v", err), http.StatusInternalServerError)
		return
	}
	defer db.CloseMongoDB(client)

	// Obter parâmetros de paginação
	query := r.URL.Query()
	page, err := strconv.Atoi(query.Get("page"))
	if err != nil || page < 1 {
		page = 1 // Padrão: primeira página
	}

	limit, err := strconv.Atoi(query.Get("limit"))
	if err != nil || limit < 1 {
		limit = 10 // Padrão: 10 registros por página
	}

	// Obter usuários paginados
	users, total, err := GetAllUsers(client, flytura.DBName, "user", page, limit)
	if err != nil {
		http.Error(w, fmt.Sprintf("erro ao buscar perfis: %v", err), http.StatusInternalServerError)
		return
	}

	// Criar resposta JSON com paginação
	response := map[string]any{
		"total": total,
		"page":  page,
		"limit": limit,
		"pages": (total + limit - 1) / limit, // Calcula o número total de páginas
		"users": users,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("erro ao codificar resposta JSON: %v", err)
	}
}

func SearchUsersHandler(w http.ResponseWriter, r *http.Request) {
	// Verificar se a requisição é do tipo POST
	if r.Method != http.MethodPost {
		http.Error(w, "Método não permitido, deve ser post", http.StatusMethodNotAllowed)
		return
	}

	// Validar Token
	status, msg := flytura.TokenValido(w, r)
	if !status {
		http.Error(w, fmt.Sprintf("erro ao validar token: %v", msg), http.StatusUnauthorized)
		return
	}

	// Conectar ao MongoDB
	client, err := db.ConnectMongoDB(flytura.ConectionString)
	if err != nil {
		http.Error(w, "Erro ao conectar ao MongoDB", http.StatusInternalServerError)
		return
	}
	defer client.Disconnect(context.Background())

	// Definir estrutura para receber os parâmetros
	var request struct {
		Name   *string `json:"name"`
		Email  *string `json:"email"`
		Perfil []int   `json:"perfil"`
		Page   int64   `json:"page"`
		Limit  int64   `json:"limit"`
	}

	// Decodificar o corpo da requisição JSON
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Erro ao decodificar JSON", http.StatusBadRequest)
		return
	}

	// Definir valores padrão para paginação
	if request.Page < 1 {
		request.Page = 1
	}
	if request.Limit < 1 {
		request.Limit = 10
	}

	// Buscar usuários com paginação
	users, total, err := SearchUsers(client, flytura.DBName, "user", request.Name, request.Email, request.Perfil, request.Page, request.Limit)
	if err != nil {
		http.Error(w, "Erro ao buscar usuários", http.StatusInternalServerError)
		return
	}

	// Criar resposta JSON com paginação
	response := map[string]any{
		"total": total,
		"page":  request.Page,
		"limit": request.Limit,
		"pages": (total + request.Limit - 1) / request.Limit, // Número total de páginas
		"users": users,
	}

	// Retornar resposta JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("erro ao codificar resposta JSON: %v", err)
	}
}

// Variáveis globais

// Função para verificar o nome de usuário e senha
func VerifyExistUser(w http.ResponseWriter, r *http.Request) {

	status, msg := flytura.TokenValido(w, r)
	if !status {
		http.Error(w, fmt.Sprintf("erro ao validar token: %v", msg), http.StatusUnauthorized)
		return
	}

	// Parse o corpo da requisição
	var email models.EmailRequest

	err := json.NewDecoder(r.Body).Decode(&email)
	if err != nil {
		http.Error(w, "Erro ao ler o corpo da requisição", http.StatusBadRequest)
		return
	}

	// fmt.Println("email", email)
	// Conectar ao MongoDB
	client, err := db.ConnectMongoDB(flytura.ConectionString)
	if err != nil {
		http.Error(w, "Erro ao conectar ao banco de dados", http.StatusInternalServerError)
		return
	}
	defer db.CloseMongoDB(client)

	// Obter a coleção de usuários
	collection := db.GetCollection(client, flytura.DBName, "user")
	// filter := bson.D{
	// 	{Key: "$or", Value: bson.A{
	// 		bson.D{{Key: "email", Value: userName}},
	// 	}},
	// }

	filter := bson.D{
		{Key: "email", Value: bson.M{
			"$regex":   email.Email,
			"$options": "i", // Ignora diferença entre maiúsculas e minúsculas
		}},
	}

	var result bson.M
	err = collection.FindOne(context.Background(), filter).Decode(&result)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			flytura.FormataRetornoHTTP(w, false, http.StatusOK)
			return
		}
		http.Error(w, "Erro interno", http.StatusInternalServerError)
		return
	}

	// Se encontrou um documento, retorna true
	flytura.FormataRetornoHTTP(w, true, http.StatusOK)

}

func GetUserByIdHandler(w http.ResponseWriter, r *http.Request) {
	// Validar token
	status, msg := flytura.TokenValido(w, r)
	if !status {
		http.Error(w, fmt.Sprintf("erro ao validar token: %v", msg), http.StatusUnauthorized)
		return
	}

	// Conectar ao MongoDB
	client, err := db.ConnectMongoDB(flytura.ConectionString)
	if err != nil {
		http.Error(w, "Erro ao conectar ao MongoDB", http.StatusInternalServerError)
		return
	}
	defer client.Disconnect(context.Background())

	// Extrair o ID da URL
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "ID não fornecido na URL", http.StatusBadRequest)
		return
	}

	// Verifica se o ID fornecido é válido
	_, err = primitive.ObjectIDFromHex(id)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	// Buscar o usuário no banco de dados pelo ID
	user, err := GetUserByID(client, flytura.DBName, "user", id)
	if err != nil {
		http.Error(w, "Erro ao buscar usuário", http.StatusInternalServerError)
		return
	}

	// Configurar o cabeçalho da resposta como JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Enviar o usuário como resposta JSON
	if err := json.NewEncoder(w).Encode(user); err != nil {
		log.Printf("erro ao codificar resposta JSON: %v", err)
		http.Error(w, "Erro ao codificar resposta", http.StatusInternalServerError)
	}
}

func InsertUserHandler(w http.ResponseWriter, r *http.Request) {
	// Validar o token de autenticação
	status, msg := flytura.TokenValido(w, r)
	if !status {
		http.Error(w, fmt.Sprintf("erro ao validar token: %v", msg), http.StatusUnauthorized)
		return
	}

	// Ler o corpo da requisição
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "erro ao decodificar corpo da requisição", http.StatusBadRequest)
		return
	}

	if user.Active == false {
		user.Active = true
	}

	if user.CreatedAt.IsZero() {
		user.CreatedAt = time.Now()
	}
	// Conectar ao MongoDB
	client, err := db.ConnectMongoDB(flytura.ConectionString)
	if err != nil {
		http.Error(w, fmt.Sprintf("erro ao conectar ao MongoDB: %v", err), http.StatusInternalServerError)
		return
	}
	defer db.CloseMongoDB(client)

	if user.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)

		if err != nil {
			log.Println("Erro ao gerar hash da senha:", err)
			return
		}

		// Atribuindo a senha hashada ao campo Password
		user.Password = string(hashedPassword)
	}

	// Inserir o usuário no MongoDB
	err = InsertUser(client, flytura.DBName, "user", user)
	if err != nil {
		http.Error(w, fmt.Sprintf("erro ao inserir usuário: %v", err), http.StatusInternalServerError)
		return
	}

	// Retornar a resposta com os dados dos usuários
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(user); err != nil {
		log.Printf("erro ao codificar resposta JSON: %v", err)
	}

}

func UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	// Validar o token de autenticação
	status, msg := flytura.TokenValido(w, r)
	if !status {
		http.Error(w, fmt.Sprintf("Erro ao validar token: %v", msg), http.StatusUnauthorized)
		return
	}

	// Decodificar o JSON recebido
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		// http.Error(w, "Erro ao decodificar JSON", http.StatusBadRequest)
		flytura.FormataRetornoHTTP(w, "Erro ao decodificar JSON", http.StatusBadRequest)

		return
	}

	// Verifica se o ID é válido
	if user.ID.IsZero() {

		flytura.FormataRetornoHTTP(w, "ID do usuário inválido", http.StatusBadRequest)

		// http.Error(w, "ID do usuário inválido", http.StatusBadRequest)
		return
	}

	if user.UpdatedAt.IsZero() {
		user.UpdatedAt = time.Now()
	}

	// Criar o objeto de atualização
	update := bson.M{
		"$set": bson.M{
			"name":           user.Name,
			"lastName":       user.LastName,
			"passportNumber": user.PassportNumber,
			"perfil":         user.Perfil,
			"updatedAt":      user.UpdatedAt,
			"idUserUpdate":   user.ID.Hex(),
			"active":         user.Active,
			"mobile":         user.Mobile,
		},
	}

	// Se uma nova senha for fornecida, gerar um hash
	if user.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			flytura.FormataRetornoHTTP(w, "Erro ao gerar hash da senha,Erro ao processar a senha ", http.StatusInternalServerError)

			// log.Println("Erro ao gerar hash da senha:", err)
			// http.Error(w, "Erro ao processar a senha", http.StatusInternalServerError)
			return
		}
		update["$set"].(bson.M)["password"] = string(hashedPassword)
	}

	// Conectar ao MongoDB e atualizar o usuário
	client, err := db.ConnectMongoDB(flytura.ConectionString)
	if err != nil {
		flytura.FormataRetornoHTTP(w, "Erro ao conectar ao banco de dados", http.StatusInternalServerError)

		// http.Error(w, "Erro ao conectar ao banco de dados", http.StatusInternalServerError)
		return
	}
	defer client.Disconnect(context.Background())

	collection := client.Database(flytura.DBName).Collection("user")
	result, err := collection.UpdateOne(context.Background(), bson.M{"_id": user.ID}, update)
	if err != nil {
		flytura.FormataRetornoHTTP(w, "Erro ao atualizar usuário, Erro ao atualizar usuário", http.StatusInternalServerError)

		// log.Println("Erro ao atualizar usuário:", err)
		// http.Error(w, "Erro ao atualizar usuário", http.StatusInternalServerError)
		return
	}

	// Verifica se algum documento foi modificado
	if result.ModifiedCount == 0 {
		flytura.FormataRetornoHTTP(w, "Nenhuma alteração realizada", http.StatusOK)

		// http.Error(w, "Nenhuma alteração realizada", http.StatusNotModified)
		return
	}

	// Responder com sucesso
	flytura.FormataRetornoHTTP(w, "Usuário atualizado com sucesso! Documento modificado", http.StatusOK)

}
