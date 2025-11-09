package main

import (
	flytura "Flytura"
	"Flytura/internal/auth"
	"Flytura/internal/awsS3"
	"Flytura/internal/db"
	"Flytura/internal/purcharseRecord"
	"context"
	"fmt"
	"io"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	airLine "Flytura/internal/airLine"
	"Flytura/internal/perfil"
	"Flytura/internal/users"
	"encoding/json"
	"log"
	"net/http"

	"github.com/rs/cors"
)

func loginHandler(w http.ResponseWriter, r *http.Request) {
	// Definir um mapa com os dados
	data := map[string]string{
		"name": "ricardo",
	}

	// Definir o header como JSON
	w.Header().Set("Content-Type", "application/json")

	// Codificar o mapa de dados para JSON e enviar como resposta
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(data)
}

func getPublicIP() (string, error) {
	resp, err := http.Get("https://api.ipify.org")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	ip, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(ip), nil
}

func main() {

	err := db.ConnectGlobalMongoDB(flytura.ConectionString)
	if err != nil {
		log.Fatalf("Erro ao conectar ao MongoDB: %v", err)
	}

	publicIP, err := getPublicIP()
	if err != nil {
		fmt.Println("Erro ao obter IP público:", err)
		publicIP = "localhost"
	}

	fmt.Println("Ip Publico", publicIP)

	var allowedOrigin string
	if strings.Contains(publicIP, flytura.UrlSiteProduction) {
		allowedOrigin = "http://" + flytura.UrlSiteProduction // domínio de produção
	} else if strings.Contains(publicIP, flytura.UrlSiteHomol) {
		allowedOrigin = "http://" + flytura.UrlSiteHomol // domínio de produção
	} else {
		allowedOrigin = "http://" + flytura.UrlSiteLocalHost // ou a porta que seu front usa
	}

	c := cors.New(cors.Options{
		AllowedOrigins: []string{allowedOrigin}, // Permitindo o domínio de onde vem a requisição
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Content-Type", "Authorization"},
	})
	// Configura as rotas para autenticação e validação de token
	http.HandleFunc("/login", auth.VerifyUser)       // Rota de login (gera o JWT)
	http.HandleFunc("/validate", auth.ValidateToken) // Rota de validação do token
	http.HandleFunc("/getPerfis", perfil.GetAllPerfilsHandler)

	//USUÁRIOS
	http.HandleFunc("/addUser", users.InsertUserHandler)
	http.HandleFunc("/getAllUsers", users.GetAllUsersHandler)
	http.HandleFunc("/verifyExistUser", users.VerifyExistUser)
	http.HandleFunc("/searchUsers", users.SearchUsersHandler)
	http.HandleFunc("/getUserById", users.GetUserByIdHandler)
	http.HandleFunc("/updateUser", users.UpdateUserHandler)

	//PURCHARSE RECORD
	/*
		Configuração criada por Ricardo Silva Ferreira
		Inicio da criação 05/09/2025 14:06
		Data Final da criação : 09/09/2025 14:10
	*/
	http.HandleFunc("/UploadPurcharseRecord", purcharseRecord.UploadPurcharseRecordHandler)
	// Não permite pesquisar por parametro apenas traz todos os registro para paginação inicialmente a primeira página
	http.HandleFunc("/GetAllPurcharseRecordPagination", purcharseRecord.GetAllPurcharseRecordPaginationHandler)
	// Obtem todos sem paginação
	http.HandleFunc("/GetAllPurcharseRecord", purcharseRecord.GetAllPurcharseRecordHandler)
	http.HandleFunc("/GetPurcharseRecordById", purcharseRecord.GetPurcharseRecordByIdHandler)
	http.HandleFunc("/InsertPurcharseRecord", purcharseRecord.InsertPurcharseRecordHandler)
	http.HandleFunc("/UpdatePurcharseRecord", purcharseRecord.UpdatePurcharseRecordHandler)
	http.HandleFunc("/VerifyExistPurcharseRecord", purcharseRecord.VerifyExistPurcharseRecordHandler)
	http.HandleFunc("/SearchPurcharseRecord", purcharseRecord.SearchPurcharseRecordHandler)
	http.HandleFunc("/UpdatePurcharseRecordMultiple", purcharseRecord.UpdatePurcharseRecordMultipleHandler)

	//Inicio da criação 30/11/2025 17:13
	http.HandleFunc("/GetAllImportStatus", purcharseRecord.GetAllImportStatussHandler)

	//Inicio da criação 01/10/2025 16:59
	http.HandleFunc("/GroupByCompanyName", purcharseRecord.GroupByCompanyNameHandler)

	//AIRLINE
	/*
		Configuração criada por Ricardo Silva Ferreira
		Inicio da criação 09/09/2025 22:39
		Data Final da criação : 09/09/2025 22:39
	*/
	http.HandleFunc("/GetAllAirline", airLine.GetAllAirLineHandler)

	//AMAZON S3f
	/*
		Configuração criada por Ricardo Silva Ferreira
		Inicio da criação 17/10/2025 13:10
		Data Final da criação : 17/10/2025 13:15
	*/

	http.HandleFunc("/SearchS3ImagesDBPagination", awsS3.SearchS3ImagesDBPaginationHandler)
	http.HandleFunc("/SearchS3ImagesDBFull", awsS3.SearchS3ImagesDBFullHandler)
	http.HandleFunc("/UpdateStatusS3Image", awsS3.UpdateStatusS3ImageHandler)
	http.HandleFunc("/UpdateMultipleStatusS3Images", awsS3.UpdateMultipleStatusS3ImagesHandler)
	http.HandleFunc("/UpdateStatusPdforXml", awsS3.UpdateStatusPdfOrXmlHandler)

	//API'S PUBLICAS
	http.HandleFunc("/GetPurcharseRecordByStatus", purcharseRecord.GetPurcharseRecordStatusHandler)
	http.HandleFunc("/UploadS3Files", awsS3.UploadS3FilesHandler)
	http.HandleFunc("/UploadS3FilesUnzip", awsS3.UploadS3FilesUnzipHandler)
	http.HandleFunc("/UploadS3MultiplesFilesUnzip", awsS3.UploadS3MultiplesFilesUnzipHandler)

	//TESTE
	http.HandleFunc("/teste", loginHandler)
	handler := c.Handler(http.DefaultServeMux)

	// balancete.GenerateBalanceteReport(2025, 01, 2025, 06)
	// http.HandleFunc("/createUser", auth.createUser) // Rota de validação do token
	// auth.CreateUser("rico", "654321")

	// Inicia o servidor na porta 8080

	go func() {
		log.Println("Servidor iniciado na porta 8080")
		if err := http.ListenAndServe(":8080", handler); err != nil {
			log.Fatalf("Erro ao iniciar servidor: %v", err)
		}
	}()

	// Canal para escutar sinais do sistema
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	<-stop
	log.Println("Encerrando servidor...")

	// Desconectar do MongoDB

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := db.MongoClient.Disconnect(ctx); err != nil {
		log.Printf("Erro ao desconectar do MongoDB: %v", err)
	}

	// Testar se a conexão ainda está ativa
	err = db.MongoClient.Ping(ctx, nil)
	if err != nil {
		log.Printf("Conexão encerrada corretamente: %v", err)
	} else {
		log.Println("⚠️ A conexão ainda está ativa!")
	}

}
