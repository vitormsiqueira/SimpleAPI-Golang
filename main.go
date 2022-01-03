package main

import (
	"encoding/json" // utilizado para transformar a lista struct em json
	"fmt"
	"io/ioutil"
	"log" // Funções de log
	"net/http"
	"strconv"

	// "strings"

	"github.com/gorilla/mux"
)

// Define a estrutura de dados do Livro
type Livro struct {
	Id     int    `json:"id"` // define o nome do campo qdo convertido p/ json
	Titulo string `json:"titulo"`
	Autor  string `json:"autor"`
}

// Essa lista servirá como um "banco de dados" na memória
var Livros []Livro = []Livro{
	Livro{
		Id:     1,
		Titulo: "Dom Casmurro",
		Autor:  "Machado de Assis",
	},
	Livro{
		Id:     2,
		Titulo: "O Filho de Mil Homens",
		Autor:  "Valter Hugo Mãe",
	},
	Livro{
		Id:     3,
		Titulo: "A Arte da Guerra",
		Autor:  "Sun Tzu",
	},
}

// Função principal (Tela inicial)
func rotaPrincipal(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Bem vindo!")
}

// Essa função codifica a lista de livros no formato json
func listBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// primeiro fazemos o encoder
	encoder := json.NewEncoder(w)
	encoder.Encode(Livros)
}

// Essa função cadastra um licro baseado no corpo enviado pelo POST
func createBook(w http.ResponseWriter, r *http.Request) {

	// Setamos o tipo da resposta da requisição
	w.Header().Set("Content-Type", "application/json")

	// Define o status para 201 (Criado)
	w.WriteHeader(http.StatusCreated)

	// Aqui temos o que foi enviado pelo body do metodo POST (dados do livro)
	body, _ := ioutil.ReadAll(r.Body)

	// Definimos uma nova variavel e atribuimos o conteudo recebido
	var newLivro Livro
	json.Unmarshal(body, &newLivro)

	// Continua a sequencia do id
	newLivro.Id = len(Livros) + 1

	// Adicionamos esse conteudo na lista de livros
	Livros = append(Livros, newLivro)

	// Codificamos
	encoder := json.NewEncoder(w)
	encoder.Encode(newLivro)
}

func searchBook(w http.ResponseWriter, r *http.Request) {
	// Setamos o tipo da resposta da requisição
	w.Header().Set("Content-Type", "application/json")

	// Variável que recebe a resposta
	vars := mux.Vars(r)
	// Extraimos o id do livro a ser procurado
	key := vars["id"]

	// converte de string para int
	key_int, _ := strconv.Atoi(key)

	// Percorre todos os livros
	// se o livro.Id for igual a chave key
	// retorna o livro codificado em JSON
	for _, livro := range Livros {
		if livro.Id == key_int {
			json.NewEncoder(w).Encode(livro)
			return
		}
	}

	w.WriteHeader(http.StatusNotFound)
}

func deleteBook(w http.ResponseWriter, r *http.Request) {
	// Variável que recebe a resposta
	vars := mux.Vars(r)
	// Extraimos o id do livro a ser deletado
	key := vars["id"]

	// converte de string para int
	key_int, _ := strconv.Atoi(key)

	// Percorremos todos os livros
	for index, livro := range Livros {
		// Se o id recebido consta na nossa lista de livros...
		if livro.Id == key_int {
			// Realizamos a exclusão
			Livros = append(Livros[:index], Livros[index+1:]...)
		}
	}
}

func updateBook(w http.ResponseWriter, r *http.Request) {

	// Setamos o tipo da resposta da requisição
	w.Header().Set("Content-Type", "application/json")

	// Aqui emos o que foi enviado pelo body do metodo POST (dados do livro)
	body, _ := ioutil.ReadAll(r.Body)

	// Definimos uma nova variavel e atribuimos o conteudo recebido
	var livroModificado Livro
	erroJson := json.Unmarshal(body, &livroModificado)

	if erroJson != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	vars := mux.Vars(r)
	key := vars["id"]

	// converte de string para int
	key_int, _ := strconv.Atoi(key)

	indiceLivro := -1
	for indice, livro := range Livros {
		if livro.Id == key_int {
			indiceLivro = indice
			break
		}
	}

	if indiceLivro < 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Adicionamos esse conteudo na lista de livros
	Livros[indiceLivro] = livroModificado

	// Codificamos
	encoder := json.NewEncoder(w)
	encoder.Encode(livroModificado)
}

// Configura o servidor definindo as rotas e a porta de acesso
func configurarServidor() {
	var port string = "1337"

	// configura as rotas
	myRouter := mux.NewRouter()

	myRouter.HandleFunc("/", rotaPrincipal).Methods("GET")
	myRouter.HandleFunc("/livros", listBook).Methods("GET")
	myRouter.HandleFunc("/livros", createBook).Methods("POST")
	myRouter.HandleFunc("/livros/{id}", searchBook).Methods("GET")
	myRouter.HandleFunc("/livros/{id}", deleteBook).Methods("DELETE")
	myRouter.HandleFunc("/livros/{id}", updateBook).Methods("PUT")

	fmt.Println("O servidor está rodando na porta " + port)

	// http.ListenAndServe(":"+port, nil) // DefaultServerMux
	log.Fatal(http.ListenAndServe(":"+port, myRouter)) // Se houver erro, será informado o log
}

func main() {
	configurarServidor()
}
