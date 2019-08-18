package main

import (
	"database/sql" // Pacote Database SQL para realizar Query
	"log"
	"net/http"      // Gerencia URLs e Servidor Web
	"text/template" // Gerencia templates

	_ "github.com/denisenkom/go-mssqldb"
)

type Usuarios struct {
	UsuarioID    int
	UsuarioNome  string
	UsuarioLogin string
	UsuarioEmail string
	UsuarioSenha string
}

// Função dbConn, abre a conexão com o banco de dados
func dbConn() (db *sql.DB) {
	dbDriver := "mssql"
	dbServer := "DESKTOP-G7K5VG3\\SQLEXPRESS"
	dbUser := "AdminGo"
	dbPass := "040195go"
	dbName := "GO"

	db, err := sql.Open(dbDriver, dbServer, dbUser+":"+dbPass+"@/"+dbName)
	if err != nil {
		panic(err.Error())
	}
	return db
}

func Index(w http.ResponseWriter, r *http.Request) {
	// Abre a conexão com o banco de dados utilizando a função dbConn()
	db := dbConn()
	// Realiza a consulta com banco de dados e trata erros
	selDB, err := db.Query("SELECT * FROM Usuario ORDER BY UsuarioID DESC")
	if err != nil {
		panic(err.Error())
	}

	// Monta a struct para ser utilizada no template
	n := Usuarios{}

	// Monta um array para guardar os valores da struct
	res := []Usuarios{}

	// Realiza a estrutura de repetição pegando todos os valores do banco
	for selDB.Next() {
		// Armazena os valores em variáveis
		var UsuarioID int
		var UsuarioNome, UsuarioLogin, UsuarioEmail, UsuarioSenha string

		// Faz o Scan do SELECT
		err = selDB.Scan(&UsuarioID, &UsuarioNome, &UsuarioLogin, &UsuarioEmail, &UsuarioSenha)
		if err != nil {
			panic(err.Error())
		}

		// Envia os resultados para a struct
		n.UsuarioID = UsuarioID
		n.UsuarioNome = UsuarioNome
		n.UsuarioEmail = UsuarioEmail
		n.UsuarioLogin = UsuarioLogin
		n.UsuarioSenha = UsuarioSenha

		// Junta a Struct com Array
		res = append(res, n)
	}

	// Abre a página Index e exibe todos os registrados na tela
	var tmpl = template.Must(template.ParseGlob("tmpl/*"))
	tmpl.ExecuteTemplate(w, "Index", res)

	// Fecha a conexão
	defer db.Close()
}

func Show(w http.ResponseWriter, r *http.Request) {
	db := dbConn()

	// Pega o ID do parametro da URL
	nId := r.URL.Query().Get("id")

	// Usa o ID para fazer a consulta e tratar erros
	selDB, err := db.Query("SELECT * FROM Usuario WHERE UsuarioID=?", nId)
	if err != nil {
		panic(err.Error())
	}

	// Monta a strcut para ser utilizada no template
	n := Usuarios{}

	// Realiza a estrutura de repetição pegando todos os valores do banco
	for selDB.Next() {
		// Armazena os valores em variaveis
		var UsuarioID int
		var UsuarioNome, UsuarioLogin, UsuarioEmail, UsuarioSenha string

		// Faz o Scan do SELECT
		err = selDB.Scan(&UsuarioID, &UsuarioNome, &UsuarioLogin, &UsuarioEmail, &UsuarioSenha)
		if err != nil {
			panic(err.Error())
		}

		// Envia os resultados para a struct
		n.UsuarioID = UsuarioID
		n.UsuarioNome = UsuarioNome
		n.UsuarioEmail = UsuarioEmail
		n.UsuarioLogin = UsuarioLogin
		n.UsuarioSenha = UsuarioSenha
	}

	// Mostra o template
	var tmpl = template.Must(template.ParseGlob("tmpl/*"))
	tmpl.ExecuteTemplate(w, "Show", n)

	// Fecha a conexão
	defer db.Close()
}

func New(w http.ResponseWriter, r *http.Request) {
	var tmpl = template.Must(template.ParseGlob("tmpl/*"))
	tmpl.ExecuteTemplate(w, "New", nil)
}

func Edit(w http.ResponseWriter, r *http.Request) {
	// Abre a conexão com banco de dados
	db := dbConn()

	// Pega o ID do parametro da URL
	nId := r.URL.Query().Get("id")

	selDB, err := db.Query("SELECT * FROM Usuario WHERE UsuarioID=?", nId)
	if err != nil {
		panic(err.Error())
	}

	// Monta a struct para ser utilizada no template
	n := Usuarios{}

	// Realiza a estrutura de repetição pegando todos os valores do banco
	for selDB.Next() {
		//Armazena os valores em variaveis
		var UsuarioID int
		var UsuarioNome, UsuarioLogin, UsuarioEmail, UsuarioSenha string

		// Faz o Scan do SELECT
		err = selDB.Scan(&UsuarioID, &UsuarioNome, &UsuarioLogin, &UsuarioEmail, &UsuarioSenha)
		if err != nil {
			panic(err.Error())
		}

		// Envia os resultados para a struct
		n.UsuarioID = UsuarioID
		n.UsuarioNome = UsuarioNome
		n.UsuarioLogin = UsuarioLogin
		n.UsuarioEmail = UsuarioEmail
		n.UsuarioSenha = UsuarioSenha
	}

	// Mostra o template com formulário preenchido para edição
	var tmpl = template.Must(template.ParseGlob("tmpl/*"))
	tmpl.ExecuteTemplate(w, "Edit", n)

	// Fecha a conexão com o banco de dados
	defer db.Close()
}

func Insert(w http.ResponseWriter, r *http.Request) {

	//Abre a conexão com banco de dados usando a função: dbConn()
	db := dbConn()

	// Verifica o METHOD do fomrulário passado
	if r.Method == "POST" {

		// Pega os campos do formulário
		UsuarioNome := r.FormValue("UsuarioNome")
		UsuarioEmail := r.FormValue("UsuarioEmail")
		UsuarioSenha := r.FormValue("UsuarioSenha")
		UsuarioLogin := r.FormValue("UsuarioLogin")

		// Prepara a SQL e verifica errors
		insForm, err := db.Prepare("INSERT INTO Usuario(UsuarioNome, UsuarioEmail, UsuarioSenha, UsuarioLogin) VALUES(?,?,?,?,?,?)")
		if err != nil {
			panic(err.Error())
		}

		// Insere valores do formulario com a SQL tratada e verifica errors
		insForm.Exec(UsuarioNome, UsuarioEmail, UsuarioSenha, UsuarioLogin)

		// Exibe um log com os valores digitados no formulário
		log.Println("INSERT: Nome: " + UsuarioNome + " | E-mail: " + UsuarioEmail)
	}

	// Encerra a conexão do dbConn()
	defer db.Close()

	//Retorna a HOME
	http.Redirect(w, r, "/", 301)
}

func Delete(w http.ResponseWriter, r *http.Request) {

	// Abre conexão com banco de dados usando a função: dbConn()
	db := dbConn()

	nId := r.URL.Query().Get("id")

	// Prepara a SQL e verifica errors
	delForm, err := db.Prepare("DELETE FROM Usuario WHERE UsuarioID=?")
	if err != nil {
		panic(err.Error())
	}

	// Insere valores do form com a SQL tratada e verifica errors
	delForm.Exec(nId)

	// Exibe um log com os valores digitados no form
	log.Println("DELETE")

	// Encerra a conexão do dbConn()
	defer db.Close()

	// Retorna a HOME
	http.Redirect(w, r, "/", 301)
}

func Update(w http.ResponseWriter, r *http.Request) {

	// Abre a conexão com o banco de dados usando a função: dbConn()
	db := dbConn()

	// Verifica o METHOD do formulário passado
	if r.Method == "POST" {

		// Pega os campos do formulário
		UsuarioNome := r.FormValue("UsuarioNome")
		UsuarioEmail := r.FormValue("UsuarioEmail")
		UsuarioSenha := r.FormValue("UsuarioSenha")
		UsuarioLogin := r.FormValue("UsuarioLogin")
		UsuarioID := r.FormValue("UsuarioID")

		// Prepara a SQL e verifica errors
		insForm, err := db.Prepare("UPDATE Usuario SET UsuarioNome=?, UsuarioEmail=?, UsuarioSenha=?, UsuarioLogin=? WHERE UsuarioID=?")
		if err != nil {
			panic(err.Error())
		}

		// Insere valores do formulário com a SQL tratada e verifica erros
		insForm.Exec(UsuarioNome, UsuarioEmail, UsuarioSenha, UsuarioLogin, UsuarioID)

		// Exibe um log com os valores digitados no formulario
		log.Println("UPDATE: Name: " + UsuarioNome + " |E-mail: " + UsuarioEmail)
	}

	// Encerra a conexão do dbConn()
	defer db.Close()

	// Retorna a HOME
	http.Redirect(w, r, "/", 301)
}

func main() {

	// Exibe mensagem que o servidor foi iniciado
	log.Println("Server started on: http://localhost:9000")

	// Gerencia as URLs
	http.HandleFunc("/", Index)
	http.HandleFunc("/show", Show)
	http.HandleFunc("/new", New)
	http.HandleFunc("/edit", Edit)

	// Ações
	http.HandleFunc("/insert", Insert)
	http.HandleFunc("/update", Update)
	http.HandleFunc("/delete", Delete)

	// Inicia o servidor na porta 9000
	http.ListenAndServe(":9000", nil)
}
