package perfil

import (
	flytura "Flytura"
	"Flytura/internal/db"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// Função de handler para a rota GET /users
func GetAllPerfilsHandler(w http.ResponseWriter, r *http.Request) {

	status, msg := flytura.TokenValido(w, r)

	if !status {
		http.Error(w, fmt.Sprintf("erro ao validar token: %v", msg), http.StatusUnauthorized)
		return
	}

	// Conectar ao MongoDB
	client, err := db.ConnectMongoDB(flytura.ConectionString)
	if err != nil {
		http.Error(w, fmt.Sprintf("erro ao conectar ao MongoDB: %v", err), http.StatusInternalServerError)
		return
	}
	defer db.CloseMongoDB(client)

	// Obter todos os usuários
	users, err := GetAllPerfil(client, flytura.DBName, "perfil")
	if err != nil {
		http.Error(w, fmt.Sprintf("erro ao buscar usuários: %v", err), http.StatusInternalServerError)
		return
	}

	// Retornar a resposta com os dados dos usuários
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(users); err != nil {
		log.Printf("erro ao codificar resposta JSON: %v", err)
	}
}
