package purcharseRecord

import (
	flytura "Flytura"
	"Flytura/internal/airLine"
	"Flytura/internal/db"
	"Flytura/internal/models"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func criarArquivoTemporario(extensao string) string {
	if extensao == ".xls" {
		return "upload-*.xls"
	}
	return "upload-*.xlsx"
}

/*
Função criada por Ricardo Silva Ferreira
Inicio da criação 03/09/2025 22:05
Data Final da criação : 05/09/2025 22:40
*/

func UploadPurcharseRecordHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(10 << 20) // 10MB
	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Erro ao receber arquivo", http.StatusBadRequest)
		return
	}
	defer file.Close()

	client, errConnectDB1 := db.ConnectMongoDB(flytura.ConectionString)
	if errConnectDB1 != nil {
		log.Println("Erro ao obter nome do arquivo:", err)
		return
	}
	defer db.CloseMongoDB(client)

	nameWithoutExt := strings.TrimSuffix(fileHeader.Filename, filepath.Ext(fileHeader.Filename))

	parts := strings.Split(nameWithoutExt, "-")

	// fmt.Println("Nome do arquivo sem extensão:", parts[1])
	codeFile := ""
	if len(parts) > 1 {

		fmt.Println("Nome do arquivo limpo 1:", strings.Join(strings.Fields(parts[1]), ""))
		codeFile = strings.Join(strings.Fields(parts[1]), "")[:4] //retira todos os espaços da string e pega apenas os 4 primeiros caracteres

		// fmt.Println("Nome do arquivo limpo:", codeFile)
	} else {
		codeFile = "error"
		log.Println("Nome do arquivo inválido:", err)
		return
	}

	airLineData, err := airLine.GetAirLineFileName(client, flytura.DBName, "airline", codeFile)
	if err != nil {
		log.Println("Erro ao obter nome do arquivo:", err)
		return
	}

	// fmt.Println("ddd", airLineData["FileName"])

	companyName := airLineData["Name"].(string)
	companyCode := airLineData["Code"].(string)

	// fileName, ok := rr["fileName"]
	// if !ok {
	// 	log.Println("Chave 'fileName' não encontrada")
	// 	return
	// }

	extensao := strings.ToLower(filepath.Ext(fileHeader.Filename))

	tempFile, errTempFile := os.CreateTemp("", criarArquivoTemporario(extensao))
	if errTempFile != nil {
		http.Error(w, "Erro ao salvar arquivo", http.StatusInternalServerError)
		return
	}

	defer os.Remove(tempFile.Name())

	io.Copy(tempFile, file)

	clients, err2 := db.ConnectMongoDB(flytura.ConectionString)
	if err2 != nil {
		http.Error(w, "Erro ao conectar ao banco de dados", http.StatusInternalServerError)
		return
	}
	defer db.CloseMongoDB(clients)

	Results := struct {
		TotalEmpty         int64 `json:"totalEmpty"`
		TotalRecordsImport int64 `json:"totalRecordsImport"`
		EmptySheet         bool  `json:"emptySheet"`
	}{
		TotalEmpty:         0,
		TotalRecordsImport: 0,
		EmptySheet:         false,
	}

	totalEmptyRegister, totalRegister, emptySheet, errProcessExcel := ProcessExcel(tempFile.Name(), fileHeader.Filename, companyName, companyCode, clients, flytura.DBName, flytura.PurcharseRecordTableName)
	if errProcessExcel != nil {
		http.Error(w, "Erro ao processar planilha", http.StatusInternalServerError)
		return
	}

	// fmt.Println("vazio", emptySheet)

	Results.EmptySheet = emptySheet
	Results.TotalRecordsImport = totalRegister

	if totalEmptyRegister > 0 {
		Results.TotalEmpty = totalEmptyRegister
	}

	// Retornar a resposta com os dados dos usuários

	flytura.FormataRetornoHTTP(w, Results, http.StatusOK)

}

/*
Função criada por Ricardo Silva Ferreira
Inicio da criação 05/09/2025 14:04
Data Final da criação : 05/09/2025 14:33
*/
// Obtem todos sem paginação
func GetAllPurcharseRecordHandler(w http.ResponseWriter, r *http.Request) {

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
	onlyFlyData, err := GetExcelData(client, flytura.DBName, flytura.PurcharseRecordTableName)
	if err != nil {
		http.Error(w, fmt.Sprintf("erro ao buscar usuários: %v", err), http.StatusInternalServerError)
		return
	}

	// Retornar a resposta com os dados dos usuários
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(onlyFlyData); err != nil {
		log.Printf("erro ao codificar resposta JSON: %v", err)
	}
}

/*
Função criada por Ricardo Silva Ferreira
Inicio da criação 05/09/2025 10:45
Data Final da criação : 05/09/2025 11:00
*/

// Não permite pesquisar por parametro apenas traz todos os registro para paginação
func GetAllPurcharseRecordPaginationHandler(w http.ResponseWriter, r *http.Request) {
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
	data, total, err := GetAllExcelData(client, flytura.DBName, flytura.PurcharseRecordTableName, page, limit)
	if err != nil {
		http.Error(w, fmt.Sprintf("erro ao buscar diários: %v", err), http.StatusInternalServerError)
		return
	}

	// Criar resposta JSON com paginação
	response := map[string]any{
		"total":           total,
		"page":            page,
		"limit":           limit,
		"pages":           (total + limit - 1) / limit, // Calcula o número total de páginas
		"purcharseRecord": data,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("erro ao codificar resposta JSON: %v", err)
	}
}

/*
Função criada por Ricardo Silva Ferreira
Inicio da criação 05/09/2025 11:23
Data Final da criação : 05/09/2025 11:34
*/

func GetPurcharseRecordByIdHandler(w http.ResponseWriter, r *http.Request) {
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
	costCenters, err := GetExcelDataByID(client, flytura.DBName, flytura.PurcharseRecordTableName, id)
	if err != nil {
		http.Error(w, "Erro ao buscar diários", http.StatusInternalServerError)
		return
	}

	// Configurar o cabeçalho da resposta como JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Enviar o usuário como resposta JSON
	if err := json.NewEncoder(w).Encode(costCenters); err != nil {
		log.Printf("erro ao codificar resposta JSON: %v", err)
		http.Error(w, "Erro ao codificar resposta", http.StatusInternalServerError)
	}
}

/*
Função criada por Ricardo Silva Ferreira
Inicio da criação 05/09/2025 11:05
Data Final da criação : 05/09/2025 11:15
*/

func InsertPurcharseRecordHandler(w http.ResponseWriter, r *http.Request) {
	// Validar o token de autenticação
	status, msg := flytura.TokenValido(w, r)
	if !status {
		http.Error(w, fmt.Sprintf("erro ao validar token: %v", msg), http.StatusUnauthorized)
		return
	}

	// Ler o corpo da requisição
	var data models.PurcharseRecord
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, "erro ao decodificar corpo da requisição", http.StatusBadRequest)
		return
	}

	if !data.Active {
		data.Active = true
	}

	if data.CreatedAt.IsZero() {
		data.CreatedAt = time.Now()
	}
	// Conectar ao MongoDB
	client, err := db.ConnectMongoDB(flytura.ConectionString)
	if err != nil {
		http.Error(w, fmt.Sprintf("erro ao conectar ao MongoDB: %v", err), http.StatusInternalServerError)
		return
	}
	defer db.CloseMongoDB(client)

	// Inserir o usuário no MongoDB
	err = InsertExcelData(client, flytura.DBName, flytura.PurcharseRecordTableName, data)
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

/*
Função criada por Ricardo Silva Ferreira
Inicio da criação 05/09/2025 12:01
Data Final da criação : 05/09/2025 12:07
*/

func UpdatePurcharseRecordHandler(w http.ResponseWriter, r *http.Request) {
	// Validar o token de autenticação
	status, msg := flytura.TokenValido(w, r)
	if !status {
		http.Error(w, fmt.Sprintf("Erro ao validar token: %v", msg), http.StatusUnauthorized)
		return
	}

	// Decodificar o JSON recebido
	var data models.PurcharseRecord
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		// http.Error(w, "Erro ao decodificar JSON", http.StatusBadRequest)
		flytura.FormataRetornoHTTP(w, "Erro ao decodificar JSON", http.StatusBadRequest)

		return
	}

	// Verifica se o ID é válido
	if data.ID.IsZero() {

		flytura.FormataRetornoHTTP(w, "ID da fatura inválido", http.StatusBadRequest)

		// http.Error(w, "ID do usuário inválido", http.StatusBadRequest)
		return
	}

	if data.UpdatedAt.IsZero() {
		data.UpdatedAt = time.Now()
	}

	// Criar o objeto de atualização
	update := bson.M{
		"$set": bson.M{

			"key":          data.Key,
			"name":         data.Name,
			"lastName":     data.LastName,
			"companyCode":  data.CompanyCode,
			"companyName":  data.CompanyName,
			"fileName":     data.FileName,
			"status":       data.Status,
			"idUserUpdate": data.ID.Hex(),
			"active":       data.Active,
		},
	}

	// Conectar ao MongoDB e atualizar o usuário
	client, err := db.ConnectMongoDB(flytura.ConectionString)
	if err != nil {
		flytura.FormataRetornoHTTP(w, "Erro ao conectar ao banco de dados", http.StatusInternalServerError)

		// http.Error(w, "Erro ao conectar ao banco de dados", http.StatusInternalServerError)
		return
	}
	defer client.Disconnect(context.Background())

	collection := client.Database(flytura.DBName).Collection(flytura.PurcharseRecordTableName)
	result, err := collection.UpdateOne(context.Background(), bson.M{"_id": data.ID}, update)
	if err != nil {
		flytura.FormataRetornoHTTP(w, "Erro ao atualizar fatura", http.StatusInternalServerError)

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
	flytura.FormataRetornoHTTP(w, "Fatura cadastrada com sucesso", http.StatusOK)

}

/*
Função criada por Ricardo Silva Ferreira
Inicio da criação 05/09/2025 13:30
Data Final da criação : 05/09/2025 13:33
*/

// Função para verificar o nome de usuário e senha
func VerifyExistPurcharseRecordHandler(w http.ResponseWriter, r *http.Request) {

	status, msg := flytura.TokenValido(w, r)
	if !status {
		http.Error(w, fmt.Sprintf("erro ao validar token: %v", msg), http.StatusUnauthorized)
		return
	}

	// Parse o corpo da requisição
	var cc models.PurcharseRecordVerifyExistRequest

	err := json.NewDecoder(r.Body).Decode(&cc)
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
	collection := db.GetCollection(client, flytura.DBName, flytura.PurcharseRecordTableName)
	// filter := bson.D{
	// 	{Key: "$or", Value: bson.A{
	// 		bson.D{{Key: "email", Value: userName}},
	// 	}},
	// }

	filter := bson.D{{Key: "key", Value: cc.Key}}

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

/*
Função criada por Ricardo Silva Ferreira
Inicio da criação 05/09/2025 13:39
Data Final da criação : 05/09/2025 14:03
*/

func SearchPurcharseRecordHandler(w http.ResponseWriter, r *http.Request) {
	// Verificar se a requisição é do tipo POST
	if r.Method != http.MethodPost {
		http.Error(w, "Método não permitido dever ser um post", http.StatusMethodNotAllowed)
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
		Key           *string    `json:"key"`
		Name          *string    `json:"name"`
		LastName      *string    `json:"lastName"`
		CompanyCode   *string    `json:"companyCode"`
		MessageReturn *string    `json:"messageReturn"`
		StartDate     *time.Time `json:"startDate"`
		EndDate       *time.Time `json:"endDate"`
		Status        *string    `json:"status"`
		Page          int64      `json:"page"`
		Limit         int64      `json:"limit"`
	}

	// Decodificar o corpo da requisição JSON
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Erro ao decodificar JSON", http.StatusBadRequest)
		return
	}

	fmt.Println("StartDate", request.StartDate)
	fmt.Println("EndDate", request.EndDate)
	// Definir valores padrão para paginação
	if request.Page < 1 {
		request.Page = 1
	}
	if request.Limit < 1 {
		request.Limit = 10
	}

	// Buscar usuários com paginação
	costCenters, total, err := SearchExcelData(
		client,
		flytura.DBName,
		flytura.PurcharseRecordTableName,
		request.Key,
		request.Name,
		request.LastName,
		request.CompanyCode,
		request.StartDate,
		request.EndDate,
		request.Status,
		request.Page,
		request.Limit)

	if err != nil {
		http.Error(w, "Erro ao buscar faturas", http.StatusInternalServerError)
		return
	}

	// Criar resposta JSON com paginação
	response := map[string]any{
		"total":           total,
		"page":            request.Page,
		"limit":           request.Limit,
		"pages":           (total + request.Limit - 1) / request.Limit, // Número total de páginas
		"purcharseRecord": costCenters,
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
Inicio da criação 30/09/2025 17:11
Data Final da criação : 30/09/2025 17:12
*/
// Obtem todos sem paginação
func GetAllImportStatussHandler(w http.ResponseWriter, r *http.Request) {

	status, msg := flytura.TokenValido(w, r)

	if !status {
		http.Error(w, fmt.Sprintf("Token inválido: %v", msg), http.StatusUnauthorized)
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
	onlyFlyData, err := GetExcelData(client, flytura.DBName, flytura.StatusImportTableName)
	if err != nil {
		http.Error(w, fmt.Sprintf("erro ao buscar usuários: %v", err), http.StatusInternalServerError)
		return
	}

	// Retornar a resposta com os dados dos usuários
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(onlyFlyData); err != nil {
		log.Printf("erro ao codificar resposta JSON: %v", err)
	}
}

/*
Função criada por Ricardo Silva Ferreira
Início da criação: 01/10/2025 17:00
Data final da criação: 01/10/2025 17:01
*/
func GroupByCompanyNameHandler(w http.ResponseWriter, r *http.Request) {

	status, msg := flytura.TokenValido(w, r)
	if !status {
		http.Error(w, fmt.Sprintf("token inválido handler onlyfly: %v", msg), http.StatusUnauthorized)
		return
	}

	// Conectar ao MongoDB
	client, err := db.ConnectMongoDB(flytura.ConectionString)
	if err != nil {
		http.Error(w, fmt.Sprintf("erro ao conectar ao MongoDB: %v", err), http.StatusInternalServerError)
		return
	}
	defer db.CloseMongoDB(client)

	// Obter parâmetros da query string
	query := r.URL.Query()
	statusParam := query.Get("status")
	companyNameParam := query.Get("companyName")
	startDateStr := query.Get("startDate")
	endDateStr := query.Get("endDate")

	// Parse das datas, se fornecidas
	var startDate, endDate *time.Time
	if startDateStr != "" && endDateStr != "" {
		start, errStart := time.Parse(time.RFC3339, startDateStr)
		end, errEnd := time.Parse(time.RFC3339, endDateStr)
		if errStart == nil && errEnd == nil {
			startDate = &start
			endDate = &end
		}
	}

	// Chamada da função com filtros
	onlyFlyData, err := GroupByCompanyNameFiltered(client, flytura.DBName, flytura.PurcharseRecordTableName, startDate, endDate, statusParam, companyNameParam)
	if err != nil {
		http.Error(w, fmt.Sprintf("erro ao buscar usuários: %v", err), http.StatusInternalServerError)
		return
	}

	// Retornar a resposta com os dados
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(onlyFlyData); err != nil {
		log.Printf("erro ao codificar resposta JSON: %v", err)
	}
}

/*
Função criada por Ricardo Silva Ferreira
Inicio da criação 14/10/2025 22:01
Data Final da criação :  14/10/2025 22:12
*/
func GetPurcharseRecordStatusHandler(w http.ResponseWriter, r *http.Request) {

	// Conectar ao MongoDB
	client, err := db.ConnectMongoDB(flytura.ConectionString)
	if err != nil {
		http.Error(w, "Erro ao conectar ao MongoDB", http.StatusInternalServerError)
		return
	}
	defer client.Disconnect(context.Background())

	// Extrair o ID da URL
	companyCode := r.URL.Query().Get("companyCode")
	if companyCode == "" {
		http.Error(w, "Codigo da compania não fornecido", http.StatusBadRequest)
		return
	}

	status := r.URL.Query().Get("status")
	if status == "" {
		http.Error(w, "status não fornecido", http.StatusBadRequest)
		return
	}

	token := r.URL.Query().Get("token")
	if token == "" {
		http.Error(w, "token não fornecido", http.StatusBadRequest)
		return
	}

	var error error
	_, error = VerifyAccessValidTokenListSheet(client, flytura.DBName, flytura.TokenAccessTableName, token)
	if error != nil {
		http.Error(w, "Token inválido", http.StatusBadRequest)
		log.Println("Token inválido", err)
		return
	}

	// Buscar o usuário no banco de dados pelo ID
	dataExcel, err := GetDataExcelByStatus(client, flytura.DBName, flytura.PurcharseRecordTableName, companyCode, status)
	if err != nil {
		http.Error(w, "Erro ao dados da planilha", http.StatusInternalServerError)
		return
	}

	// Configurar o cabeçalho da resposta como JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Enviar o usuário como resposta JSON
	if err := json.NewEncoder(w).Encode(dataExcel); err != nil {
		log.Printf("erro ao codificar resposta JSON: %v", err)
		http.Error(w, "Erro ao codificar resposta", http.StatusInternalServerError)
	}
}

/*
Função criada por Ricardo Silva Ferreira
Inicio da criação 14/10/2025 22:24
Data Final da criação :  14/10/2025 22:32
*/
func VerifyAccessValidTokenListSheet(client *mongo.Client, dbName, collectionName, token string) (map[string]any, error) {

	collection := client.Database(dbName).Collection(collectionName)

	// objectID, erroId := primitive.ObjectIDFromHex(excelId)
	// if erroId != nil {
	// 	log.Fatalf("Erro ao converter string para ObjectID: %v", erroId)
	// }

	filter := bson.M{
		"token":  token,
		"active": true,
	}

	// Variável para armazenar o usuário retornado
	var tokenAccessSheet models.TokenAccessSheet

	// Usar FindOne para pegar apenas um único registro
	err := collection.FindOne(context.Background(), filter).Decode(&tokenAccessSheet)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("token não encontrado")
		}
		return nil, fmt.Errorf("erro ao buscar plano de contas: %v", err)
	}

	// Converter o _id para string

	// Retornar o usuário como um mapa
	_return := map[string]any{
		"ID":   tokenAccessSheet.ID.Hex(), // Agora o campo ID é uma string
		"Name": tokenAccessSheet.Name,
		"Code": tokenAccessSheet.Token,
	}

	return _return, nil
}

type RequestModelPurcharseRecord struct {
	Key           string `json:"key" bson:"key,omitempty"`
	Status        string `json:"status" bson:"status,omitempty"`
	MessageReturn string `json:"messageReturn" bson:"messageReturn,omitempty"`
}

/*
Função criada por Ricardo Silva Ferreira
Inicio da criação 21/10/2025 16:19
Data Final da criação : 21/10/2025 16:21
*/
func UpdatePurcharseRecordMultipleHandler(w http.ResponseWriter, r *http.Request) {
	// Validar o token de autenticação

	// Decodificar o JSON recebido como array
	var records []RequestModelPurcharseRecord
	if err := json.NewDecoder(r.Body).Decode(&records); err != nil {
		flytura.FormataRetornoHTTP(w, "Erro ao decodificar JSON", http.StatusBadRequest)
		return
	}

	fmt.Println("Header", records)

	// Conectar ao MongoDB
	client, err := db.ConnectMongoDB(flytura.ConectionString)
	if err != nil {
		flytura.FormataRetornoHTTP(w, "Erro ao conectar ao banco de dados", http.StatusInternalServerError)
		return
	}
	defer client.Disconnect(context.Background())

	collection := client.Database(flytura.DBName).Collection(flytura.PurcharseRecordTableName)

	token := r.Header.Get("token")

	if token == "" {
		http.Error(w, "Token não fornecido", http.StatusUnauthorized)
		return
	}

	var error error
	_, error = VerifyAccessValidTokenListSheet(client, flytura.DBName, flytura.TokenAccessTableName, token)
	if error != nil {
		http.Error(w, "Token inválido", http.StatusBadRequest)
		log.Println("Token inválido", err)
		return
	}

	// Contador de atualizações
	var updatedCount int64
	//Retira os espaços das string
	re := regexp.MustCompile(`\s+`)

	for _, data := range records {

		// if data.UpdatedAt.IsZero() {
		fmt.Println(data)

		dtUpdatedAt := time.Now()
		// }

		fmt.Println("key ", data.Key)
		update := bson.M{
			"$set": bson.M{
				"status":        data.Status,
				"messageReturn": data.MessageReturn,
				"updatedAt":     dtUpdatedAt,
			},
		}

		result, err := collection.UpdateOne(context.Background(), bson.M{"key": re.ReplaceAllString(data.Key, "")}, update)
		if err == nil {
			updatedCount += result.ModifiedCount
		} else {
			fmt.Println("key err ", err)
		}
	}

	if updatedCount == 0 {
		flytura.FormataRetornoHTTP(w, "Nenhuma fatura foi atualizada", http.StatusOK)
	} else {
		flytura.FormataRetornoHTTP(w, fmt.Sprintf("%d faturas atualizadas com sucesso", updatedCount), http.StatusOK)
	}
}
