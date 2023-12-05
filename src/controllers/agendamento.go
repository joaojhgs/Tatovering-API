package controllers

import (
	"fmt"
	"net/http"
	"reflect"
	extract "tatovering/src/middlewares"
	models "tatovering/src/models/cadastros"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/mitchellh/mapstructure"
	supabase "github.com/nedpals/supabase-go"
)

func EfetuarAgendamentoUsuario(client *supabase.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extrair token do cabeçalho da autorização
		token, erroToken := extract.ExtractBearerToken(c.GetHeader("Authorization"))
		if erroToken != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": erroToken.Error()})
			return
		}

		// Decodificar o token
		claims, err := decodeToken(token)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Preencher os objetos agendamentoUsuario e servicoUsuario
		var dados map[string]interface{}
		if err := c.BindJSON(&dados); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Erro ao decodificar JSON"})
			return
		}

		fmt.Println(222, dados)
		agendamentoUsuarioMap, ok := dados["agendamento"].(map[string]interface{})
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Chave 'agendamento' não encontrada ou não é um mapa"})
			return
		}
		fmt.Println(1.5, agendamentoUsuarioMap)

		servicoUsuarioMap, ok := dados["servico"].(map[string]interface{})
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Chave 'servico' não encontrada ou não é um mapa"})
			return
		}

		// Converte os mapas para structs
		var agendamentoUsuario models.CadastroAgendamentoUsuario
		var servicoUsuario models.CadastroServicoUsuario

		// Converte o mapa do serviço para a struct correspondente
		if err := mapstructure.Decode(servicoUsuarioMap, &servicoUsuario); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Erro ao converter dados do serviço"})
			return
		}

		// Converte o mapa do agendamento para a struct correspondente
		fmt.Println(2.0, agendamentoUsuarioMap)
		if err := mapstructure.Decode(agendamentoUsuarioMap, &agendamentoUsuario); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Erro ao converter dados do agendamento"})
			return
		}

		fmt.Println(45)
		fmt.Println(2.5, agendamentoUsuario)

		subClaim, ok := claims["sub"]

		if !ok {
			// Lidar com a ausência da chave "user"
			fmt.Println("Chave 'user' não encontrada nas reivindicações.")
			// ...
			return
		}

		// Tentar assertar o valor como uma string
		id, ok := subClaim.(string)

		if !ok {
			// Lidar com o caso em que o valor não é uma string
			fmt.Println("Valor da chave 'sub' não é uma string.")
			// ...
			return
		}

		tatuadorId := agendamentoUsuarioMap["tatuador_id"]
		estudioId := agendamentoUsuarioMap["estudio_id"]
		tatuagemId := servicoUsuarioMap["tatuagem_id"]

		tatuagemIdString, ok := tatuagemId.(string)

		str, ok := tatuadorId.(string)

		str2, ok := estudioId.(string)

		servicoid := uuid.New()
		agendamentoid := uuid.New()

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"ruimNoParseClienteId": err.Error()})
			return
		}

		// Definir cliente_id usando o ID do usuário decodificado
		agendamentoUsuario.Id = agendamentoid.String()
		agendamentoUsuario.ClienteId = id
		agendamentoUsuario.ServicoId = servicoid.String()
		agendamentoUsuario.TatuadorId = str
		agendamentoUsuario.EstudioId = str2

		// Executar a inserção no banco de dados
		var result models.CadastroAgendamentoUsuario

		result.ClienteId = agendamentoUsuario.ClienteId
		result.ServicoId = agendamentoUsuario.ServicoId

		tipo := reflect.TypeOf(servicoUsuario)

		// Iterar sobre os campos do tipo
		for i := 0; i < tipo.NumField(); i++ {
			// Obter o campo atual
			campo := tipo.Field(i)
			// Imprimir o nome do campo
			fmt.Printf("Campo: %s, Tipo: %s\n", campo.Name, campo.Type)
		}

		err = client.DB.From("agendamentos").Insert(agendamentoUsuario).Execute(&result)

		var resultServico models.CadastroServicoUsuario

		servicoUsuario.Id = servicoid.String()
		servicoUsuario.ClienteId = id
		servicoUsuario.TatuadorId = str
		servicoUsuario.EstudioId = str2
		servicoUsuario.TatuagemId = tatuagemIdString

		resultServico.Id = servicoid.String()

		fmt.Println(333, servicoUsuario)

		err = client.DB.From("servicos").Insert(servicoUsuario).Execute(&resultServico)

		c.JSON(http.StatusCreated, result)
	}
}

// Função para decodificar o token JWT
func decodeToken(token string) (jwt.MapClaims, error) {
	// Chave secreta para verificar a assinatura do token
	chaveSecreta := []byte("3xBRcSCY2xTUjfL+ELWskobjMqUFez0sCCGu9hxDfqacWL7FdbYb6bQlVAXK48hMoQYp0PeEy3eHzawk9/XJDA==")

	// Analisar e validar o token
	tokenObj, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		// Verifica o método de assinatura
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Método de assinatura inválido: %v", token.Header["alg"])
		}
		return chaveSecreta, nil
	})

	if err != nil {
		return nil, err
	}

	// Verificar se o token é válido
	if !tokenObj.Valid {
		return nil, fmt.Errorf("Token inválido")
	}

	// Acessar as reivindicações (claims) do token
	claims, ok := tokenObj.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("Erro ao obter reivindicações do token")
	}

	return claims, nil
}
