package awsS3

import (
	"fmt"
	"net/http"
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

	err = UploadToS3(file, "uploads/"+header.Filename)
	if err != nil {

		fmt.Println("Erro detalhado ao enviar para S3:", err)
		// return fmt.Errorf("erro ao enviar para S3: %w", err)

		return
	}

	fmt.Fprintf(w, "Arquivo %s enviado com sucesso!", header.Filename)
}
