package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

type Livro struct {
	Id     int    `json:"id"`
	Titulo string `json:"titulo"`
	Autor  string `json:"autor"`
}

type Usuario struct {
	Id    int    `json:"id"`
	Login string `json:"login"`
	Senha string `json:"senha"`
}

var Livros []Livro = []Livro{
	Livro{
		Id:     1,
		Titulo: "O ladrão de Raios",
		Autor:  "Rick Riordan",
	},

	Livro{
		Id:     2,
		Titulo: "Mar de monstros",
		Autor:  "Rick Riordan",
	},
	Livro{
		Id:     3,
		Titulo: "A maldição dos titãns",
		Autor:  "Rick Riordan",
	},
}

var Usuarios []Usuario = []Usuario{}

//Criando a função para chamar o servidor:
func criarServidor() {
	rotas()
	porta := ":8080"
	fmt.Println("O servidor está rodando com sucesso! na porta: ", porta)
	http.ListenAndServe(porta, nil)
}

func rotas() {
	http.HandleFunc("/", criarRotaPrincipal)
	http.HandleFunc("/sobre", criarRotaSobre)
	http.HandleFunc("/livros", verificarMetodos)
	http.HandleFunc("/livros/", buscarLivro)
	http.HandleFunc("/cadastro", criarLogin)
	http.HandleFunc("/login", login)

}

//criando a rota principal:
func criarRotaPrincipal(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Seja bem vindo a minha rota principal")
}

func verificarMetodos(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		criarRotasListarLivros(w, r)
	} else if r.Method == "POST" {
		partes := strings.Split(r.URL.Path, "/")
		if partes[1] == "livros" {
			criarRotasCadastrarLivro(w, r)
		}
		if partes[1] == "cadastro" {
			criarLogin(w, r)
		}
		if partes[1] == "login" {
			login(w, r)
		}

	}

}

func criarRotaSobre(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Meu nome é Gabriel Aderaldo.")
}

//Regra de negocio de Listar livros
func criarRotasListarLivros(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {

		return
	}

	w.Header().Set("Content-type", "Application/json")
	enconder := json.NewEncoder(w)
	enconder.Encode(Livros)
}

//Regra de negocio de Cadastrar livros
func criarRotasCadastrarLivro(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "Application/json")
	w.WriteHeader(http.StatusCreated) //Basicamente estou setando o estatus 201 na resposta do HTTP
	body, error := ioutil.ReadAll(r.Body)
	if error != nil {
		fmt.Fprint(w, error)
	}
	var novoLivro Livro
	json.Unmarshal(body, &novoLivro)
	novoLivro.Id = len(Livros) + 1
	Livros = append(Livros, novoLivro)
	encoder := json.NewEncoder(w)
	encoder.Encode(novoLivro)
}

//Regra de negocio de Buscar livros
func buscarLivro(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "Application/json")
	partes := strings.Split(r.URL.Path, "/")

	if len(partes) > 3 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	id, _ := strconv.Atoi(partes[2])

	for _, livro := range Livros {
		if livro.Id == id {
			json.NewEncoder(w).Encode(livro)
			fmt.Fprintln(w, "Nome do livro: ", livro.Titulo, "/", "Nome do Autor: ", livro.Autor)
		}
	}

	w.WriteHeader(http.StatusNotFound)
}

//rota para criar login e senha
func criarLogin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "Application/json")
	w.WriteHeader(http.StatusCreated)

	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		fmt.Fprintln(w, err)
	}

	var novoLogin Usuario
	json.Unmarshal(body, &novoLogin)
	novoLogin.Id = len(Usuarios) + 1
	Usuarios = append(Usuarios, novoLogin)
	enconder := json.NewEncoder(w)
	enconder.Encode(novoLogin)
}

//login e senha
func login(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintln(w, err)
	}

	var novoUsuario Usuario
	json.Unmarshal(body, &novoUsuario)

	for _, usuario := range Usuarios {
		if novoUsuario.Login == usuario.Login && novoUsuario.Senha == usuario.Senha {
			fmt.Fprintln(w, "LOGADO")
		}

		if novoUsuario.Login != usuario.Login || novoUsuario.Senha != usuario.Senha {
			w.WriteHeader(401)
			fmt.Fprintln(w, "Login ou senha incorreto")
		}
	}
}

func main() {
	criarServidor()
}
