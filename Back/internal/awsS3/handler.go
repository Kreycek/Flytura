package awsS3

import (
	flytura "Flytura"
	"Flytura/internal/airLine"
	"Flytura/internal/db"
	"Flytura/internal/models"
	"archive/zip"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

/*
Função criada por Ricardo Silva Ferreira
Início da criação: 17/10/2025 13:02
Data final da criação:  17/10/2025 13:00
Última Data modificação:  19/11/2025 16:14, adicionado campo billedFlytura
Última Data modificação:  03/12/2025 17:36, adicionado o workProcess com go routines
OBS: Função que enviar os dados para o S3 E grava no banco de dados, ela apenas importa um arquivo seja qual for
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

	// Lê o conteúdo do ZIP para memória
	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, file)
	if err != nil {
		http.Error(w, "Erro ao copiar conteúdo do ZIP", http.StatusInternalServerError)
		return
	}

	companyCode := r.FormValue("companyCode")
	key := r.FormValue("key")

	flytura.UploadChan <- flytura.UploadTask{
		FileContent: buf.Bytes(),
		FileName:    header.Filename,
		CompanyCode: companyCode,
		Key:         key,
	}

	// err = UploadToS3Only(file, header.Filename, companyCode, key)
	// if err != nil {

	// 	fmt.Println("Erro ao enviar para S3:", err)
	// 	// return fmt.Errorf("erro ao enviar para S3: %w", err)

	// 	return
	// }

	// clientDb, err := db.ConnectMongoDB(flytura.ConectionString)
	// if err != nil {
	// 	http.Error(w, "Erro ao conectar ao MongoDB", http.StatusInternalServerError)
	// 	return
	// }
	// defer clientDb.Disconnect(context.Background())

	airLineData, errAirLineName := airLine.GetAirLineFileName(db.MongoClient, flytura.DBName, "airline", companyCode)
	if errAirLineName != nil {
		log.Println("Erro ao obter nome do arquivo:", errAirLineName)
	}
	companyName := airLineData["Name"].(string)

	nowUTC := time.Now().UTC()

	dh, err := flytura.DiffHours(nowUTC, flytura.Fuso1, flytura.Fuso2)
	if err != nil {
		panic(err)
	}

	image := models.ImagesDB{
		ID:              primitive.NewObjectID(),
		FileName:        header.Filename,
		DtImport:        nowUTC.Add(-time.Duration(dh) * time.Hour),
		PDFFileName:     "",
		XMLFileName:     "",
		CompanyCode:     companyCode,
		CompanyName:     companyName,
		DownloadDone:    false,
		DownloadPDFDone: false,
		DownloadXMLDone: false,
		Key:             key,
		ZipFileName:     "",
		OriginData:      "Automation",
	}

	InsertIMGS3(db.MongoClient, flytura.DBName, "imagesDB", image)

	// Retornar resposta JSON
	response := map[string]any{
		"Result": "Arquivo %s enviado com sucesso!",
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("erro ao codificar resposta JSON: %v", err)
	}
}

/*
Função criada por Ricardo Silva Ferreira
Início da criação: 23/10/2025 01:03
Data final da criação:  23/10/2025 01:57
Última Data modificação:  19/11/2025 16:14, adicionado campo billedFlytura
Última Data modificação:  03/12/2025 17:36, adicionado o workProcess com go routines
OBS: Função recebe um arquivo .ZIP e que descompacta arquivos envia vários arquivos para o AWS S3
*/

func UploadS3FilesUnzipHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseMultipartForm(50 << 20) // até 50 MB
	if err != nil {
		http.Error(w, "Erro ao processar formulário", http.StatusBadRequest)
		return
	}

	companyCode := r.FormValue("companyCode")
	key := r.FormValue("key")

	zipFile, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Erro ao ler arquivo ZIP", http.StatusBadRequest)
		return
	}
	defer zipFile.Close()

	// fmt.Println("FileName:", header.Filename)

	// Lê o conteúdo do ZIP para memória
	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, zipFile)
	if err != nil {
		http.Error(w, "Erro ao copiar conteúdo do ZIP", http.StatusInternalServerError)
		return
	}

	// Upload do arquivo ZIP original para o S3
	// err = UploadToS3Only(bytes.NewReader(buf.Bytes()), header.Filename, companyCode, key)
	// if err != nil {
	// 	http.Error(w, "Erro ao enviar ZIP original para o S3", http.StatusInternalServerError)
	// 	fmt.Println("Erro ao enviar ZIP original:", err)
	// 	return
	// }

	flytura.UploadChan <- flytura.UploadTask{
		FileContent: buf.Bytes(),
		FileName:    header.Filename,
		CompanyCode: companyCode,
		Key:         key,
	}

	// Abre o ZIP para leitura dos arquivos internos
	zipReader, err := zip.NewReader(bytes.NewReader(buf.Bytes()), int64(buf.Len()))
	if err != nil {
		http.Error(w, "Erro ao abrir arquivo ZIP", http.StatusInternalServerError)
		return
	}

	pdfFileName := ""
	xmlFileName := ""
	zipFileName := ""

	for _, file := range zipReader.File {
		if file.FileInfo().IsDir() {
			continue
		}

		f, err := file.Open()
		if err != nil {
			fmt.Printf("Erro ao abrir %s: %v\n", file.Name, err)
			continue
		}

		fileContent, err := io.ReadAll(f)
		f.Close()
		if err != nil {
			fmt.Printf("Erro ao ler conteúdo de %s: %v\n", file.Name, err)
			continue
		}

		// Detecta extensão com segurança
		ext := strings.ToLower(filepath.Ext(file.Name))
		switch ext {
		case ".xml":
			xmlFileName = file.Name
		case ".pdf":
			pdfFileName = file.Name
		case ".zip":
			zipFileName = file.Name
		}

		flytura.UploadChan <- flytura.UploadTask{
			FileContent: fileContent,
			FileName:    file.Name,
			CompanyCode: companyCode,
			Key:         key,
		}

		// err = UploadToS3Only(bytes.NewReader(fileContent), file.Name, companyCode, key)
		// if err != nil {
		// 	fmt.Printf("Erro ao enviar %s para S3: %v\n", file.Name, err)
		// 	continue
		// }
	}

	// Conectar ao MongoDB
	// clientDb, err := db.ConnectMongoDB(flytura.ConectionString)
	// if err != nil {
	// 	http.Error(w, "Erro ao conectar ao MongoDB", http.StatusInternalServerError)
	// 	return
	// }
	// defer clientDb.Disconnect(context.Background())

	airLineData, errAirLineName := airLine.GetAirLineFileName(db.MongoClient, flytura.DBName, "airline", companyCode)
	if errAirLineName != nil {
		log.Println("Erro ao obter nome do arquivo:", errAirLineName)
	}

	nowUTC := time.Now().UTC()

	dh, err := flytura.DiffHours(nowUTC, flytura.Fuso1, flytura.Fuso2)
	if err != nil {
		panic(err)
	}
	companyName := airLineData["Name"].(string)

	image := models.ImagesDB{
		ID:              primitive.NewObjectID(),
		FileName:        header.Filename,
		DtImport:        nowUTC.Add(-time.Duration(dh) * time.Hour),
		PDFFileName:     pdfFileName,
		XMLFileName:     xmlFileName,
		CompanyCode:     companyCode,
		CompanyName:     companyName,
		DownloadDone:    false,
		DownloadPDFDone: false,
		DownloadXMLDone: false,
		Key:             key,
		ZipFileName:     zipFileName,
		OriginData:      "Automation",
	}

	InsertIMGS3(db.MongoClient, flytura.DBName, "imagesDB", image)

	// Retornar resposta JSON
	response := map[string]any{
		"Result": "Arquivos importados com sucesso",
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("erro ao codificar resposta JSON: %v", err)
	}
}

/*
Função criada por Ricardo Silva Ferreira
Início da criação: 23/10/2025 13:14
Data final da criação:  23/10/2025 13:17
Última Data modificação:  19/11/2025 16:14, adicionado campo billedFlytura
Última Data modificação:  03/12/2025 17:36, adicionado o workProcess com go routines
OBS: Função recebe um arquivo .ZIP e que descompacta arquivos envia vários arquivos para o AWS S3
*/

func UploadS3MultiplesFilesUnzipHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseMultipartForm(50 << 20) // até 50 MB
	if err != nil {
		http.Error(w, "Erro ao processar formulário", http.StatusBadRequest)
		return
	}

	companyCode := r.FormValue("companyCode")

	//O Parâmetro abaixo diz se foi processado pela flytura
	billedFlytura, erroBF := strconv.ParseBool(r.FormValue("billedFlytura"))
	if erroBF != nil {
		http.Error(w, "Erro ao enviar parâmetro billedFlytura", http.StatusBadRequest)
		return
	}

	key := r.FormValue("key")

	files := r.MultipartForm.File["file"]
	if len(files) == 0 {
		http.Error(w, "Nenhum arquivo ZIP enviado", http.StatusBadRequest)
		return
	}

	var importedFiles []string
	var filesNotZip []string
	var filesName []string

	for _, header := range files {
		zipFile, err := header.Open()
		if err != nil {
			fmt.Printf("Erro ao abrir %s: %v\n", header.Filename, err)
			continue
		}
		defer zipFile.Close()

		// fmt.Println("Processando:", header.Filename)

		buf := new(bytes.Buffer)
		_, err = io.Copy(buf, zipFile)
		if err != nil {
			fmt.Printf("Erro ao copiar conteúdo de %s: %v\n", header.Filename, err)
			continue
		}

		// err = UploadToS3Only(bytes.NewReader(buf.Bytes()), header.Filename, companyCode, key)

		// if err != nil {
		// 	fmt.Printf("Erro ao enviar ZIP %s para S3: %v\n", header.Filename, err)
		// 	continue
		// }

		flytura.UploadChan <- flytura.UploadTask{
			FileContent: buf.Bytes(),
			FileName:    header.Filename,
			CompanyCode: companyCode,
			Key:         key,
		}

		//Aqui pega o arquivo zip
		zipReader, err := zip.NewReader(bytes.NewReader(buf.Bytes()), int64(buf.Len()))
		if err != nil {
			fmt.Printf("Erro ao abrir ZIP %s: %v\n", header.Filename, err)
			continue
		}

		pdfFileName := ""
		xmlFileName := ""
		zipFileName := ""

		ext := strings.ToLower(filepath.Ext(header.Filename))

		if ext == ".zip" {
			zipFileName = header.Filename
		} else {
			filesNotZip = append(filesNotZip, "O Arquivo "+header.Filename+" não foi importado porque não é um arquivo .zip")
			continue
		}

		for _, file := range zipReader.File {
			if file.FileInfo().IsDir() {
				continue
			}

			f, err := file.Open()
			if err != nil {
				fmt.Printf("Erro ao abrir %s: %v\n", file.Name, err)
				continue
			}

			fileContent, err := io.ReadAll(f)
			f.Close()
			if err != nil {
				fmt.Printf("Erro ao ler conteúdo de %s: %v\n", file.Name, err)
				continue
			}

			ext := strings.ToLower(filepath.Ext(file.Name))
			switch ext {
			case ".xml":
				xmlFileName = file.Name
			case ".pdf":
				pdfFileName = file.Name
				filesName = append(filesName, pdfFileName)
			case ".zip":
				zipFileName = file.Name
			}

			// err = UploadToS3Only(bytes.NewReader(fileContent), file.Name, companyCode, key)
			// if err != nil {
			// 	fmt.Printf("Erro ao enviar %s para S3: %v\n", file.Name, err)
			// 	continue
			// }

			//Aqui envia para o S3
			flytura.UploadChan <- flytura.UploadTask{
				FileContent: fileContent,
				FileName:    file.Name,
				CompanyCode: companyCode,
				Key:         key,
			}

		}

		// clientDb, err := db.ConnectMongoDB(flytura.ConectionString)
		// if err != nil {
		// 	fmt.Printf("Erro ao conectar ao MongoDB: %v\n", err)
		// 	continue
		// }
		// defer clientDb.Disconnect(context.Background())

		airLineData, errAirLineName := airLine.GetAirLineFileName(db.MongoClient, flytura.DBName, "airline", companyCode)
		if errAirLineName != nil {
			fmt.Printf("Erro ao obter nome da companhia: %v\n", errAirLineName)
			continue
		}

		companyName := airLineData["Name"].(string)

		nowUTC := time.Now().UTC()

		dh, err := flytura.DiffHours(nowUTC, flytura.Fuso1, flytura.Fuso2)
		if err != nil {
			panic(err)
		}

		image := models.ImagesDB{
			ID:              primitive.NewObjectID(),
			FileName:        strings.Join(filesName, ";"),
			DtImport:        nowUTC.Add(-time.Duration(dh) * time.Hour),
			PDFFileName:     strings.Join(filesName, ";"),
			XMLFileName:     xmlFileName,
			CompanyCode:     companyCode,
			CompanyName:     companyName,
			DownloadPDFDone: false,
			DownloadXMLDone: false,
			DownloadDone:    false,
			Key:             key,
			ZipFileName:     zipFileName,
			BilledFlytura:   billedFlytura,
			OriginData:      "Automation",
		}

		InsertIMGS3(db.MongoClient, flytura.DBName, "imagesDB", image)
		importedFiles = append(importedFiles, header.Filename)
	}

	response := map[string]any{
		"Result":           "Arquivos importados com sucesso",
		"ImportedFiles":    importedFiles,
		"NotImportedFiles": filesNotZip,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		fmt.Printf("Erro ao codificar resposta JSON: %v\n", err)
	}

}

/*
Função criada por Ricardo Silva Ferreira
Início da criação: 23/10/2025 13:14
Data final da criação:  23/10/2025 13:17
Última Data modificação:  19/11/2025 16:14, adicionado campo billedFlytura
Última Data modificação:  03/12/2025 17:36, adicionado o workProcess com go routines
OBS: Função recebe um arquivo .ZIP e que descompacta arquivos envia vários arquivos para o AWS S3
*/

// func UploadS3MultiplesFilesUnzipHandler(w http.ResponseWriter, r *http.Request) {

// 	if r.Method != http.MethodPost {
// 		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
// 		return
// 	}

// 	err := r.ParseMultipartForm(50 << 20) // até 50 MB
// 	if err != nil {
// 		http.Error(w, "Erro ao processar formulário", http.StatusBadRequest)
// 		return
// 	}

// 	companyCode := r.FormValue("companyCode")

// 	//O Parâmetro abaixo diz se foi processado pela flytura
// 	billedFlytura, erroBF := strconv.ParseBool(r.FormValue("billedFlytura"))
// 	if erroBF != nil {
// 		http.Error(w, "Erro ao enviar parâmetro billedFlytura", http.StatusBadRequest)
// 		return
// 	}

// 	key := r.FormValue("key")

// 	files := r.MultipartForm.File["file"]
// 	if len(files) == 0 {
// 		http.Error(w, "Nenhum arquivo ZIP enviado", http.StatusBadRequest)
// 		return
// 	}

// 	var importedFiles []string
// 	var filesNotZip []string

// 	for _, header := range files {
// 		zipFile, err := header.Open()
// 		if err != nil {
// 			fmt.Printf("Erro ao abrir %s: %v\n", header.Filename, err)
// 			continue
// 		}
// 		defer zipFile.Close()

// 		// fmt.Println("Processando:", header.Filename)

// 		buf := new(bytes.Buffer)
// 		_, err = io.Copy(buf, zipFile)
// 		if err != nil {
// 			fmt.Printf("Erro ao copiar conteúdo de %s: %v\n", header.Filename, err)
// 			continue
// 		}

// 		// err = UploadToS3Only(bytes.NewReader(buf.Bytes()), header.Filename, companyCode, key)

// 		// if err != nil {
// 		// 	fmt.Printf("Erro ao enviar ZIP %s para S3: %v\n", header.Filename, err)
// 		// 	continue
// 		// }

// 		flytura.UploadChan <- flytura.UploadTask{
// 			FileContent: buf.Bytes(),
// 			FileName:    header.Filename,
// 			CompanyCode: companyCode,
// 			Key:         key,
// 		}

// 		//Aqui pega o arquivo zip
// 		zipReader, err := zip.NewReader(bytes.NewReader(buf.Bytes()), int64(buf.Len()))
// 		if err != nil {
// 			fmt.Printf("Erro ao abrir ZIP %s: %v\n", header.Filename, err)
// 			continue
// 		}

// 		pdfFileName := ""
// 		xmlFileName := ""
// 		zipFileName := ""

// 		ext := strings.ToLower(filepath.Ext(header.Filename))

// 		if ext == ".zip" {
// 			zipFileName = header.Filename
// 		} else {
// 			filesNotZip = append(filesNotZip, "O Arquivo "+header.Filename+" não foi importado porque não é um arquivo .zip")
// 			continue
// 		}

// 		for _, file := range zipReader.File {
// 			if file.FileInfo().IsDir() {
// 				continue
// 			}

// 			f, err := file.Open()
// 			if err != nil {
// 				fmt.Printf("Erro ao abrir %s: %v\n", file.Name, err)
// 				continue
// 			}

// 			fileContent, err := io.ReadAll(f)
// 			f.Close()
// 			if err != nil {
// 				fmt.Printf("Erro ao ler conteúdo de %s: %v\n", file.Name, err)
// 				continue
// 			}

// 			ext := strings.ToLower(filepath.Ext(file.Name))
// 			switch ext {
// 			case ".xml":
// 				xmlFileName = file.Name
// 			case ".pdf":
// 				pdfFileName = file.Name
// 			case ".zip":
// 				zipFileName = file.Name
// 			}

// 			// err = UploadToS3Only(bytes.NewReader(fileContent), file.Name, companyCode, key)
// 			// if err != nil {
// 			// 	fmt.Printf("Erro ao enviar %s para S3: %v\n", file.Name, err)
// 			// 	continue
// 			// }

// 			//Aqui envia para o S3
// 			flytura.UploadChan <- flytura.UploadTask{
// 				FileContent: fileContent,
// 				FileName:    file.Name,
// 				CompanyCode: companyCode,
// 				Key:         key,
// 			}

// 		}

// 		// clientDb, err := db.ConnectMongoDB(flytura.ConectionString)
// 		// if err != nil {
// 		// 	fmt.Printf("Erro ao conectar ao MongoDB: %v\n", err)
// 		// 	continue
// 		// }
// 		// defer clientDb.Disconnect(context.Background())

// 		airLineData, errAirLineName := airLine.GetAirLineFileName(db.MongoClient, flytura.DBName, "airline", companyCode)
// 		if errAirLineName != nil {
// 			fmt.Printf("Erro ao obter nome da companhia: %v\n", errAirLineName)
// 			continue
// 		}

// 		companyName := airLineData["Name"].(string)

// 		nowUTC := time.Now().UTC()

// 		dh, err := flytura.DiffHours(nowUTC, flytura.Fuso1, flytura.Fuso2)
// 		if err != nil {
// 			panic(err)
// 		}

// 		image := models.ImagesDB{
// 			ID:              primitive.NewObjectID(),
// 			FileName:        header.Filename,
// 			DtImport:        nowUTC.Add(-time.Duration(dh) * time.Hour),
// 			PDFFileName:     pdfFileName,
// 			XMLFileName:     xmlFileName,
// 			CompanyCode:     companyCode,
// 			CompanyName:     companyName,
// 			DownloadPDFDone: false,
// 			DownloadXMLDone: false,
// 			DownloadDone:    false,
// 			Key:             key,
// 			ZipFileName:     zipFileName,
// 			BilledFlytura:   billedFlytura,
// 			OriginData:      "RPA",
// 		}

// 		InsertIMGS3(db.MongoClient, flytura.DBName, "imagesDB", image)
// 		importedFiles = append(importedFiles, header.Filename)
// 	}

// 	response := map[string]any{
// 		"Result":           "Arquivos importados com sucesso",
// 		"ImportedFiles":    importedFiles,
// 		"NotImportedFiles": filesNotZip,
// 	}
// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(http.StatusOK)
// 	if err := json.NewEncoder(w).Encode(response); err != nil {
// 		fmt.Printf("Erro ao codificar resposta JSON: %v\n", err)
// 	}

// }

/*
Função criada por Ricardo Silva Ferreira
Inicio da criação 19/10/2025 19:14
Data Final da criação : 19/10/2025 19:17
Alteração : 01/12/2025 20:56 a pedido do joão acrescentada busca por se fez donwload sim ou não
*/

func SearchS3ImagesDBPaginationHandler(w http.ResponseWriter, r *http.Request) {
	// Verificar se a requisição é do tipo GET

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
	// client, err := db.ConnectMongoDB(flytura.ConectionString)
	// if err != nil {
	// 	http.Error(w, "Erro ao conectar ao MongoDB", http.StatusInternalServerError)
	// 	return
	// }
	// defer client.Disconnect(context.Background())

	// Definir estrutura para receber os parâmetros

	query := r.URL.Query()

	key := query.Get("key")
	companyCode := query.Get("companyCode")
	startDateStr := query.Get("startDate")
	billedFlyturaStr := query.Get("billedFlytura")
	doDonwloadStr := query.Get("doDownload")

	var billedFlytura *bool

	if billedFlyturaStr != "" {
		parsed, err := strconv.ParseBool(billedFlyturaStr)
		if err != nil {
			http.Error(w, "Erro ao fornecer o parâmetro billedFlytura", http.StatusBadRequest)
			return
		}
		billedFlytura = &parsed // atribui o endereço do bool convertido
	}

	var doDownload *bool

	// fmt.Println("billedFlyturaStr", billedFlyturaStr)
	// fmt.Println("doDonwloadStr", doDonwloadStr)

	if doDonwloadStr != "" {
		parsed, err := strconv.ParseBool(doDonwloadStr)
		if err != nil {
			http.Error(w, "Erro ao fornecer o parâmetro billedFlytura", http.StatusBadRequest)
			return
		}
		doDownload = &parsed // atribui o endereço do bool convertido
	}

	var startDate *time.Time
	if startDateStr != "" {
		t, err := time.Parse(time.RFC3339, startDateStr)
		if err == nil {
			startDate = &t
		} else {
			// lidar com erro de parsing, se necessário
			fmt.Println("Erro ao converter startDate SearchS3ImagesDBPaginationHandler:", err)
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
			fmt.Println("Erro ao converter startDate SearchS3ImagesDBPaginationHandler 2:", err)
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

	// fmt.Println("StartDate", startDate)
	// fmt.Println("EndDate", endDateStr)
	// Definir valores padrão para paginação

	if limit < 1 {
		limit = 10
	}

	// Buscar imagens com paginação
	imagesDb, total, err := SearchImagesDBPagination(
		db.MongoClient,
		flytura.DBName,
		flytura.ImagesDBTableName,
		&companyCode,
		&key,
		billedFlytura,
		doDownload,
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

/*
Função criada por Ricardo Silva Ferreira
Inicio da criação 20/10/2025 13:31
Data Final da criação : 20/10/2025 13:32
*/

func SearchS3ImagesDBFullHandler(w http.ResponseWriter, r *http.Request) {
	// Verificar se a requisição é do tipo GET

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
	// client, err := db.ConnectMongoDB(flytura.ConectionString)
	// if err != nil {
	// 	http.Error(w, "Erro ao conectar ao MongoDB", http.StatusInternalServerError)
	// 	return
	// }
	// defer client.Disconnect(context.Background())

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
			fmt.Println("Erro ao converter startDat :", err)
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
			fmt.Println("Erro ao converter startDate SearchS3ImagesDBFullHandler:", err)
		}
	}

	// Buscar imagens com paginação
	imagesDb, total, err := SearchImagesDBFull(
		db.MongoClient,
		flytura.DBName,
		flytura.ImagesDBTableName,
		&companyCode,
		startDate,
		endDate)

	if err != nil {
		http.Error(w, "Erro ao buscar imagens", http.StatusInternalServerError)
		return
	}

	// Criar resposta JSON com paginação
	response := map[string]any{
		"total":    total,
		"imagesDB": imagesDb,
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
Inicio da criação 20/10/2025 15:02
Data Final da criação : 20/10/2025 15:05
*/

func UpdateStatusS3ImageHandler(w http.ResponseWriter, r *http.Request) {
	// Validar o token de autenticação
	status, msg := flytura.TokenValido(w, r)
	if !status {
		http.Error(w, fmt.Sprintf("Erro ao validar token: %v", msg), http.StatusUnauthorized)
		return
	}

	// Decodificar o JSON recebido
	var data models.ImagesDB
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

		nowUTC := time.Now().UTC()

		dh, err := flytura.DiffHours(nowUTC, flytura.Fuso1, flytura.Fuso2)
		if err != nil {
			panic(err)
		}

		data.UpdatedAt = nowUTC.Add(-time.Duration(dh) * time.Hour)
	}

	// Criar o objeto de atualização
	update := bson.M{
		"$set": bson.M{
			"downloadDone": data.DownloadDone,
			"updatedAt":    data.UpdatedAt,
		},
	}

	// Conectar ao MongoDB e atualizar o usuário
	// client, err := db.ConnectMongoDB(flytura.ConectionString)
	// if err != nil {
	// 	flytura.FormataRetornoHTTP(w, "Erro ao conectar ao banco de dados", http.StatusInternalServerError)

	// 	// http.Error(w, "Erro ao conectar ao banco de dados", http.StatusInternalServerError)
	// 	return
	// }
	// defer client.Disconnect(context.Background())

	collection := db.MongoClient.Database(flytura.DBName).Collection(flytura.ImagesDBTableName)
	result, err := collection.UpdateOne(context.Background(), bson.M{"_id": data.ID}, update)
	if err != nil {
		flytura.FormataRetornoHTTP(w, "Erro ao atualizar dados da imagem", http.StatusInternalServerError)

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
	flytura.FormataRetornoHTTP(w, "Imagem atualizda com sucesso", http.StatusOK)

}

/*
Função criada por Ricardo Silva Ferreira
Inicio da criação 20/10/2025 15:21
Data Final da criação : 20/10/2025 15:40
*/
func UpdateMultipleStatusS3ImagesHandler(w http.ResponseWriter, r *http.Request) {
	// Validar o token de autenticação
	status, msg := flytura.TokenValido(w, r)
	if !status {
		http.Error(w, fmt.Sprintf("Erro ao validar token: %v", msg), http.StatusUnauthorized)
		return
	}

	// Estrutura esperada no JSON
	type UpdateRequest struct {
		IDs          []string `json:"ids"`
		DownloadDone bool     `json:"DownloadDone"`
	}

	var req UpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		flytura.FormataRetornoHTTP(w, "Erro ao decodificar JSON", http.StatusBadRequest)
		return
	}

	if len(req.IDs) == 0 {
		flytura.FormataRetornoHTTP(w, "Lista de IDs está vazia", http.StatusBadRequest)
		return
	}

	// Converter os IDs para ObjectID
	var objectIDs []primitive.ObjectID
	for _, idStr := range req.IDs {
		objID, err := primitive.ObjectIDFromHex(idStr)
		if err != nil {
			flytura.FormataRetornoHTTP(w, fmt.Sprintf("ID inválido: %s", idStr), http.StatusBadRequest)
			return
		}
		objectIDs = append(objectIDs, objID)
	}

	// Conectar ao MongoDB
	// client, err := db.ConnectMongoDB(flytura.ConectionString)
	// if err != nil {
	// 	flytura.FormataRetornoHTTP(w, "Erro ao conectar ao banco de dados", http.StatusInternalServerError)
	// 	return
	// }
	// defer client.Disconnect(context.Background())

	// collection := client.Database(flytura.DBName).Collection(flytura.ImagesDBTableName)

	collection := db.MongoClient.Database(flytura.DBName).Collection(flytura.ImagesDBTableName)
	nowUTC := time.Now().UTC()

	dh, err := flytura.DiffHours(nowUTC, flytura.Fuso1, flytura.Fuso2)
	if err != nil {
		panic(err)
	}
	// Criar filtro e atualização
	filter := bson.M{"_id": bson.M{"$in": objectIDs}}
	update := bson.M{
		"$set": bson.M{
			"downloadDone": req.DownloadDone,
			"updatedAt":    nowUTC.Add(-time.Duration(dh) * time.Hour),
		},
	}

	result, err := collection.UpdateMany(context.Background(), filter, update)
	if err != nil {
		flytura.FormataRetornoHTTP(w, "Erro ao atualizar imagens", http.StatusInternalServerError)
		return
	}

	if result.ModifiedCount == 0 {
		flytura.FormataRetornoHTTP(w, "Nenhuma imagem foi atualizada", http.StatusOK)
		return
	}

	flytura.FormataRetornoHTTP(w, fmt.Sprintf("Atualizadas %d imagens com sucesso", result.ModifiedCount), http.StatusOK)
}

/*
Função criada por Ricardo Silva Ferreira
Inicio da criação 23/10/2025 17:50
Data Final da criação : 23/10/2025 17:54
*/

func UpdateStatusPdfOrXmlHandler(w http.ResponseWriter, r *http.Request) {
	// Validar o token de autenticação
	status, msg := flytura.TokenValido(w, r)
	if !status {
		http.Error(w, fmt.Sprintf("Erro ao validar token: %v", msg), http.StatusUnauthorized)
		return
	}

	type request struct {
		ID         primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
		FileType   string             `json:"fileType" bson:"fileType"`
		DownloadOk bool               `json:"downloadOk" bson:"downloadOk,omitempty"`
		UpdatedAt  time.Time          `json:"updatedAt" bson:"updatedAt,omitempty"`
	}
	// Decodificar o JSON recebido
	var data request
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

		nowUTC := time.Now().UTC()

		dh, err := flytura.DiffHours(nowUTC, flytura.Fuso1, flytura.Fuso2)
		if err != nil {
			panic(err)
		}
		data.UpdatedAt = nowUTC.Add(-time.Duration(dh) * time.Hour)
	}

	// fmt.Print("AAAAA", data)

	update := bson.M{}
	// Criar o objeto de atualização
	if data.FileType == "pdf" {
		update = bson.M{
			"$set": bson.M{
				"downloadPDFDone": data.DownloadOk,
				"updatedAt":       data.UpdatedAt,
			},
		}
	}

	if data.FileType == "xml" {
		update = bson.M{
			"$set": bson.M{
				"downloadXMLDone": data.DownloadOk,
				"updatedAt":       data.UpdatedAt,
			},
		}
	}

	// Conectar ao MongoDB e atualizar o usuário
	// client, err := db.ConnectMongoDB(flytura.ConectionString)
	// if err != nil {
	// 	flytura.FormataRetornoHTTP(w, "Erro ao conectar ao banco de dados", http.StatusInternalServerError)

	// 	// http.Error(w, "Erro ao conectar ao banco de dados", http.StatusInternalServerError)
	// 	return
	// }
	// defer client.Disconnect(context.Background())

	// collection := client.Database(flytura.DBName).Collection(flytura.ImagesDBTableName)

	collection := db.MongoClient.Database(flytura.DBName).Collection(flytura.ImagesDBTableName)
	result, err := collection.UpdateOne(context.Background(), bson.M{"_id": data.ID}, update)
	if err != nil {
		flytura.FormataRetornoHTTP(w, "Erro ao atualizar dados da imagem", http.StatusInternalServerError)

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
	flytura.FormataRetornoHTTP(w, "Imagem atualizda com sucesso", http.StatusOK)

}

/*
Função criada por Ricardo Silva Ferreira
Inicio da criação 25/11/2025 16:22
Data Final da criação : 25/11/2025 16:22
OBS : 11/06/2026 foi adicionado para excluir varios que são separados por ponto e virgula
*/

func DeleteImagesDBByIDHandler(w http.ResponseWriter, r *http.Request) {
	// Validar método HTTP
	if r.Method != http.MethodDelete {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	// Validar token
	status, msg := flytura.TokenValido(w, r)
	if !status {
		http.Error(w, fmt.Sprintf("erro ao validar token: %v", msg), http.StatusUnauthorized)
		return
	}

	id := r.URL.Query().Get("id")
	xmlFile := r.URL.Query().Get("xmlFile")
	pdfFile := r.URL.Query().Get("pdfFile")
	zipFile := r.URL.Query().Get("zipFile")

	fmt.Println("pdfFile params ", pdfFile)
	// fmt.Println("xmlFile ", xmlFile)
	// fmt.Println("zipFile ", zipFile)

	if id == "" {
		http.Error(w, "ID não informado na rota", http.StatusBadRequest)
		return
	}

	if pdfFile != "" {

		pdfFiles := strings.Split(pdfFile, ";")

		// fmt.Println("pdf files ", pdfFiles)

		for _, f := range pdfFiles {
			f = strings.TrimSpace(f) // ✅ limpa espaços

			if f == "" {
				continue
			}
			f = path.Base(f)

			// fmt.Println("PDF", f)

			err := DeleteFromS3(f)
			if err != nil {
				http.Error(w, fmt.Sprintf("erro ao excluir arquivo PDF: %v", err), http.StatusInternalServerError)
				return
			}
		}
	}

	if xmlFile != "" {

		xmlFiles := strings.Split(xmlFile, ";")

		for _, f := range xmlFiles {
			f = strings.TrimSpace(f)

			if f == "" {
				continue
			}

			f = path.Base(f)
			err := DeleteFromS3(f)
			if err != nil {
				http.Error(w, fmt.Sprintf("erro ao excluir arquivo XML: %v", err), http.StatusInternalServerError)
				return
			}
		}
	}

	errorZipFile := DeleteFromS3(zipFile)
	if errorZipFile != nil {
		http.Error(w, fmt.Sprintf("erro ao excluir arquivo ZIP: %v", errorZipFile), http.StatusInternalServerError)
		return
	}

	// Executar exclusão
	err := DeleteImagesDBByID(db.MongoClient, flytura.DBName, flytura.ImagesDBTableName, id)
	if err != nil {
		http.Error(w, fmt.Sprintf("erro ao excluir registro: %v", err), http.StatusInternalServerError)
		return
	}

	// Retornar resposta JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := map[string]string{
		"message": "registro excluído com sucesso",
		"id":      id,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("erro ao codificar resposta JSON: %v", err)
	}
}

// */
// func DeleteImagesDBByIDHandler(w http.ResponseWriter, r *http.Request) {
// 	// Validar método HTTP
// 	if r.Method != http.MethodDelete {
// 		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
// 		return
// 	}

// 	// Validar token
// 	status, msg := flytura.TokenValido(w, r)
// 	if !status {
// 		http.Error(w, fmt.Sprintf("erro ao validar token: %v", msg), http.StatusUnauthorized)
// 		return
// 	}

// 	id := r.URL.Query().Get("id")
// 	xmlFile := r.URL.Query().Get("xmlFile")
// 	pdfFile := r.URL.Query().Get("pdfFile")
// 	zipFile := r.URL.Query().Get("zipFile")

// 	// fmt.Println("pdfFile ", pdfFile)
// 	// fmt.Println("xmlFile ", xmlFile)
// 	// fmt.Println("zipFile ", zipFile)

// 	if id == "" {
// 		http.Error(w, "ID não informado na rota", http.StatusBadRequest)
// 		return
// 	}

// 	if pdfFile != "" {
// 		errorPdfFIle := DeleteFromS3(pdfFile)
// 		if errorPdfFIle != nil {
// 			http.Error(w, fmt.Sprintf("erro ao excluir arquivo PDF: %v", errorPdfFIle), http.StatusInternalServerError)
// 			return
// 		}
// 	}

// 	if xmlFile != "" {
// 		errorXmlFile := DeleteFromS3(xmlFile)
// 		if errorXmlFile != nil {
// 			http.Error(w, fmt.Sprintf("erro ao excluir arquivo PDF: %v", errorXmlFile), http.StatusInternalServerError)
// 			return
// 		}
// 	}

// 	errorZipFile := DeleteFromS3(zipFile)
// 	if errorZipFile != nil {
// 		http.Error(w, fmt.Sprintf("erro ao excluir arquivo ZIP: %v", errorZipFile), http.StatusInternalServerError)
// 		return
// 	}

// 	// Executar exclusão
// 	err := DeleteImagesDBByID(db.MongoClient, flytura.DBName, flytura.ImagesDBTableName, id)
// 	if err != nil {
// 		http.Error(w, fmt.Sprintf("erro ao excluir registro: %v", err), http.StatusInternalServerError)
// 		return
// 	}

// 	// Retornar resposta JSON
// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(http.StatusOK)

// 	response := map[string]string{
// 		"message": "registro excluído com sucesso",
// 		"id":      id,
// 	}

// 	if err := json.NewEncoder(w).Encode(response); err != nil {
// 		log.Printf("erro ao codificar resposta JSON: %v", err)
// 	}
// }

/*
	Função criada por Ricardo Silva Ferreira
	Inicio da criação 10/06/2026 21:28
	Data Final da criação : 10/06/2026 22:25
*/

func UploadManualImportInvoicesRecordHandler(w http.ResponseWriter, r *http.Request) {

	fmt.Println("files teste")
	err := r.ParseMultipartForm(10 << 20) // 10MB
	if err != nil {
		http.Error(w, "Erro ao processar formulário", http.StatusBadRequest)
		return
	}

	files := r.MultipartForm.File["file"]

	fmt.Println("files", files)

	if len(files) == 0 {
		http.Error(w, "Nenhum arquivo enviado", http.StatusBadRequest)
		return
	}

	nameFiles := ""

	for cont, fileHeader := range files {

		// fmt.Println("📄 Nome:", fileHeader.Filename)
		// fmt.Println("📄 Tipo:", fileHeader.Header.Get("Content-Type"))

		if cont == len(files)-1 {
			nameFiles += fileHeader.Filename

		} else {
			nameFiles += fileHeader.Filename + ";"
		}

		// ✅ validar extensão
		if !strings.HasSuffix(strings.ToLower(fileHeader.Filename), ".pdf") {
			http.Error(w, "Apenas arquivos PDF são permitidos", http.StatusBadRequest)
			return
		}

		// ✅ validar content-type
		contentType := fileHeader.Header.Get("Content-Type")
		if contentType != "application/pdf" {
			http.Error(w, "Arquivo não é PDF válido", http.StatusBadRequest)
			return
		}

		file, err := fileHeader.Open()
		if err != nil {
			http.Error(w, "Erro ao ler arquivo", http.StatusBadRequest)
			return
		}
		defer file.Close()

		// file, err := fileHeader.Open()
		// if err != nil {
		// 	http.Error(w, "Erro ao abrir arquivo", http.StatusInternalServerError)
		// 	return
		// }

		// ✅ exemplo: ler conteúdo (opcional)
		// data, err := io.ReadAll(file)
		// if err != nil {
		// 	http.Error(w, "Erro ao ler PDF", http.StatusInternalServerError)
		// 	file.Close()
		// 	return
		// }
		// fmt.Println("✅ PDF carregado:", fileHeader.Filename, "| Tamanho:", len(data))

		// Lê o conteúdo do ZIP para memória
		buf := new(bytes.Buffer)
		_, err = io.Copy(buf, file)

		if err != nil {
			http.Error(w, "Erro ao copiar conteúdo do ZIP", http.StatusInternalServerError)
			return
		}

		// //MNADA PARA O PARALELISMO E ENVIA PARA O S3
		flytura.UploadChan <- flytura.UploadTask{
			FileContent: buf.Bytes(),
			FileName:    fileHeader.Filename,
			CompanyCode: "",
			Key:         "",
		}

		// 👉 aqui você pode:
		// salvar no disco
		// enviar para S3
		// processar conteúdo

		file.Close() // ✅ correto (não usar defer no loop)
	}

	// ✅ campos do formdata
	userNameImport := r.FormValue("userName")
	userInserted := r.FormValue("idUserInserted")
	companyCode := r.FormValue("companyCode")
	key := r.FormValue("key")
	billedFlytura := r.FormValue("billedFlytura")

	// fmt.Println("👤 Arquivos:", nameFiles)
	// fmt.Println("👤 userName:", userNameImport)
	// fmt.Println("👤 userInserted:", userInserted)
	// fmt.Println("🏢 companyCode:", companyCode)
	// fmt.Println("🔑 key:", key)
	// fmt.Println("📊 billedFlytura:", billedFlytura)

	airLineData, errAirLineName := airLine.GetAirLineFileName(db.MongoClient, flytura.DBName, "airline", companyCode)
	if errAirLineName != nil {
		log.Println("Erro ao obter nome do arquivo:", errAirLineName)
	}

	nowUTC := time.Now().UTC()
	dh, err := flytura.DiffHours(nowUTC, flytura.Fuso1, flytura.Fuso2)
	if err != nil {
		panic(err)
	}

	var bf bool

	if billedFlytura == "true" {
		bf = true
	} else {
		bf = false
	}

	image := models.ImagesDB{
		ID:              primitive.NewObjectID(),
		FileName:        nameFiles,
		DtImport:        nowUTC.Add(-time.Duration(dh) * time.Hour),
		PDFFileName:     nameFiles,
		XMLFileName:     "",
		CompanyCode:     airLineData["Code"].(string),
		CompanyName:     airLineData["Name"].(string),
		DownloadDone:    false,
		DownloadPDFDone: false,
		DownloadXMLDone: false,
		Key:             key,
		ZipFileName:     "",
		IdUserInserted:  userInserted,
		UserNameImport:  userNameImport,
		OriginData:      "Upload Manual",
		ServerDate:      nowUTC,
		BilledFlytura:   bf,
	}

	InsertIMGS3(db.MongoClient, flytura.DBName, "imagesDB", image)

	// Retornar resposta JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := map[string]string{
		"message": "Upload feito com sucesso",
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("erro ao codificar resposta JSON: %v", err)
	}
}

/*
Função criada por Ricardo Silva Ferreira
Inicio da criação 11/06/2026 13:03
Data Final da criação : 11/06/2026 13:04
*/
type CheckDuplicatesRequest struct {
	Files []string `json:"files"`
}

func CheckDuplicatePDFsHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Método deve ser POST", http.StatusMethodNotAllowed)
		return
	}

	// ✅ Validar token
	status, msg := flytura.TokenValido(w, r)
	if !status {
		http.Error(w, fmt.Sprintf("erro ao validar token: %v", msg), http.StatusUnauthorized)
		return
	}

	// ✅ Ler body JSON
	var req CheckDuplicatesRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	if len(req.Files) == 0 {
		http.Error(w, "files é obrigatório", http.StatusBadRequest)
		return
	}

	// ✅ Chamar função
	duplicates, err := checkExistingPDFs(
		db.MongoClient,
		flytura.DBName,
		flytura.ImagesDBTableName,
		req.Files,
	)
	if err != nil {
		http.Error(w, "Erro ao verificar duplicados", http.StatusInternalServerError)
		return
	}

	// ✅ Resposta
	response := map[string]any{
		"duplicates": duplicates,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

/*
Função criada por Ricardo Silva Ferreira
Inicio da criação 12/06/2026 14:13
Data Final da criação : 12/06/2026 14:15
*/
func CountLast30DaysByDtImportsHandler(w http.ResponseWriter, r *http.Request) {

	// 🔐 Validação do token
	status, msg := flytura.TokenValido(w, r)
	if !status {
		http.Error(w, fmt.Sprintf("erro ao validar token: %v", msg), http.StatusUnauthorized)
		return
	}

	// 🔹 Chamada da nova função (sem filtros)
	total, err := CountLast30DaysByDtImport(
		db.MongoClient,
		flytura.DBName,
		flytura.ImagesDBTableName,
	)
	if err != nil {
		http.Error(w, fmt.Sprintf("erro ao contar registros: %v", err), http.StatusInternalServerError)
		return
	}

	// 🔹 Estrutura de resposta
	response := map[string]interface{}{
		"totalLast30Days": total,
	}

	// 🔹 Retorno JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("erro ao codificar JSON: %v", err)
	}
}

/*
Função criada por Ricardo Silva Ferreira
Inicio da criação 12/06/2026 14:40
Data Final da criação : 12/06/2026 14:45
*/
func CountLast30DaysByAmountRangeHandler(w http.ResponseWriter, r *http.Request) {

	// 🔐 Validação do token
	status, msg := flytura.TokenValido(w, r)
	if !status {
		http.Error(w, fmt.Sprintf("erro ao validar token: %v", msg), http.StatusUnauthorized)
		return
	}

	// 🔹 Chamada da função (agora retorna slice)
	data, err := CountLast30DaysByAmountRange(
		db.MongoClient,
		flytura.DBName,
		flytura.ConciliationTableName,
	)
	if err != nil {
		http.Error(w, fmt.Sprintf("erro ao buscar dados: %v", err), http.StatusInternalServerError)
		return
	}

	// 🔹 Retorno direto (já no formato correto)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("erro ao codificar JSON: %v", err)
	}
}
