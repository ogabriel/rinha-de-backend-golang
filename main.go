package main

import (
	"context"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Pessoa struct {
	ID         uuid.UUID `json:"id"`
	Apelido    string    `json:"apelido" binding:"required,max=32"`
	Nome       string    `json:"nome" binding:"required,max=100"`
	Nascimento string    `json:"nascimento" binding:"required,len=10"`
	Stack      []string  `json:"stack"`
}

var pool *pgxpool.Pool

func main() {
	connString := "postgres://postgres:postgres@127.0.0.1:5432/rinha?sslmode=disable&pool_min_conns=1&pool_max_conns=15"

	config, err := pgxpool.ParseConfig(connString)

	if err != nil {
		panic(err)
	}

	pool, err = pgxpool.NewWithConfig(context.Background(), config)

	defer pool.Close()

	if err != nil {
		panic(err)
	}

	// router := gin.Default()
	router := gin.New()
	router.Use(gin.Recovery())

	router.POST("/pessoas", postPessoas)
	router.GET("/pessoas/:id", getPessoas)
	router.GET("/pessoas", indexPessoas)
	router.GET("/contagem-pessoas", contagemPessoas)

	port, ok := os.LookupEnv("PORT")

	if !ok {
		port = "9999"
	}

	router.Run("localhost:" + port)
}

func postPessoas(c *gin.Context) {
	var person Pessoa

	if err := c.BindJSON(&person); err != nil {
		c.String(http.StatusBadRequest, "")
		return
	}

	if invalidStack(person.Stack) {
		c.String(http.StatusBadRequest, "")
		return
	}

	if _, err := time.Parse("2006-01-02", person.Nascimento); err != nil {
		c.String(http.StatusBadRequest, "")
		return
	}

	uuid, err := uuid.NewRandom()

	if err != nil {
		c.String(http.StatusUnprocessableEntity, "")
		return
	}

	_, err = pool.Exec(
		context.Background(),
		"INSERT INTO pessoas (id, apelido, nome, nascimento, stack, busca) VALUES ($1, $2, $3, $4, $5, $6)",
		uuid,
		person.Apelido,
		person.Nome,
		person.Nascimento,
		person.Stack,
		buildBusca(person),
	)

	if err != nil {
		c.String(http.StatusUnprocessableEntity, "")
		return
	}

	c.Header("Location", "/pessoas/"+uuid.String())
	c.JSON(http.StatusCreated, person)
}

func invalidStack(stack []string) bool {
	for _, v := range stack {
		if v == "" || len(v) > 32 {
			return true
		}
	}

	return false
}

func buildBusca(person Pessoa) string {
	var busca strings.Builder

	busca.WriteString(strings.Map(unicode.ToLower, person.Apelido))
	busca.WriteString(" ")
	busca.WriteString(strings.Map(unicode.ToLower, person.Nome))
	busca.WriteString(" ")

	if person.Stack != nil {
		for _, v := range person.Stack {
			busca.WriteString(strings.Map(unicode.ToLower, v))
			busca.WriteString(" ")
		}
	}

	return busca.String()
}

func getPessoas(c *gin.Context) {
	id := c.Param("id")

	var person Pessoa

	if err := pool.QueryRow(context.Background(), "SELECT id, apelido, nome, nascimento, stack FROM pessoas WHERE id = $1", id).Scan(&person.ID, &person.Apelido, &person.Nome, &person.Nascimento, &person.Stack); err != nil {
		c.String(http.StatusNotFound, err.Error())
		return
	}

	c.JSON(http.StatusOK, person)
}

func indexPessoas(c *gin.Context) {
	term, passed := c.GetQuery("t")

	if !passed || term == "" {
		c.String(http.StatusBadRequest, "")
		return
	}

	term = strings.ToLower(term)

	rows, err := pool.Query(context.Background(), "SELECT id, apelido, nome, nascimento, stack FROM pessoas WHERE busca LIKE '%' || $1 || '%' LIMIT 50", term)

	defer rows.Close()

	if err != nil {
		c.String(http.StatusUnprocessableEntity, err.Error())
		return
	}

	var people []Pessoa

	for rows.Next() {
		var person Pessoa

		if err := rows.Scan(&person.ID, &person.Apelido, &person.Nome, &person.Nascimento, &person.Stack); err != nil {
			c.String(http.StatusUnprocessableEntity, err.Error())
			return
		}

		people = append(people, person)
	}

	c.JSON(http.StatusOK, people)
}

func contagemPessoas(c *gin.Context) {
	var count int
	_ = pool.QueryRow(context.Background(), "SELECT COUNT(*) FROM pessoas").Scan(&count)

	c.String(http.StatusOK, strconv.Itoa(count))
}
