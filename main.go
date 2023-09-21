package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	supabase "github.com/nedpals/supabase-go"
)

func Cors() gin.HandlerFunc {
    return func(ctx *gin.Context) {
        method := ctx.Request.Method
        if method == "OPTIONS" {
            ctx.Header("Access-Control-Max-Age", "1728000")
            ctx.Header("Access-Control-Allow-Credentials", "true")
            ctx.Header("Access-Control-Allow-Methods", "GET,PUT,POST,DELETE,PATCH,OPTIONS")
            ctx.Header("Access-Control-Allow-Origin", ctx.Request.Header.Get("origin"))
            ctx.Header("Access-Control-Allow-Headers", "Content-Type,Cookie,Authorization,Access-Control-Request-Headers,Access-Control-Request-Method,Origin,Referer,Sec-Fetch-Dest,Accept-Language,Accept-Encoding,Sec-Fetch-Mode,Sec-Fetch-Site,User-Agent,Pragma,Host,Connection,Cache-Control,Accept-Language,Accept-Encoding,X-Requested-With,X-Forwarded-For,X-Forwarded-Host,X-Forwarded-Proto,X-Forwarded-Port,X-Forwarded-Prefix,X-Real-IP,Accept")
            ctx.AbortWithStatus(http.StatusNoContent)
            return
        }
        ctx.Header("Access-Control-Allow-Origin", ctx.Request.Header.Get("origin"))
        ctx.Header("Access-Control-Allow-Credentials", "true")
        ctx.Next()
    }
}

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
	router.Use(Cors())
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
	private := router.Group("/", jwtTokenCheck)

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
	/*********************************************************
	* 				   	  CRUD TATUAGENS 				   	 *
	**********************************************************/
	router.GET("/tatuagens", func(c *gin.Context) {
		var listaTatuagens []Tatuagem
		err := client.DB.From("tatuagens").Select("*").Execute(&listaTatuagens)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
			}
	
		c.JSON(http.StatusOK, listaTatuagens)

	})

	private.POST("/tatuagens", func(c *gin.Context) {
		var requestBody Tatuagem

		if err := c.ShouldBindJSON(&requestBody); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var results []Tatuagem
		err := client.DB.From("tatuagens").Insert(requestBody).Execute(&results)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, results)
	})

	private.PATCH("/tatuagens/:id", func(c *gin.Context) {
		tatuagemId := c.Param("id")

		var requestBody Tatuagem

		if err := c.ShouldBindJSON(&requestBody); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var results Tatuagem
		err := client.DB.From("tatuagens").Select("*").Single().Eq("id", tatuagemId).Execute(&results)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		updateErr := client.DB.From("tatuagens").Update(requestBody).Eq("id", tatuagemId)

		if updateErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"errror": err.Error()})
		}

		c.JSON(http.StatusOK, results)
	})

	router.GET("/tatuagens/:id", func(c *gin.Context) {
		tatuadorId := c.Param("id")

		var results []Tatuagem

		err := client.DB.From("tatuagens").Select("*").Eq("tatuador_id", tatuadorId).Execute(&results)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, results)
	})

	private.DELETE("/tatuagens/:id", func(c *gin.Context) {
		tatuagemId := c.Param("id")

		var results Tatuagem
		err := client.DB.From("tatuagens").Delete().Eq("id", tatuagemId).Execute(&results)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, results)
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
		var listaEstudios []Estudio
		err := client.DB.From("estudios").Select("*").Execute(&listaEstudios)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
			}
	
		c.JSON(http.StatusOK, listaEstudios)

	})
	// Create a estudio
	router.POST("/estudios", func(c *gin.Context) {

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

	router.PATCH("estudios/:id", func(c *gin.Context) {
		var estudioId = c.Param("id")

		var requestBody interface{}

		errBody := c.ShouldBindJSON(&requestBody)
		if errBody != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": errBody.Error()})
			return
		}

		var result []interface{}
		var err = client.DB.From("estudios").Update(requestBody).Eq("id", estudioId).Execute(&result)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if len(result) == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Estudio não encontrado"})
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
	Nome            string `json:"nome"`
	Email           string `json:"email"`
	TelefoneCelular string `json:"telefone_celular"`
	Cpf             string `json:"cpf"`
	Rg              string `json:"rg"`
	Status          string `json:"status"`
	Endereco        string `json:"endereco"`
}

type Tatuador struct {
	EstudioId      int    `json:"estudio_id"`
	Experiencia    int    `json:"experiencia"`
	EstiloTatuagem string `json:"estilo_tatuagem"`
	Status         string `json:"status"`
	Tipo           string `json:"tipo"`
	RedesSociais   *struct {
		Instagram string `json:"instagram"`
		X         string `json:"x"`
		Facebook  string `json:"facebook"`
	} `json:"redes_sociais"`
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

type Estudio struct {
	ProprietarioId       int     `json:"proprietario_id"`
	Nome                 string  `json:"nome"`
	Email                string  `json:"email"`
	HorarioDeFuncionamento *struct {
    Segunda []string `json:"segunda"`
    Terca   []string `json:"terca"`
    Quarta  []string `json:"quarta"`
    Quinta  []string `json:"quinta"`
    Sexta   []string `json:"sexta"`
} `json:"horario_funcionamento"`
}
