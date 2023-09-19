package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"os"
	supabase "github.com/nedpals/supabase-go"
)

func main() {
	if err := godotenv.Load(); err != nil {
        fmt.Printf("Erro ao carregar arquivo .env: %v\n", err)
        os.Exit(1)
    }
	// Initialize your Supabase client
	supabaseURL := os.Getenv("SUPABASE_URL")
	supabaseKey := os.Getenv("SUPABASE_KEY")

	// Create a Gin router
	router := gin.Default()

	//Initialize a single supabase client instead of one for each query received
	client := supabase.CreateClient(supabaseURL, supabaseKey)

	extractBearerToken := func(header string) (string, error) {
		if header == "" {
			return "", errors.New("Missing authorization header")
		}

		jwtToken := strings.Split(header, " ")
		if len(jwtToken) != 2 {
			return "", errors.New("Incorrectly formatted authorization header")
		}

		return jwtToken[1], nil
	}

	jwtTokenCheck := func(c *gin.Context) {
		jwtToken, err := extractBearerToken(c.GetHeader("Authorization"))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		client.DB.AddHeader("Authorization", "Bearer "+jwtToken)
		c.Next()
	}

	// Create a group, all routes initialized with this group will pass through the
	// jwtTokenCheck middleware function and be located like: /private/...
	private := router.Group("/private", jwtTokenCheck);
	
	//Initialize a single supabase client instead of one for each query received

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

		fmt.Println(user)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"user": user})
	})

	// Define CRUD routes for "usuarios"
	private.POST("/usuarios", func(c *gin.Context) {
		// Create a new usuario
		var row Usuario

		if errBind := c.ShouldBindJSON(&row); errBind != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": errBind.Error()})
			return
		}

		var results []Usuario
		errInsert := client.DB.From("usuarios").Insert(row).Execute(&results)

		if errInsert != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": errInsert.Error()})
			return
		}

		c.JSON(http.StatusCreated, results)
	})

	router.GET("/usuarios/:id", func(c *gin.Context) {
		id := c.Param("id")
		var user Usuario
		err := client.DB.From("usuarios").Select("*").Single().Eq("id", id).Execute(&user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, user)
	})

	private.PATCH("/usuarios/:id", func(c *gin.Context) {
		id := c.Param("id")
		var user Usuario
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		var results []Usuario
		err := client.DB.From("usuarios").Update(user).Eq("id", id).Execute(&results)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}
		c.JSON(http.StatusOK, user)
	})

	router.GET("/tatuadores/:tatuador_id", func(c *gin.Context) {
		tatuador_id := c.Param("tatuador_id")
		var tatuador Tatuador
		err := client.DB.From("tatuadores").Select("*").Single().Eq("id", tatuador_id).Execute(&tatuador)

		fmt.Println(tatuador) // Selected rows

		if err != nil {
			c.JSON(http.StatusNoContent, gin.H{"caiu": err.Error()})
			return
		}
	
		c.JSON(http.StatusOK, tatuador)
	})

	router.GET("/tatuadores", func(c *gin.Context) {
		var listaTatuadores []Tatuador
		err := client.DB.From("tatuadores").Select("*").Execute(&listaTatuadores)
		if err != nil {
			panic(err)
		}

		fmt.Println(listaTatuadores)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	
		c.JSON(http.StatusOK, listaTatuadores)
	})

	private.POST("/tatuadores", func(c *gin.Context) {
		var tatuador Tatuador

		if err := c.BindJSON(&tatuador); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Erro ao decodificar JSON"})
			return
		}
		
		var results []Tatuador

		err := client.DB.From("tatuadores").Insert(tatuador).Execute(&results)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"ruim": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, results)
	})

	private.PATCH("/tatuadores/:tatuador_id", func(c *gin.Context) {
		tatuador_id := c.Param("tatuador_id")
		var tatuador Tatuador

		if err := c.BindJSON(&tatuador); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Erro ao decodificar JSON"})
			return
		}
		var results []Tatuador
		err := client.DB.From("tatuadores").Update(tatuador).Eq("id", tatuador_id).Execute(&results)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"ruim": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, tatuador)
	})

	// Start the Gin server
	port := 8080 // Change to the desired port
	router.Run(fmt.Sprintf(":%d", port))
}

// Define the Usuario struct to match your database structure
type Usuario struct {
	Nome string `json:"nome"`
	Email string `json:"email"`
	TelefoneCelular string `json:"telefone_celular"`
	Cpf string `json:"cpf"`
	Rg string `json:"rg"`
	Status string `json:"status"`
	Endereco string `json:"endereco"`
}

type Tatuador struct {
	EstudioId int `json:"estudio_id"`
	Experiencia int `json:"experiencia"`
	EstiloTatuagem string `json:"estilo_tatuagem"`
	Status string `json:"status"`
	Tipo string `json:"tipo"`
	RedesSociais   *struct {
		Instagram string `json:"instagram"`
		X         string `json:"x"`
		Facebook  string `json:"facebook"`
	} `json:"redes_sociais"`
}
