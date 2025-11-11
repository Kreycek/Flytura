package outPutInvoices

import (
	flytura "Flytura"
	"Flytura/internal/db"
	"Flytura/internal/models"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
)

/*
	Função criada por Ricardo Silva Ferreira
	Inicio da criação 11/11/2025 11:25
	Data Final da criação : 11/11/2025 11:31
*/

func SearchOutPutInvoicesHandler(w http.ResponseWriter, r *http.Request) {
	// Verificar se a requisição é do tipo POST
	if r.Method != http.MethodGet {
		http.Error(w, "Método não permitido dever ser um get", http.StatusMethodNotAllowed)
		return
	}

	// Validar Token
	status, msg := flytura.TokenValido(w, r)
	if !status {
		http.Error(w, fmt.Sprintf("erro ao validar token: %v", msg), http.StatusUnauthorized)
		return
	}

	query := r.URL.Query()
	page, err := strconv.Atoi(query.Get("page"))
	if err != nil || page < 1 {
		page = 1 // Padrão: primeira página
	}

	limit, err := strconv.Atoi(query.Get("limit"))
	if err != nil || limit < 1 {
		limit = 10 // Padrão: 10 registros por página
	}

	keyCode := query.Get("key")

	companyCode := query.Get("companyCode")

	startDate := query.Get("startDate")

	endDate := query.Get("endDate")

	var _startDate *time.Time = nil
	var _endDate *time.Time = nil

	fmt.Println("Data 2", startDate)

	if startDate != "" && endDate != "" {

		var stIniError error

		parsedStart, stIniError := time.Parse("2006-01-02T15:04:05Z", startDate)
		if stIniError != nil {
			http.Error(w, "Formato de data inválido", http.StatusInternalServerError)
			return
		}
		_startDate = &parsedStart

		parsedEnd, stIniError := time.Parse("2006-01-02T15:04:05Z", endDate)
		if stIniError != nil {
			http.Error(w, "Formato de data inválido", http.StatusInternalServerError)
			return
		}
		_endDate = &parsedEnd
	}
	// fmt.Println("StartDate", request.StartDate)
	// fmt.Println("EndDate", request.EndDate)
	// Definir valores padrão para paginação
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	// Buscar usuários com paginação
	outPutInvoices, total, err := SearchOutPutInvoicesPagination(
		db.MongoClient,
		flytura.DBName,
		flytura.OutPutInvoices,
		&keyCode,
		&companyCode,
		_startDate,
		_endDate,
		page,
		limit)

	if err != nil {
		fmt.Println(err)
		http.Error(w, "Erro ao buscar faturas", http.StatusInternalServerError)
		return
	}

	// Criar resposta JSON com paginação
	response := map[string]any{
		"total":          total,
		"page":           page,
		"limit":          limit,
		"pages":          (total + int64(limit) - 1) / int64(limit), // Número total de páginas
		"outPutInvoices": outPutInvoices,
	}

	// Retornar resposta JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("erro ao codificar resposta JSON: %v", err)
	}
}

/*
Função criada por Ricardo Silva Ferreira
Inicio da criação 11/11/2025 11:40
Data Final da criação : 11/11/2025 11:48
*/

func InsertOutPutInvoicesHandler(w http.ResponseWriter, r *http.Request) {
	// Validar o token de autenticação
	if r.Method != http.MethodPost {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	token := r.Header.Get("token")

	if token == "" {
		http.Error(w, "Token não fornecido", http.StatusUnauthorized)
		return
	}

	var errorValidaToken error
	_, errorValidaToken = flytura.VerifyAccessValidTokenListSheet(db.MongoClient, flytura.DBName, flytura.TokenAccessTableName, token)
	if errorValidaToken != nil {
		http.Error(w, "Token inválido", http.StatusBadRequest)
		return
	}

	// Ler o corpo da requisição
	var data []models.OutPutInvoices
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, "erro ao decodificar corpo da requisição", http.StatusBadRequest)
		return
	}

	// Conectar ao MongoDB
	// client, err := db.ConnectMongoDB(flytura.ConectionString)
	// if err != nil {
	// 	http.Error(w, fmt.Sprintf("erro ao conectar ao MongoDB: %v", err), http.StatusInternalServerError)
	// 	return
	// }
	// defer db.CloseMongoDB(client)

	// Inserir o usuário no MongoDB
	err = InsertOutPutInvoices(db.MongoClient, flytura.DBName, flytura.OutPutInvoices, data)
	if err != nil {
		http.Error(w, fmt.Sprintf("erro ao inserir fatura: %v", err), http.StatusInternalServerError)
		return
	}

	// Retornar a resposta com os dados dos usuários
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("erro ao codificar resposta JSON: %v", err)
	}

}
