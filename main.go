package main

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

func postPessoas(c *gin.Context) {

}

func getPessoas(c *gin.Context) {

}

func indexPessoas(c *gin.Context) {

}

func contagemPessoas(c *gin.Context) {

}
