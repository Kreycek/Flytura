package airLine

import (
	flytura "Flytura"
	"Flytura/internal/db"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

/*
Função criada por Ricardo Silva Ferreira
Inicio da criação 05/09/2025 14:04
Data Final da criação : 05/09/2025 14:33
*/
// Obtem todos sem paginação
func GetAllAirLineHandler(w http.ResponseWriter, r *http.Request) {

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
	airline, err := GetAirLines(client, flytura.DBName, flytura.Airline)
	if err != nil {
		http.Error(w, fmt.Sprintf("erro ao buscar companias: %v", err), http.StatusInternalServerError)
		return
	}

	// Retornar a resposta com os dados dos usuários
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(airline); err != nil {
		log.Printf("erro ao codificar resposta JSON: %v", err)
	}
}
