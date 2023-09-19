package main

import (
	"fmt"
	"net/http"
	_ "tatovering/src/controllers"
	_ "tatovering/src/models"

	"github.com/gin-gonic/gin"
	supabase "github.com/nedpals/supabase-go"
)

func main() {
	// Initialize your Supabase client
	supabaseURL := "https://pfzlboeaonsookzcnniv.supabase.co"
	supabaseKey := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzdXBhYmFzZSIsInJlZiI6InBmemxib2Vhb25zb29remNubml2Iiwicm9sZSI6ImFub24iLCJpYXQiOjE2OTIwNjI3NTYsImV4cCI6MjAwNzYzODc1Nn0.KuEEX9EBIQmLTA02iPtqqNIewDmXITDxnIfD4qEqTN8"

	// Create a Gin router
	router := gin.Default()
	client := supabase.CreateClient(supabaseURL, supabaseKey)

	/*********************************************************
	* 				   	  CRUD TATUAGENS 				   	 *
	**********************************************************/
	router.POST("/tatuagens", func(c *gin.Context) {

		var requestBody struct {
			TatuadorId    int     `json:"tatuador_id"`
			AgendamentoId int     `json:"agendamento_id"`
			Preco         float32 `json:"preco"`
			Desenho       string  `json:"desenho"`
			Tamaho        int     `json:"tamanho"`
			Cor           string  `json:"cor"`
			Estilo        string  `json:"estilo"`
		}

		if err := c.ShouldBindJSON(&requestBody); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		row := Tatuagem{
			TatuadorId:    requestBody.TatuadorId,
			AgendamentoId: requestBody.AgendamentoId,
			Preco:         requestBody.Preco,
			Desenho:       requestBody.Desenho,
			Tamaho:        requestBody.Tamaho,
			Cor:           requestBody.Cor,
			Estilo:        requestBody.Estilo,
		}

		var results []Tatuagem

		// inserting data and receive error if exist
		err := client.DB.From("tatuagens").Insert(row).Execute(&results)

		// chack error returned
		if err != nil {
			// ginh.H used to returnd a json file
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// return a json file
		c.JSON(http.StatusOK, results)

	})

	router.PATCH("/tatuagens/:id", func(c *gin.Context) {
		tatuagemId := c.Param("id")

		var requestBody struct {
			TatuadorId    int     `json:"tatuador_id"`
			AgendamentoId int     `json:"agendamento_id"`
			Preco         float32 `json:"preco"`
			Desenho       string  `json:"desenho"`
			Tamaho        int     `json:"tamanho"`
			Cor           string  `json:"cor"`
			Estilo        string  `json:"estilo"`
		}

		if err := c.ShouldBindJSON(&requestBody); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		row := Tatuagem{
			TatuadorId:    requestBody.TatuadorId,
			AgendamentoId: requestBody.AgendamentoId,
			Preco:         requestBody.Preco,
			Desenho:       requestBody.Desenho,
			Tamaho:        requestBody.Tamaho,
			Cor:           requestBody.Cor,
			Estilo:        requestBody.Estilo,
		}

		var results Tatuagem
		err := client.DB.From("tatuagens").Select("*").Eq("id", tatuagemId).Execute(&results)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"errror": err.Error()})
			return
		}

		updateErr := client.DB.From("tatuagens").Update(row).Eq("id", tatuagemId)

		if updateErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"errror": err.Error()})
		}

		c.JSON(http.StatusOK, results)
	})

	// Find all tattoo by tattoo artist
	router.GET("/tatuagens/:id", func(c *gin.Context) {
		// extract of param the tatuador id
		tatuadorId := c.Param("id")

		// variable of return function execute databse
		var results []Tatuagem

		err := client.DB.From("tatuagens").Select("*").Eq("tatuador_id", tatuadorId).Execute(&results)

		// tratament error case exists
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// response
		c.JSON(http.StatusOK, results)
	})

	// Delete tattoo per id tattoo artist
	router.DELETE("/tatuagens/:id", func(c *gin.Context) {
		tatuagemId := c.Param("id")

		var results Tatuagem
		err := client.DB.From("tatuagens").Delete().Eq("id", tatuagemId).Execute(&results)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, results)
	})

	/*********************************************************
	* 				   	  CRUD USUARIOS 				   	 *
	**********************************************************/
	// Define CRUD routes for "usuarios"
	router.POST("/usuarios", func(c *gin.Context) {

		// Create a new usuario
		row := Usuario{
			Nome: "Gabriel Medrado",
		}
		var results []Usuario
		client := supabase.CreateClient(supabaseURL, supabaseKey)
		err := client.DB.From("usuarios").Insert(row).Execute(&results)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, results)
	})

	// router.GET("/usuarios/:id", func(c *gin.Context) {
	// 	// Get a usuario by ID
	// 	tatuadorID := c.Param("id")
	// 	resp, err := client.
	// 		Table("usuarios").
	// 		Select("*").
	// 		Equals("id", tatuadorID).
	// 		Single()
	// 	if err != nil {
	// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	// 		return
	// 	}

	// 	c.JSON(http.StatusOK, resp)
	// })

	// router.PUT("/usuarios/:id", func(c *gin.Context) {
	// 	// Update a usuario by ID
	// 	tatuadorID := c.Param("id")
	// 	var usuario Usuario // Define the Usuario struct as needed
	// 	if err := c.BindJSON(&usuario); err != nil {
	// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	// 		return
	// 	}

	// 	resp, err := client.
	// 		Table("usuarios").
	// 		Update(tatuadorToMap(usuario)).
	// 		Equals("id", tatuadorID).
	// 		Single()
	// 	if err != nil {
	// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	// 		return
	// 	}

	// 	c.JSON(http.StatusOK, resp)
	// })

	// router.DELETE("/usuarios/:id", func(c *gin.Context) {
	// 	// Delete a usuario by ID
	// 	tatuadorID := c.Param("id")
	// 	resp, err := client.
	// 		Table("usuarios").
	// 		Delete().
	// 		Equals("id", tatuadorID).
	// 		Single()
	// 	if err != nil {
	// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	// 		return
	// 	}

	// 	c.JSON(http.StatusOK, resp)
	// })

	// Start the Gin server
	port := 8080 // Change to the desired port
	router.Run(fmt.Sprintf(":%d", port))
}

// Define the Usuario struct to match your database structure
type Usuario struct {
	Nome string `json:"nome"`
}

type Tatuagem struct {
	TatuadorId    int     `json:"tatuador_id"`
	AgendamentoId int     `json:"agendamento_id"`
	Preco         float32 `json:"preco"`
	Desenho       string  `json:"desenho"`
	Tamaho        int     `json:"tamanho"`
	Cor           string  `json:"cor"`
	Estilo        string  `json:"estilo"`
}

// Helper function to convert Usuario struct to map for Supabase
func tatuadorToMap(usuario Usuario) map[string]interface{} {
	return map[string]interface{}{
		"nome": usuario.Nome,
	}
}
