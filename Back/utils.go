// config/config.go
package clarion

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt"
)

// Variável global que contém a chave secreta para JWT
var SecretKey = []byte("my_secret_key")
var UrlSite = "http://localhost:52912"
var ConectionString = "mongodb://admin:secret@localhost:27017"
var DBName = "flytura"

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

func FormataRetornoHTTP(w http.ResponseWriter, mensagem any, codHttp int) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(codHttp) // Código 200 OK
	return json.NewEncoder(w).Encode(map[string]any{"message": mensagem})
}
func FormataRetornoHTTPGeneric(w http.ResponseWriter, bodyName string, body any, codHttp int) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(codHttp) // Código 200 OK
	return json.NewEncoder(w).Encode(map[any]any{"users": body})
}
