package auth

import (
	flytura "Flytura"
	"Flytura/internal/db"
	"Flytura/internal/models"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

// Variáveis globais

// Função para verificar o nome de usuário e senha
func VerifyUser(w http.ResponseWriter, r *http.Request) {
	// Parse o corpo da requisição
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Erro ao ler o corpo da requisição", http.StatusBadRequest)
		return
	}

	// Conectar ao MongoDB
	client, err := db.ConnectMongoDB(flytura.ConectionString)
	if err != nil {
		http.Error(w, "Erro ao conectar ao banco de dados", http.StatusInternalServerError)
		return
	}
	defer db.CloseMongoDB(client)

	// Obter a coleção de usuários
	collection := db.GetCollection(client, flytura.DBName, flytura.UserDBTableName)

	filter := bson.D{
		{Key: "$or", Value: bson.A{
			bson.D{{Key: "userName", Value: user.Username}},
			bson.D{{Key: "email", Value: user.Username}},
		}},
		{Key: "active", Value: true}, // Adicionando a condição para o campo active ser true
	}

	// Verificar se o usuário existe no banco
	var result bson.M
	err = collection.FindOne(context.Background(), filter).Decode(&result)
	if err != nil {
		flytura.FormataRetornoHTTP(w, "Erro geral", http.StatusUnauthorized)
		return
		// log.Fatal(err)
	}
	if err == mongo.ErrNoDocuments {
		flytura.FormataRetornoHTTP(w, "Usuário não encontrado", http.StatusUnauthorized)
		// http.Error(w, "Usuário não encontrado", http.StatusUnauthorized)
		return
	}
	if err != nil {
		http.Error(w, fmt.Sprintf("Erro ao buscar o usuário: %v", err), http.StatusInternalServerError)
		return
	}

	// Comparar a senha fornecida com a senha armazenada no banco de dados
	storedPassword, ok := result["password"].(string)
	if !ok {
		flytura.FormataRetornoHTTP(w, "Erro na senha armazenada", http.StatusUnauthorized)
		// http.Error(w, "Erro na senha armazenada", http.StatusInternalServerError)
		return
	}

	// Verificar se a senha fornecida corresponde à senha armazenada
	err = bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(user.Password))
	if err != nil {
		flytura.FormataRetornoHTTP(w, "Senha incorreta", http.StatusUnauthorized)
		// http.Error(w, "Senha incorreta", http.StatusUnauthorized)
		return
	}

	// Criar o token JWT
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["username"] = user.Username
	claims["perfis"] = result["perfil"]
	claims["name"] = result["name"].(string)
	claims["lastName"] = result["lastName"].(string)
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	// Criar o token
	tokenString, err := token.SignedString(flytura.SecretKey)
	if err != nil {
		flytura.FormataRetornoHTTP(w, "Erro ao criar o token", 401)
		// http.Error(w, "Erro ao criar o token", http.StatusInternalServerError)
		return
	}

	// Retornar o token para o cliente
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
}

// Função para validar o token JWT
func ValidateToken(w http.ResponseWriter, r *http.Request) {
	// Recuperar o token do cabeçalho 'Authorization'
	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
		flytura.FormataRetornoHTTP(w, "Token não fornecido", http.StatusUnauthorized)
		// http.Error(w, "Token não fornecido", http.StatusUnauthorized)
		return
	}

	// Remover o prefixo 'Bearer ' caso esteja presente
	if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
		tokenString = tokenString[7:]
	}

	// Parse e validação do token JWT
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Método de assinatura inválido")
		}
		return flytura.SecretKey, nil // jwtKey é a chave secreta para validação
	})

	if err != nil {
		// formataRetornoHTTP(w, "Token inválido")
		http.Error(w, "Token inválido", http.StatusUnauthorized)
		return
	}

	// Validar o token
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		fmt.Println("Token válido. Usuário:", claims["username"])

		// Enviar uma resposta com o status de sucesso e a mensagem "Token válido"

		flytura.FormataRetornoHTTP(w, "Token válido", http.StatusOK)
	} else {
		http.Error(w, "Token inválido", http.StatusUnauthorized)
	}
}
