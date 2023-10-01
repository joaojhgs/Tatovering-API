package controllers
import (
	"fmt"
	"net/http"
	"reflect"
	"github.com/gin-gonic/gin"
	supabase "github.com/nedpals/supabase-go"
	"tatovering/src/models"
)

func GetByIdTatuador(client *supabase.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		tatuador_id := c.Param("tatuador_id")
	
		var tatuador models.Tatuador

		err := client.DB.From("tatuadores").Select("*").Single().Eq("id", tatuador_id).Execute(&tatuador)

		fmt.Println(tatuador)

		if err != nil {
			c.JSON(http.StatusNoContent, gin.H{"caiu": err.Error()})
			return
		}

		c.JSON(http.StatusOK, tatuador)
	}
}

func ListagemTatuador(client *supabase.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		var listaTatuadores []models.Tatuador

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
	}
}

func CadastroTatuador(client *supabase.Client) gin.HandlerFunc {
	return func(c *gin.Context) {

		var tatuador models.Tatuador

		if err := c.BindJSON(&tatuador); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Erro ao decodificar JSON"})
			return
		}

		var results []models.Tatuador

		err := client.DB.From("tatuadores").Insert(tatuador).Execute(&results)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"ruim": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, results)
	}
}

func EdicaoTatuador(client *supabase.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		var tatuadorAtual models.Tatuador
		err := client.DB.From("tatuadores").Select("*").Single().Eq("id", id).Execute(&tatuadorAtual); 
		
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"bla": err.Error()})
			return
		}
		// BindJSON tentará analisar o corpo da solicitação JSON na variável 'user'
		var tatuadorUpdate models.Tatuador
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
		var results []models.Tatuador

		errUpdate := client.DB.From("tatuadores").Update(tatuadorAtual).Eq("id", id).Execute(&results)

		if errUpdate != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"blablabla": errUpdate.Error()})
			return
		}

		c.JSON(http.StatusOK, tatuadorAtual)
	}
}