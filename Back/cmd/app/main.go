package main

import (
	flytura "Flytura"
	"Flytura/internal/auth"
	awsS3 "Flytura/internal/awsS3"
	"Flytura/internal/db"
	"Flytura/internal/integration/bigQueryOnfly"
	"Flytura/internal/logs"
	"Flytura/internal/models"
	"Flytura/internal/outPutInvoices"
	"Flytura/internal/purcharseRecord"
	"bytes"
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	airLine "Flytura/internal/airLine"
	"Flytura/internal/perfil"
	"Flytura/internal/users"
	"log"
	"net/http"
	"net/url"

	"github.com/cenkalti/backoff/v4"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	jobRunning bool
	jobMutex   sync.Mutex
)

/*
	Função criada por Ricardo Silva Ferreira
	Inicio da criação 10/05/2026 19:36
	Data Final da criação : 10/05/2026 19:36
*/
// agenda para rodar imediatamente e depois a cada 8 horas
func startImportConciliationJob() {
	ticker := time.NewTicker(8 * time.Hour)
	defer ticker.Stop()

	// executa uma vez ao subir a aplicação
	runImportSafely()

	// executa a cada 8 horas
	for range ticker.C {
		runImportSafely()
	}
}

/*
	Função criada por Ricardo Silva Ferreira
	Inicio da criação 10/05/2026 19:36
	Data Final da criação : 10/05/2026 19:36

*/
// execução protegida do job
func runImportSafely() {

	// 🔒 evita rodar em paralelo
	jobMutex.Lock()
	if jobRunning {
		jobMutex.Unlock()
		log.Println("⚠️ Import já está rodando, pulando esta execução")
		return
	}
	jobRunning = true
	jobMutex.Unlock()

	defer func() {
		jobMutex.Lock()
		jobRunning = false
		jobMutex.Unlock()

		// 🛑 captura panic para NÃO derrubar a aplicação
		if r := recover(); r != nil {
			log.Printf("❌ Panic no ImportConciliationDataOnflys: %v", r)
		}
	}()

	log.Println("🚀 Iniciando ImportConciliationDataOnflys...")
	start := time.Now()

	// ✅ chamada da função SEM parâmetros e SEM retorno
	bigQueryOnfly.ImportConciliationDataOnflys()

	log.Printf("✅ Import finalizado em %s", time.Since(start))
}

/*
	Função criada por Ricardo Silva Ferreira
	Inicio da criação 03/12/2025 15:32
	Data Final da criação : 03/12/2025 15:40
	Obs: Cria Workers para subir as imagens para o S3
*/

func startWorkers(n int, dbConnection *mongo.Client) {
	for i := 0; i < n; i++ {
		go func(id int) {
			for task := range flytura.UploadChan {
				// log.Printf("[Worker %d] Processando arquivo: %s", id, task.FileName)
				err := retryUpload(task)
				if err != nil {
					dataLog := models.Log{Module: "S3", Class: "AwsS3->handler.go e s3.go", Method: "UploadS3FilesUnzipHandler", ErrorDescription: err.Error()}
					logs.InsertLog(dbConnection, flytura.DBName, flytura.LogsTableName, dataLog)
					// log.Printf("[Worker %d] Falha ao enviar %s: %v", id, task.FileName, err)
				} else {
					log.Printf("[Worker %d] Upload concluído: %s", id, task.FileName)
				}
			}
		}(i)
	}
}

/*
	Função criada por Ricardo Silva Ferreira
	Inicio da criação 03/12/2025 15:32
	Data Final da criação : 03/12/2025 15:40
	Obs: Fica tentando até enviar no caso depois de um minuto ele para e manda um erro
*/

func retryUpload(task flytura.UploadTask) error {
	operation := func() error {
		// log.Printf("[Retry] Tentando enviar arquivo: %s", task.FileName)
		return awsS3.UploadToS3Only(bytes.NewReader(task.FileContent), task.FileName)
	}

	// Configura backoff exponencial
	b := backoff.NewExponentialBackOff()
	b.InitialInterval = 2 * time.Second // primeira espera, 2 segundos
	b.MaxInterval = 10 * time.Second    // intervalo máximo entre tentativas
	b.MaxElapsedTime = 1 * time.Minute  // tempo total máximo para tentar, um minuto

	// Adiciona callback para logar cada tentativa
	notify := func(err error, t time.Duration) {
		log.Printf("[Retry] Falha ao enviar %s. Tentando novamente em %v. Erro: %v", task.FileName, t, err)
	}

	return backoff.RetryNotify(operation, b, notify)
}

func main() {

	value := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
	fmt.Println("GOOGLE_APPLICATION_CREDENTIALS:", value)
	// inicia o job em background
	go startImportConciliationJob()

	// Instante atual (UTC)
	nowUTC := time.Now().UTC()

	// Descobrir diferença de horas entre Portugal e México neste instante
	dh, err := flytura.DiffHours(nowUTC, flytura.Fuso1, flytura.Fuso2)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Diferença atual (Portugal - México City): %d horas\n", dh)

	errd := db.ConnectGlobalMongoDB(flytura.ConectionString)
	if errd != nil {
		log.Fatalf("Erro ao conectar ao MongoDB: %v", errd)
	}

	flytura.UploadChan = make(chan flytura.UploadTask, 200)
	startWorkers(10, db.MongoClient)

	// Descobrir IP público (fallback para localhost)
	publicIP, err := flytura.GetPublicIP()
	if err != nil {
		fmt.Println("Erro ao obter IP público:", err)
		publicIP = "localhost"
	}
	fmt.Println("Ip Publico", publicIP)

	// -------------------------------------------------------------
	// CORS: aceitar http/https dos hosts permitidos (prod/homol/dev)
	// -------------------------------------------------------------
	allowedHosts := map[string]struct{}{}

	// Helper: adiciona host mesmo se a entrada for uma URL completa
	addAllowedHost := func(urlOrHost string) {
		fmt.Println("Inicio addAllowedHost ", urlOrHost)

		s := strings.TrimSpace(urlOrHost)
		if s == "" {
			return
		}
		if u, err := url.Parse(s); err == nil && u.Host != "" {
			// Entrada era URL; usa apenas host[:port]
			allowedHosts[u.Host] = struct{}{}
			return
		}
		// Entrada era host:porta (ou IP)
		allowedHosts[s] = struct{}{}

		fmt.Println("Inicio addAllowedHost allowedHosts[s] ", allowedHosts[s])
	}

	addAllowedHost(flytura.UrlSiteExterno)
	addAllowedHost(flytura.UrlSiteHomol)
	addAllowedHost(flytura.UrlSiteProduction)
	addAllowedHost(flytura.UrlSiteLocalHost)
	addAllowedHost(flytura.UrlSiteSubDominio)
	addAllowedHost(flytura.UrlSiteLocalHostIP)

	// Decidir ambiente (melhor usar APP_ENV; mas se quiser manter o IP público):
	// if publicIP == flytura.UrlSiteProduction {
	// 	fmt.Println("Prod")
	// 	flytura.Environment = "Prod"
	// 	addAllowedHost(flytura.UrlSiteProduction) // "54.156.244.197"
	// } else if publicIP == flytura.UrlSiteHomol {
	// 	flytura.Environment = "Homol"
	// 	fmt.Println("Homol")
	// 	addAllowedHost(flytura.UrlSiteHomol) // "18.210.18.180"
	// }

	// Opcional: permitir hosts via env
	// if envHosts := strings.TrimSpace(os.Getenv("ALLOWED_HOSTS")); envHosts != "" {
	// 	for _, h := range strings.Split(envHosts, ",") {
	// 		addAllowedHost(h)
	// 	}
	// }

	// Função de validação de Origin (CORS)
	originAllowed := func(origin string) bool {
		if origin == "" {
			return false
		}
		u, err := url.Parse(origin)
		if err != nil {
			return false
		}
		if u.Scheme != "http" && u.Scheme != "https" {
			return false
		}
		if u.Host == "" {
			return false
		}

		_, ok := allowedHosts[u.Host]
		return ok
	}

	// fmt.Println("Teste originAllowed(flytura.com):", originAllowed("https://flytura.com"))

	c := cors.New(cors.Options{
		AllowOriginFunc:  originAllowed,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})

	// Configura as rotas para autenticação e validação de token
	http.HandleFunc("/api/login", auth.VerifyUser)       // Rota de login (gera o JWT)
	http.HandleFunc("/api/validate", auth.ValidateToken) // Rota de validação do token
	http.HandleFunc("/api/getPerfis", perfil.GetAllPerfilsHandler)

	//USUÁRIOS
	http.HandleFunc("/api/addUser", users.InsertUserHandler)
	http.HandleFunc("/api/getAllUsers", users.GetAllUsersHandler)
	http.HandleFunc("/api/verifyExistUser", users.VerifyExistUser)
	http.HandleFunc("/api/searchUsers", users.SearchUsersHandler)
	http.HandleFunc("/api/getUserById", users.GetUserByIdHandler)
	http.HandleFunc("/api/updateUser", users.UpdateUserHandler)

	//BIGQUERY
	/*
		Configuração criada por Ricardo Silva Ferreira
		Inicio da criação 22/05/2026 15:35
		Data Final da criação : 22/05/2026 15:35
	*/

	http.HandleFunc("/api/SearchConciliationPagination", bigQueryOnfly.SearchConciliationPaginationHandler)
	http.HandleFunc("/api/SearchConciliationExcel", bigQueryOnfly.SearchConciliationExcelHandler)

	//PURCHARSE RECORD
	/*
		Configuração criada por Ricardo Silva Ferreira
		Inicio da criação 05/09/2025 14:06
		Data Final da criação : 09/09/2025 14:10
	*/
	http.HandleFunc("/api/UploadExcelPurcharseRecord", purcharseRecord.UploadPurcharseRecordHandler)
	// Não permite pesquisar por parametro apenas traz todos os registro para paginação inicialmente a primeira página
	http.HandleFunc("/api/GetAllPurcharseRecordPagination", purcharseRecord.GetAllPurcharseRecordPaginationHandler)
	// Obtem todos sem paginação
	http.HandleFunc("/api/GetAllPurcharseRecord", purcharseRecord.GetAllPurcharseRecordHandler)
	http.HandleFunc("/api/GetPurcharseRecordById", purcharseRecord.GetPurcharseRecordByIdHandler)
	http.HandleFunc("/api/InsertPurcharseRecord", purcharseRecord.InsertPurcharseRecordHandler)
	http.HandleFunc("/api/UpdatePurcharseRecord", purcharseRecord.UpdatePurcharseRecordHandler)
	http.HandleFunc("/api/VerifyExistPurcharseRecord", purcharseRecord.VerifyExistPurcharseRecordHandler)
	http.HandleFunc("/api/SearchPurcharseRecord", purcharseRecord.SearchPurcharseRecordHandler)
	//Inserido em 09/04/2026 08:52
	http.HandleFunc("/api/SearchPurcharseRecordByPeriod", purcharseRecord.SearchPurcharseRecordByPeriodHandler)

	http.HandleFunc("/api/UpdatePurcharseRecordMultiple", purcharseRecord.UpdatePurcharseRecordMultipleHandler)
	//Inserido em 07/01/2026 20:58
	http.HandleFunc("/api/InsertPurcharseRecordMultiple", purcharseRecord.InsertPurcharseRecordMultipleHandler)
	//Inserido em 27/11/2025 13:21
	doibi := mux.NewRouter()
	doibi.HandleFunc("/api/DeletePurcharseRecordByID/{id}", purcharseRecord.DeletePurcharseRecordByIDHandler).Methods("DELETE")
	//Inicio da criação 30/11/2025 17:13
	http.HandleFunc("/api/GetAllImportStatus", purcharseRecord.GetAllImportStatussHandler)
	//Inicio da criação 01/10/2025 16:59
	http.HandleFunc("/api/GroupByCompanyName", purcharseRecord.GroupByCompanyNameHandler)

	//Inicio da criação 27/04/2026 15:59
	http.HandleFunc("/api/AgregateByImportDateAndStatus", purcharseRecord.AgregateByImportDateAndStatusHandler)

	//AIRLINE
	/*
		Configuração criada por Ricardo Silva Ferreira
		Inicio da criação 09/09/2025 22:39
		Data Final da criação : 09/09/2025 22:39
	*/
	http.HandleFunc("/api/GetAllAirline", airLine.GetAllAirLineHandler)

	/*
		Configuração criada por Ricardo Silva Ferreira
		Inicio da criação 15/03/2026 23:51
		Data Final da criação : 15/03/2026 23:52
	*/
	http.HandleFunc("/api/GetAirlineByCode", airLine.GetAirLineByCodeHandler)

	//OUTPUTINVOICES
	/*
		Configuração criada por Ricardo Silva Ferreira
		Inicio da criação 11/11/2025 11:50
		Data Final da criação : 11/11/2025 11:55
	*/

	http.HandleFunc("/api/SearchOutPutInvoices", outPutInvoices.SearchOutPutInvoicesHandler)

	//Inserido em 09/04/2026 23:38
	http.HandleFunc("/api/SearchOutPutInvoicesInforme", outPutInvoices.SearchOutPutInvoicesInformeHandler)

	http.HandleFunc("/api/GroupByCompanySumSection", outPutInvoices.GroupByCompanySumSectionHandler)

	/*
		Configuração criada por Ricardo Silva Ferreira
		Inicio da criação 10/06/2026 11:32
		Data Final da criação : 10/06/2026 11:36
	*/
	http.HandleFunc("/api/UploadExcelOutputInvoicesRecord", outPutInvoices.UploadExcelOutputInvoicesRecordHandler)

	//Inserido em 20/11/2025 17:26
	r := mux.NewRouter()
	r.HandleFunc("/api/DeleteOutPutInvoicesByID/{id}", outPutInvoices.DeleteOutPutInvoicesByIDHandler).Methods("DELETE")

	//AMAZON S3f
	/*
		Configuração criada por Ricardo Silva Ferreira
		Inicio da criação 17/10/2025 13:10
		Data Final da criação : 17/10/2025 13:15
	*/

	//Inserido em 10/06/2026 21:34
	http.HandleFunc("/api/UploadManualImportInvoicesRecord", awsS3.UploadManualImportInvoicesRecordHandler)
	http.HandleFunc("/api/SearchS3ImagesDBPagination", awsS3.SearchS3ImagesDBPaginationHandler)
	http.HandleFunc("/api/SearchS3ImagesDBFull", awsS3.SearchS3ImagesDBFullHandler)
	http.HandleFunc("/api/UpdateStatusS3Image", awsS3.UpdateStatusS3ImageHandler)
	http.HandleFunc("/api/UpdateMultipleStatusS3Images", awsS3.UpdateMultipleStatusS3ImagesHandler)
	http.HandleFunc("/api/UpdateStatusPdforXml", awsS3.UpdateStatusPdfOrXmlHandler)
	//Inserido em 25/11/2025 16:24
	s := mux.NewRouter()
	s.HandleFunc("/api/DeleteImagesDBByIDHandler", awsS3.DeleteImagesDBByIDHandler).Methods("DELETE")
	//API'S PUBLICAS
	http.HandleFunc("/api/GetPurcharseRecordByStatus", purcharseRecord.GetPurcharseRecordStatusHandler)
	//FUNÇÃO ABAIXO NÃO USADA
	http.HandleFunc("/api/UploadS3Files", awsS3.UploadS3FilesHandler)
	//FUNÇÃO ABAIXO NÃO USADA
	http.HandleFunc("/api/UploadS3FilesUnzip", awsS3.UploadS3FilesUnzipHandler)

	//Inserido em 11/06/2026 13:06
	http.HandleFunc("/api/CheckDuplicatePDFs", awsS3.CheckDuplicatePDFsHandler)

	http.HandleFunc("/api/UploadS3MultiplesFilesUnzip", awsS3.UploadS3MultiplesFilesUnzipHandler)
	http.HandleFunc("/api/InsertOutPutInvoices", outPutInvoices.InsertOutPutInvoicesHandler)

	//TESTE
	http.HandleFunc("/api/teste", flytura.LoginHandler)

	// Junta mux com DefaultServeMux
	mainHandler := http.NewServeMux()
	mainHandler.Handle("/api/", http.DefaultServeMux)       // rotas antigas
	mainHandler.Handle("/api/DeleteOutPutInvoicesByID/", r) // rotas mux
	mainHandler.Handle("/api/DeleteImagesDBByIDHandler", s)
	mainHandler.Handle("/api/DeletePurcharseRecordByID/", doibi)

	handler := c.Handler(mainHandler)

	// balancete.GenerateBalanceteReport(2025, 01, 2025, 06)
	// http.HandleFunc("/createUser", auth.createUser) // Rota de validação do token
	// auth.CreateUser("rico", "654321")

	// token, _ := flytura.GerarTokenSemExpiracao()
	// fmt.Println("Token para o cliente ", token)

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
