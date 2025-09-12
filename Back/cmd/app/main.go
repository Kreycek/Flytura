package main

import (
	flytura "Flytura"
	"Flytura/internal/auth"

	airLine "Flytura/internal/airLine"
	"Flytura/internal/onlyFly"
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

func main() {

	c := cors.New(cors.Options{
		AllowedOrigins: []string{flytura.UrlSite}, // Permitindo o domínio de onde vem a requisição
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

	//ONLY FLY
	/*
		Configuração criada por Ricardo Silva Ferreira
		Inicio da criação 05/09/2025 14:06
		Data Final da criação : 09/09/2025 14:10
	*/
	http.HandleFunc("/uploadOnlyFlyExcelData", onlyFly.UploadOnlyFlyHandler)
	// Não permite pesquisar por parametro apenas traz todos os registro para paginação inicialmente a primeira página
	http.HandleFunc("/GetAllOnlyFlyExcelData", onlyFly.GetAllExcelDatasHandler)
	// Obtem todos sem paginação
	http.HandleFunc("/GetOnlyFlyExcelData", onlyFly.GetAllOnlyExcelDatasHandler)
	http.HandleFunc("/GetOnlyFlyExcelDataById", onlyFly.GetExcelDataByIdHandler)
	http.HandleFunc("/InsertOnlyFlyExcelData", onlyFly.InsertExcelDataHandler)
	http.HandleFunc("/UpdateOnlyFlyExcelData", onlyFly.UpdateExcelDataHandler)
	http.HandleFunc("/VerifyExistOnlyFlyExcelData", onlyFly.VerifyExistExcelDataHandler)
	http.HandleFunc("/SearchOnlyFlyExcelData", onlyFly.SearchExcelsHandler)

	//AIRLINE
	/*
		Configuração criada por Ricardo Silva Ferreira
		Inicio da criação 09/09/2025 22:39
		Data Final da criação : 09/09/2025 22:39
	*/
	http.HandleFunc("/GetAllAirline", airLine.GetAllAirLineHandler)

	//TESTE
	http.HandleFunc("/teste", loginHandler)
	handler := c.Handler(http.DefaultServeMux)

	// balancete.GenerateBalanceteReport(2025, 01, 2025, 06)
	// http.HandleFunc("/createUser", auth.createUser) // Rota de validação do token
	// auth.CreateUser("rico", "654321")

	// Inicia o servidor na porta 8080
	log.Println("Servidor rodando na porta 8080...")
	log.Fatal(http.ListenAndServe(":8080", handler))
}
