// config/config.go
package flytura

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt"
)

// Variável global que contém a chave secreta para JWT
var SecretKey = []byte("my_secret_key")

var UrlSiteLocalHost = "localhost:4200"
var UrlSiteProduction = "54.156.244.197"
var UrlSiteHomol = "18.210.18.180"
var BucketProd = "flytura-bucket"
var BucketRegion = "us-east-1"
var AccessKeyAws = os.Getenv("AWS_ACCESS_KEY_ID")
var SecretKeyAws = os.Getenv("AWS_SECRET_ACCESS_KEY")

//para criar as variáveis no powerShell
//[System.Environment]::SetEnvironmentVariable("AWS_ACCESS_KEY_ID", "AKIATTTYBPAD3ENZW3KB", "User")
//[System.Environment]::SetEnvironmentVariable("AWS_SECRET_ACCESS_KEY", "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", "User")

var ConectionString = "mongodb://admin:secret@localhost:27017"

// var ConectionString = "mongodb://localhost:27017"
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
