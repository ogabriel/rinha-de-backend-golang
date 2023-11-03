package main

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Pessoa struct {
	ID         uuid.UUID `json:"id"`
	Apelido    string    `json:"apelido"`
	Nome       string    `json:"nome"`
	Nascimento string    `json:"nascimento"`
	Stack      []string  `json:"stack"`
}

func main() {
	router := gin.Default()

	router.POST("/pessoas", postPessoas)
	router.GET("/pessoas/:id", getPessoas)
	router.GET("/pessoas", indexPessoas)
	router.GET("/contagem-pessoas", contagemPessoas)

	router.Run("localhost:9999")
}

func pool() {
	connString := "postgres://postgres:postgres@127.0.0.1:5432/rinha&pool_max_conns=100"
	dbpool, err := pgxpool.New(context.Background(), connString)

	defer dbpool.Close()

	log.SetPrefix("database: ")

	if err != nil {
		log.Panic("could not connect to database %v", err)
	}
}

func postPessoas(c *gin.Context) {

}

func getPessoas(c *gin.Context) {

}

func indexPessoas(c *gin.Context) {

}

func contagemPessoas(c *gin.Context) {

}
