package main

import (
	"fmt"
	"net/http"
	"context"
	"errors"
	"strings"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"os"
	"github.com/jinzhu/copier"
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

	extractBearerToken := func (header string) (string, error) {
		if header == "" {
			return "", errors.New("bad header value given")
		}
	
		jwtToken := strings.Split(header, " ")
		if len(jwtToken) != 2 {
			return "", errors.New("incorrectly formatted authorization header")
		}
	
		return jwtToken[1], nil
	}

	jwtTokenCheck := func (c *gin.Context) {
		jwtToken, err := extractBearerToken(c.GetHeader("Authorization"))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		client.DB.AddHeader("Authorization", "Bearer "+jwtToken);
		c.Next()
	}

	// Create a group, all routes initialized with this group will pass through the 
	// jwtTokenCheck middleware function and be located like: /private/...
	private := router.Group("/private", jwtTokenCheck);
	
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
		var usuario Usuario

		// BindJSON tentará analisar o corpo da solicitação JSON na variável 'usuario'
		if err := c.BindJSON(&usuario); err != nil {
			// Em caso de erro, responder com status de erro e mensagem
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var results []Usuario
		client := supabase.CreateClient(supabaseURL, supabaseKey)
		err := client.DB.From("usuarios").Insert(usuario).Execute(&results)

			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			c.JSON(http.StatusCreated, results)
	})

	router.POST("/tatuadores", func(c *gin.Context) {
		client := supabase.CreateClient(supabaseURL, supabaseKey)
		// Crie uma variável para armazenar os dados do usuário a partir do corpo da solicitação

		var data map[string]map[string]interface{}
		
		// BindJSON tentará analisar o corpo da solicitação JSON na variável 'usuario'
		if err := c.BindJSON(&data); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Erro ao decodificar JSON"})
			return
		}

		usuarioData := data["usuario"]
		tatuadorData := data["tatuador"] 
		
		fmt.Println("AQUI", usuarioData)
		fmt.Println("AQUI", tatuadorData)
		
		var usuario Usuario
		var tatuador Tatuador

		if usuarioData != nil {
			// Use o pacote github.com/jinzhu/copier para copiar os valores do mapa para a struct
			if err := copier.Copy(&usuario, usuarioData); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao copiar dados para a struct de usuario"})
				return
			}
		}

		if tatuadorData != nil {
			// Use o pacote github.com/jinzhu/copier para copiar os valores do mapa para a struct
			if err := copier.Copy(&tatuador, tatuadorData); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao copiar dados para a struct de tatuador"})
				return
			}
		}

		var resul []Usuario

		err := client.DB.From("usuarios").Insert(usuario).Execute(&resul)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"caca": err.Error()})
				return
			} else {
				var results []Tatuador

				err := client.DB.From("tatuadores").Insert(tatuador).Execute(&results)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"ruim": err.Error()})
					return
				}

				c.JSON(http.StatusCreated, results)
			}
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
	UsuarioUuid string `json:"usuario_uuid"`
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
