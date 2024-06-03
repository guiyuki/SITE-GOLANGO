package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

/*
func handleForm(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Método não suportado", http.StatusMethodNotAllowed)
		return
	}

	// Obter os valores dos campos de entrada
	label1 := r.FormValue("label1")
	label2 := r.FormValue("label2")

	// Processar os valores dos labels
	fmt.Println("Label 1:", label1)
	fmt.Println("Label 2:", label2)

	// Enviar uma resposta para o cliente
	fmt.Fprintf(w, "Texto dos labels recebido com sucesso!")
}
*/

// Esse serve para página estática
func primeiraResposta(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
}

type Pagina struct {
	Titulo string
	Corpo  string
}

func Erro(w http.ResponseWriter) {
	fmt.Printf("Deu ruim")
	fmt.Fprintf(w, `<h1>A página não foi encontrada.</hl>`)
	return
}

// Modelo para página Template
func Modelo(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("src/modelo.html")

	if err != nil {
		Erro(w)
	}
	/*
		fmt.Printf("Deu ruim")
			fmt.Fprintf(w, `<h1>A página não foi encontrada.</hl>`)
			return
		}
	*/
	pagina := Pagina{Titulo: "Titulo Poggers", Corpo: "Dados importantes"}
	tmpl.Execute(w, pagina)
}

func Boleto(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "src/boleto.html")
}

/*
func Pasta(w http.ResponseWriter, r *http.Request) {
	w.ResponseWriter(w, http.Dir("./src/"))
}
*/

func CSS(w http.ResponseWriter, r *http.Request) {
	// Leia o conteúdo do arquivo CSS
	cssData, _ := ioutil.ReadFile("src/estilo.css")
	/*
		if err != nil {
			http.Error(w, "Erro ao ler o arquivo CSS", http.StatusInternalServerError)
			return
		}*/

	// Defina o tipo de conteúdo como CSS
	w.Header().Set("Content-Type", "text/css")

	// Escreva o conteúdo CSS na resposta
	w.Write(cssData)
}

func main() {
	// Definir a porta do servidor
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Porta padrão
	}

	http.HandleFunc("/estilo.css", CSS)

	http.HandleFunc("/boleto", Boleto)
	//http.HandleFunc("/", Pasta)
	http.HandleFunc("/", Modelo)

	// Iniciar o servidor HTTP
	fmt.Println("Servidor iniciado na porta", port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal("Erro ao iniciar o servidor:", err)
	}
}
