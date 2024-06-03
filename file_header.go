package filehandler

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// Manipulador de arquivos para lidar com o envio de arquivos
func HandleFiles(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Método não suportado", http.StatusMethodNotAllowed)
		return
	}

	// Extrair os arquivos do formulário
	files, err := r.MultipartForm.Parse()
	if err != nil {
		log.Println("Erro ao analisar o formulário:", err)
		http.Error(w, "Erro interno do servidor", http.StatusInternalServerError)
		return
	}

	// Processar cada arquivo enviado
	for _, file := range files["file"] {
		// Validar o arquivo
		if err := validateFile(file); err != nil {
			log.Println("Erro ao validar o arquivo:", err)
			http.Error(w, "Arquivo inválido", http.StatusBadRequest)
			return
		}

		// Salvar o arquivo em um local temporário
		tmpFile, err := createTempFile(file)
		if err != nil {
			log.Println("Erro ao criar arquivo temporário:", err)
			http.Error(w, "Erro interno do servidor", http.StatusInternalServerError)
			return
		}
		defer os.Remove(tmpFile.Name())

		// Processar o arquivo (por exemplo, ler o conteúdo, converter formato)
		processFile(tmpFile)

		// Mover o arquivo para o local de destino (opcional)
		// moveFile(tmpFile)
	}

	// Responder com sucesso
	fmt.Fprintf(w, "Arquivos enviados com sucesso!")
}

// Função para validar o arquivo
func validateFile(fileHeader *multipart.FileHeader) error {
	// Validar o tamanho do arquivo
	if fileHeader.Size > 1024*1024*5 { // 5 MB
		return fmt.Errorf("Tamanho do arquivo excede o limite (5 MB)")
	}

	// Validar o tipo de arquivo
	allowedTypes := []string{"image/jpeg", "image/png", "application/pdf"}
	for _, allowedType := range allowedTypes {
		if fileHeader.Header.Get("Content-Type") == allowedType {
			return nil
		}
	}
	return fmt.Errorf("Tipo de arquivo não permitido")
}

// Função para criar um arquivo temporário
func createTempFile(fileHeader *multipart.FileHeader) (*os.File, error) {
	// Gerar nome de arquivo temporário
	tmpFileName := fmt.Sprintf("upload-%s-%d", strings.TrimSpace(fileHeader.Filename()), time.Now().UnixNano())
	tmpFilePath := filepath.Join(os.TempDir(), tmpFileName)

	// Criar o arquivo temporário
	tmpFile, err := os.Create(tmpFilePath)
	if err != nil {
		return nil, err
	}
	defer tmpFile.Close()

	// Copiar o conteúdo do arquivo enviado para o arquivo temporário
	_, err = io.Copy(tmpFile, fileHeader.Open())
	if err != nil {
		return nil, err
	}

	return tmpFile, nil
}

// Função para processar o arquivo (opcional)
func processFile(file *os.File) {
	// Leia o conteúdo do arquivo
	fileContent, err := ioutil.ReadAll(file)
	if err != nil {
		log.Println("Erro ao ler o arquivo:", err)
		return
	}

	// Processe o conteúdo do arquivo (por exemplo, converta o formato)
	// Imprima o conteúdo do arquivo (opcional)
	fmt.Println("Conteúdo do arquivo:", string(fileContent))
}

// Função para mover o arquivo para o local de destino (opcional)
func moveFile(tmpFile *os.File) error {
	// Defina o local de destino<p>
	destPath := "/caminho/para/o/local/de/destino"

	// Renomeie o arquivo temporário para o nome original
	err := os.Rename(tmpFile.Name(), filepath.Join(destPath, tmpFile.Name()))
