package bigQueryOnfly

import (
	flytura "Flytura"
	"Flytura/internal/db"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
)

/*
Função criada por Ricardo Silva Ferreira
Inicio da criação 22/05/2026 15:08
Data Final da criação :  22/05/2026 15:08
*/
func SearchConciliationPaginationHandler(w http.ResponseWriter, r *http.Request) {
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

	originLocator := query.Get("originLocator")
	returnLocator := query.Get("returnLocator")

	originETicket := query.Get("originETicket")
	returnETicket := query.Get("returnETicket")

	startDate := query.Get("startDate")

	endDate := query.Get("endDate")

	var _startDate *time.Time = nil
	var _endDate *time.Time = nil

	// fmt.Println("Data 2", startDate)

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
	conciliation, total, err := SearchConciliationPagination(
		db.MongoClient,
		flytura.DBName,
		flytura.ConciliationTableName,
		originLocator,
		returnLocator,
		originETicket,
		returnETicket,
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
		"total":        total,
		"page":         page,
		"limit":        limit,
		"pages":        (total + int64(limit) - 1) / int64(limit), // Número total de páginas
		"conciliation": conciliation,
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
Inicio da criação 22/05/2026 15:30
Data Final da criação :  22/05/2026 15:32
*/
func SearchConciliationExcelHandler(w http.ResponseWriter, r *http.Request) {
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

	originLocator := query.Get("originLocator")
	returnLocator := query.Get("returnLocator")

	originETicket := query.Get("originETicket")
	returnETicket := query.Get("returnETicket")

	startDate := query.Get("startDate")

	endDate := query.Get("endDate")

	var _startDate *time.Time = nil
	var _endDate *time.Time = nil

	// fmt.Println("Data 2", startDate)

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

	// Buscar usuários com paginação
	outPutInvoices, total, err := SearchConciliationExcel(
		db.MongoClient,
		flytura.DBName,
		flytura.ConciliationTableName,
		originLocator,
		returnLocator,
		originETicket,
		returnETicket,
		_startDate,
		_endDate)

	if err != nil {
		fmt.Println(err)
		http.Error(w, "Erro ao buscar faturas", http.StatusInternalServerError)
		return
	}

	// Criar resposta JSON com paginação
	response := map[string]any{
		"total":        total,
		"conciliation": outPutInvoices,
	}

	// Retornar resposta JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("erro ao codificar resposta JSON: %v", err)
	}
}
