package main

import (
	"context"
	"fmt"
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
	connString := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable&pool_min_conns=%s&pool_max_conns=%[6]s",
		getEnv("DATABASE_USER"),
		getEnv("DATABASE_PASS"),
		getEnv("DATABASE_HOST"),
		getEnv("DATABASE_PORT"),
		getEnv("DATABASE_NAME"),
		getEnv("DATABASE_POOL"),
	)

	config, err := pgxpool.ParseConfig(connString)

	if err != nil {
		panic(err)
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), config)

	if err != nil {
		panic(err)
	}

	defer pool.Close()

	if err := pool.Ping(context.Background()); err != nil {
		panic(err)
	}

	router := gin.New()
	router.Use(gin.Recovery())

	router.POST("/pessoas", postPessoas(pool))
	router.GET("/pessoas/:id", getPessoas(pool))
	router.GET("/pessoas", indexPessoas(pool))
	router.GET("/contagem-pessoas", contagemPessoas(pool))

	router.Run("localhost:" + getEnv("PORT"))
}

func getEnv(envName string) string {
	env, ok := os.LookupEnv(envName)

	if !ok {
		message := fmt.Sprintf("env not declared $%s", envName)
		panic(message)
	}

	return env
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
			buildBusca(&person),
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
	if person.Apelido == "" || person.Nome == "" || person.Nascimento == "" {
		return true
	}

	return false
}

func invalidFields(person *Pessoa) bool {
	if len(person.Apelido) > 32 || len(person.Nome) > 100 {
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

func buildBusca(person *Pessoa) string {
	size := len(person.Apelido) + len(person.Nome) + 2 + (32 + ((len(person.Stack) - 1) * 33))

	var busca strings.Builder
	busca.Grow(size)

	busca.WriteString(strings.Map(unicode.ToLower, person.Apelido))
	busca.WriteString(" ")
	busca.WriteString(strings.Map(unicode.ToLower, person.Nome))
	busca.WriteString(" ")

	if person.Stack != nil {
		busca.WriteString(person.Stack[0])

		for _, v := range person.Stack[1:] {
			busca.WriteString(" ")
			busca.WriteString(strings.Map(unicode.ToLower, v))
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

		term = strings.Map(unicode.ToLower, term)

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
