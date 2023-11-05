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
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Pessoa struct {
	ID         string   `json:"id"`
	Apelido    string   `json:"apelido"`
	Nome       string   `json:"nome"`
	Nascimento string   `json:"nascimento"`
	Stack      []string `json:"stack"`
}

func main() {
	connString := "postgres://postgres:postgres@127.0.0.1:5432/rinha?sslmode=disable&pool_min_conns=1&pool_max_conns=15"

	config, err := pgxpool.ParseConfig(connString)

	if err != nil {
		panic(err)
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), config)

	defer pool.Close()

	if err != nil {
		panic(err)
	}

	// router := gin.Default()
	router := gin.New()
	router.Use(gin.Recovery())

	router.POST("/pessoas", postPessoas(pool))
	router.GET("/pessoas/:id", getPessoas(pool))
	router.GET("/pessoas", indexPessoas(pool))
	router.GET("/contagem-pessoas", contagemPessoas(pool))

	port, ok := os.LookupEnv("PORT")

	if !ok {
		port = "9999"
	}

	router.Run("localhost:" + port)
}

func postPessoas(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var person Pessoa

		if err := ctx.BindJSON(&person); err != nil {
			ctx.Status(http.StatusBadRequest)
			return
		}

		if missingFields(&person) {
			ctx.Status(http.StatusUnprocessableEntity)
			return
		}

		if invalidFields(&person) {
			ctx.Status(http.StatusBadRequest)
			return
		}

		row := pool.QueryRow(
			ctx,
			"INSERT INTO pessoas (apelido, nome, nascimento, stack, busca) VALUES ($1, $2, $3, $4, $5) RETURNING id",
			person.Apelido,
			person.Nome,
			person.Nascimento,
			person.Stack,
			buildBusca(person.Apelido, person.Nome, &person.Stack),
		)

		if err := row.Scan(&person.ID); err != nil {
			ctx.Status(http.StatusUnprocessableEntity)
			return
		}

		ctx.Header("Location", "/pessoas/"+person.ID)
		ctx.Status(http.StatusCreated)
	}
}

func missingFields(person *Pessoa) bool {
	if len(person.Apelido) == 0 || len(person.Nome) == 0 || len(person.Nascimento) == 0 {
		return true
	}

	return false
}

func invalidFields(person *Pessoa) bool {
	if len(person.Apelido) > 32 && len(person.Nome) > 100 {
		return true
	}

	if _, err := time.Parse("2006-01-02", person.Nascimento); err != nil {
		return true
	}

	for _, v := range person.Stack {
		if v == "" || len(v) > 32 {
			return true
		}
	}

	return false
}

func buildBusca(apelido string, nome string, stack *[]string) string {
	size := len(apelido) + len(nome) + 2 + (len(*stack) * 33)

	var busca strings.Builder
	busca.Grow(size)

	busca.WriteString(strings.Map(unicode.ToLower, apelido))
	busca.WriteString(" ")
	busca.WriteString(strings.Map(unicode.ToLower, nome))
	busca.WriteString(" ")

	if stack != nil {
		for _, v := range *stack {
			busca.WriteString(strings.Map(unicode.ToLower, v))
			busca.WriteString(" ")
		}
	}

	return busca.String()
}

func getPessoas(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.Param("id")

		var person Pessoa

		if err := pool.QueryRow(ctx, "SELECT id, apelido, nome, nascimento, stack FROM pessoas WHERE id = $1", id).Scan(&person.ID, &person.Apelido, &person.Nome, &person.Nascimento, &person.Stack); err != nil {
			ctx.Status(http.StatusNotFound)
			return
		}

		ctx.JSON(http.StatusOK, person)

	}
}

func indexPessoas(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		term, passed := ctx.GetQuery("t")

		if !passed || term == "" {
			ctx.Status(http.StatusBadRequest)
			return
		}

		term = strings.ToLower(term)

		rows, err := pool.Query(ctx, "SELECT id, apelido, nome, nascimento, stack FROM pessoas WHERE busca LIKE '%' || $1 || '%' LIMIT 50", term)

		if err != nil {
			ctx.Status(http.StatusUnprocessableEntity)
			return
		}

		people, err := pgx.CollectRows(rows, pgx.RowToStructByPos[Pessoa])

		ctx.JSON(http.StatusOK, people)
	}
}

func contagemPessoas(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var count int
		_ = pool.QueryRow(ctx, "SELECT COUNT(*) FROM pessoas").Scan(&count)

		ctx.String(http.StatusOK, strconv.Itoa(count))
	}
}
