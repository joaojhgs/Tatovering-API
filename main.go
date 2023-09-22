package main

import (
	"reflect"
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
	router.POST("/usuarios", func(c *gin.Context) {
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

	router.PATCH("/usuarios/:id", func(c *gin.Context) {
		id := c.Param("id")
		var existingUser Usuario
		err := client.DB.From("usuarios").Select("*").Single().Eq("id", id).Execute(&existingUser); 
		
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		// BindJSON tentará analisar o corpo da solicitação JSON na variável 'user'
		var user Usuario
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Usar reflexão para iterar sobre os campos da struct e atualizar os valores
		existingUserReflect := reflect.ValueOf(&existingUser).Elem()
		userReflect := reflect.ValueOf(&user).Elem()

		for i := 0; i < existingUserReflect.NumField(); i++ {
			fieldName := existingUserReflect.Type().Field(i).Name
			existingFieldValue := existingUserReflect.Field(i)
			userFieldValue := userReflect.FieldByName(fieldName)

			// Verifique se o campo existe na solicitação e é diferente de zero
			if userFieldValue.IsValid() && userFieldValue.Interface() != reflect.Zero(userFieldValue.Type()).Interface() {
				existingFieldValue.Set(userFieldValue)
			}
		}

		// Atualize o registro no banco de dados
		var results []Usuario
		errUpdate := client.DB.From("usuarios").Update(existingUser).Eq("id", id).Execute(&results)

		if errUpdate != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, existingUser)
	})

	router.DELETE("/usuarios/:id", func(c *gin.Context) {
		id := c.Param("id")

		var usuario Usuario
		err := client.DB.From("usuarios").Delete().Eq("id", id).Execute(&usuario)
		if err != nil {
			panic(err)
		}

		fmt.Println(usuario) // Empty - nothing returned from delete
		c.JSON(http.StatusOK, usuario)
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

	router.PATCH("/tatuadores/:id", func(c *gin.Context) {
		id := c.Param("id")

		var tatuadorAtual Tatuador
		err := client.DB.From("tatuadores").Select("*").Single().Eq("id", id).Execute(&tatuadorAtual); 
		
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"bla": err.Error()})
			return
		}
		// BindJSON tentará analisar o corpo da solicitação JSON na variável 'user'
		var tatuadorUpdate Tatuador
		errDadosUpdate := c.ShouldBindJSON(&tatuadorUpdate);
		
		if errDadosUpdate != nil {
			c.JSON(http.StatusBadRequest, gin.H{"blabla": errDadosUpdate.Error()})
			return
		}

		// Usar reflexão para iterar sobre os campos da struct e atualizar os valores
		tatuadorAtualReflect := reflect.ValueOf(&tatuadorAtual).Elem()
		tatuadorUpdateReflect := reflect.ValueOf(&tatuadorUpdate).Elem()

		for i := 0; i < tatuadorAtualReflect.NumField(); i++ {
			fieldName := tatuadorAtualReflect.Type().Field(i).Name
			atualCampoValor := tatuadorAtualReflect.Field(i)
			updateCampoValor := tatuadorUpdateReflect.FieldByName(fieldName)

			// Verifique se o campo existe na solicitação e é diferente de zero
			if updateCampoValor.IsValid() && updateCampoValor.Interface() != reflect.Zero(updateCampoValor.Type()).Interface() {
				atualCampoValor.Set(updateCampoValor)
			}
		}

		// Atualize o registro no banco de dados
		var results []Tatuador
	
		errUpdate := client.DB.From("tatuadores").Update(tatuadorAtual).Eq("id", id).Execute(&results)

		if errUpdate != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"blablabla": errUpdate.Error()})
			return
		}

		c.JSON(http.StatusOK, tatuadorAtual)
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

	router.PATCH("/tatuagens/:id", func(c *gin.Context) {
		id := c.Param("id")

		var tatuagemAtual Tatuagem
		err := client.DB.From("tatuagens").Select("*").Single().Eq("id", id).Execute(&tatuagemAtual); 
		
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"bla": err.Error()})
			return
		}
		// BindJSON tentará analisar o corpo da solicitação JSON na variável 'user'
		var tatuagemUpdate Tatuagem
		errDadosUpdate := c.ShouldBindJSON(&tatuagemUpdate);
		
		if errDadosUpdate != nil {
			c.JSON(http.StatusBadRequest, gin.H{"blabla": errDadosUpdate.Error()})
			return
		}

		// Usar reflexão para iterar sobre os campos da struct e atualizar os valores
		tatuagemAtualReflect := reflect.ValueOf(&tatuagemAtual).Elem()
		tatuagemUpdateReflect := reflect.ValueOf(&tatuagemUpdate).Elem()

		for i := 0; i < tatuagemAtualReflect.NumField(); i++ {
			fieldName := tatuagemAtualReflect.Type().Field(i).Name
			atualCampoValor := tatuagemAtualReflect.Field(i)
			updateCampoValor := tatuagemUpdateReflect.FieldByName(fieldName)

			// Verifique se o campo existe na solicitação e é diferente de zero
			if updateCampoValor.IsValid() && updateCampoValor.Interface() != reflect.Zero(updateCampoValor.Type()).Interface() {
				atualCampoValor.Set(updateCampoValor)
			}
		}

		// Atualize o registro no banco de dados
		var results []Tatuagem
	
		errUpdate := client.DB.From("tatuagens").Update(tatuagemAtual).Eq("id", id).Execute(&results)

		if errUpdate != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"blablabla": errUpdate.Error()})
			return
		}

		c.JSON(http.StatusOK, tatuagemAtual)
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
	Preco         float64 `json:"preco"`
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
