package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	supabase "github.com/nedpals/supabase-go"
)

func main() {
	// Initialize your Supabase client
	supabaseURL := "https://pfzlboeaonsookzcnniv.supabase.co"
	supabaseKey := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzdXBhYmFzZSIsInJlZiI6InBmemxib2Vhb25zb29remNubml2Iiwicm9sZSI6ImFub24iLCJpYXQiOjE2OTIwNjI3NTYsImV4cCI6MjAwNzYzODc1Nn0.KuEEX9EBIQmLTA02iPtqqNIewDmXITDxnIfD4qEqTN8"

	extractBearerToken := func(header string) (string, error) {
		if header == "" {
			return "", errors.New("bad header value given")
		}

		jwtToken := strings.Split(header, " ")
		if len(jwtToken) != 2 {
			return "", errors.New("incorrectly formatted authorization header")
		}

		return jwtToken[1], nil
	}

	jwtTokenCheck := func(c *gin.Context) {
		jwtToken, err := extractBearerToken(c.GetHeader("Authorization"))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		client := supabase.CreateClient(supabaseURL, supabaseKey)
		client.DB.AddHeader("Authorization", "Bearer "+jwtToken)
		c.Next()
	}

	// Create a Gin router
	router := gin.Default()

	// Create a group, all routes initialized with this group will pass through the
	// jwtTokenCheck middleware function and be located like: /private/...
	private := router.Group("/private", jwtTokenCheck)

	//Initialize a single supabase client instead of one for each query received
	client := supabase.CreateClient(supabaseURL, supabaseKey)
	// Route for user sign-up
	router.POST("/signup", func(c *gin.Context) {
		// Defines the input data and validation
		var requestBody struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}
		// Bind the request to the defined model and throw error if some validation fails.
		if err := c.ShouldBindJSON(&requestBody); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Create user credentials
		credentials := supabase.UserCredentials{
			Email:    requestBody.Email,
			Password: requestBody.Password,
		}
		ctx := context.Background()
		// Sign up the user with Supabase
		user, err := client.Auth.SignUp(ctx, credentials)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"user": user})
	})

	// Route for user sign-in
	router.POST("/signin", func(c *gin.Context) {
		var requestBody struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}

		if err := c.ShouldBindJSON(&requestBody); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Create user credentials
		credentials := supabase.UserCredentials{
			Email:    requestBody.Email,
			Password: requestBody.Password,
		}

		ctx := context.Background()
		// Sign up the user with Supabase
		user, err := client.Auth.SignIn(ctx, credentials)

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"user": user})
	})

	// Define CRUD routes for "usuarios"
	private.POST("/usuarios", func(c *gin.Context) {
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

	// Estudios Logic

	// Get estudio by ID
	router.GET("/estudios/:id", func(c *gin.Context) {

		var estudioId = c.Param("id")

		var estudio interface{}

		var err = client.DB.From("estudios").Select("*").Single().Eq("id", estudioId).Execute(&estudio)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, estudio)
	})

	// Get estudios
	router.GET("/estudios", func(c *gin.Context) {

		var estudios []interface{}

		var err = client.DB.From("estudios").Select("*").Execute(&estudios)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, estudios)
	})

	// Create a estudio
	router.POST("/estudios", func(c *gin.Context) {

		// type Estudio struct {
		// 	ProprietarioId       int     `json:"proprietario_id"`
		// 	Nome                 string  `json:"nome"`
		// 	Email                string  `json:"email"`
		// 	HorarioFuncionamento string  `json:"horario_funcionamento"`
		// 	Endereco             string  `json:"endereco"`
		// 	Localizacao          string  `json:"localizacao"`
		// 	Telefone             string  `json:"telefone"`
		// 	Descricao            string  `json:"descricao"`
		// 	TaxaAgendamento      float32 `json:"taxa_agendamento"`
		// }

		// var requestBody = Estudio{}
		var requestBody interface{}

		err := c.ShouldBindJSON(&requestBody)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// var body = Estudio{
		// 	ProprietarioId:       requestBody.ProprietarioId,
		// 	Nome:                 requestBody.Nome,
		// 	Email:                requestBody.Email,
		// 	HorarioFuncionamento: requestBody.HorarioFuncionamento,
		// 	Endereco:             requestBody.Endereco,
		// 	Localizacao:          requestBody.Localizacao,
		// 	Telefone:             requestBody.Telefone,
		// 	Descricao:            requestBody.Descricao,
		// 	TaxaAgendamento:      requestBody.TaxaAgendamento,
		// }

		var result interface{}
		erro := client.DB.From("estudios").Insert(requestBody).Execute(&result)

		if erro != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": erro.Error()})
			return
		}

		c.JSON(http.StatusOK, result)
	})

	router.DELETE("estudios/:id", func(c *gin.Context) {
		var estudioId = c.Param("id")

		var result interface{}

		var errSelect = client.DB.From("estudios").Select("*").Single().Eq("id", estudioId).Execute(&result)
		if errSelect != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": errSelect.Error()})
			return
		}

		var deleteReturn interface{}
		var err = client.DB.From("estudios").Delete().Eq("id", estudioId).Execute(&deleteReturn)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, result)
	})

	// Start the Gin server
	port := 8080 // Change to the desired port
	router.Run(fmt.Sprintf(":%d", port))
}

// Define the Usuario struct to match your database structure
type Usuario struct {
	Nome string `json:"nome"`
}

// Helper function to convert Usuario struct to map for Supabase
func tatuadorToMap(usuario Usuario) map[string]interface{} {
	return map[string]interface{}{
		"nome": usuario.Nome,
	}
}
