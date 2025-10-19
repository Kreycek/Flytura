package awsS3

import (
	flytura "Flytura"
	"Flytura/internal/db"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
)

/*
Função criada por Ricardo Silva Ferreira
Início da criação: 17/10/2025 13:02
Data final da criação:  17/10/2025 13:005
*/

func UploadS3FilesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	// Processa o corpo da requisição como multipart
	err := r.ParseMultipartForm(10 << 20) // até 10 MB
	if err != nil {
		http.Error(w, "Erro ao processar formulário", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Erro ao ler arquivo", http.StatusBadRequest)
		return
	}
	defer file.Close()

	companyCode := r.FormValue("companyCode")

	err = UploadToS3(file, header.Filename, companyCode)
	if err != nil {

		fmt.Println("Erro detalhado ao enviar para S3:", err)
		// return fmt.Errorf("erro ao enviar para S3: %w", err)

		return
	}

	fmt.Println("File name ", header.Filename)
	fmt.Println("CompanyCode ", companyCode)

	fmt.Fprintf(w, "Arquivo %s enviado com sucesso!", header.Filename)
}

/*
Função criada por Ricardo Silva Ferreira
Inicio da criação 19/10/2025 19:14
Data Final da criação : 19/10/2025 19:17
*/

func SearchS3ImagesDBHandler(w http.ResponseWriter, r *http.Request) {
	// Verificar se a requisição é do tipo POST

	if r.Method != http.MethodGet {
		http.Error(w, "Método não permitido, deve ser GET", http.StatusMethodNotAllowed)
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

	query := r.URL.Query()

	companyCode := query.Get("companyCode")
	startDateStr := query.Get("startDate")

	var startDate *time.Time
	if startDateStr != "" {
		t, err := time.Parse(time.RFC3339, startDateStr)
		if err == nil {
			startDate = &t
		} else {
			// lidar com erro de parsing, se necessário
			fmt.Println("Erro ao converter startDate:", err)
		}
	}

	endDateStr := query.Get("endDate")

	var endDate *time.Time
	if endDateStr != "" {
		t, err := time.Parse(time.RFC3339, endDateStr)
		if err == nil {
			endDate = &t
		} else {
			// lidar com erro de parsing, se necessário
			fmt.Println("Erro ao converter startDate:", err)
		}
	}

	pageStr := query.Get("page")

	page, err := strconv.ParseInt(pageStr, 10, 64)
	if err != nil {
		page = 10 // valor padrão se a conversão falhar
	}
	if page < 1 {
		page = 1
	}

	limitStr := query.Get("limit")

	limit, err := strconv.ParseInt(limitStr, 10, 64)
	if err != nil {
		limit = 10 // valor padrão se a conversão falhar
	}

	fmt.Println("StartDate", startDate)
	fmt.Println("EndDate", endDateStr)
	// Definir valores padrão para paginação

	if limit < 1 {
		limit = 10
	}

	// Buscar imagens com paginação
	imagesDb, total, err := SearchImagesDB(
		client,
		flytura.DBName,
		"imagesDB",
		&companyCode,
		startDate,
		endDate,
		page,
		limit)

	if err != nil {
		http.Error(w, "Erro ao buscar imagens", http.StatusInternalServerError)
		return
	}

	// Criar resposta JSON com paginação
	response := map[string]any{
		"total":    total,
		"page":     pageStr,
		"limit":    limit,
		"pages":    (total + limit - 1) / limit, // Número total de páginas
		"imagesDB": imagesDb,
	}

	// Retornar resposta JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("erro ao codificar resposta JSON: %v", err)
	}
}
