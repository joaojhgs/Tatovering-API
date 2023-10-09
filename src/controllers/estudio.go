package controllers

import (
	"fmt"
	"net/http"
	"reflect"
	"github.com/gin-gonic/gin"
	supabase "github.com/nedpals/supabase-go"
	"tatovering/src/models"
)

func ListagemEstudio(client *supabase.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		var listaEstudios []models.Estudio

		err := client.DB.From("estudios").Select("*").Execute(&listaEstudios)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
			}
	
		c.JSON(http.StatusOK, listaEstudios)
	}
}

func GetByIdEstudio(client *supabase.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		var estudioId = c.Param("id")

		var estudio models.Estudio

		var err = client.DB.From("estudios").Select("*").Single().Eq("id", estudioId).Execute(&estudio)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, estudio)
	}
}

func CadastrarEstudio(client *supabase.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		var novoEstudio models.Estudio

		err := c.ShouldBindJSON(&novoEstudio)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var result []models.Estudio
		erro := client.DB.From("estudios").Insert(novoEstudio).Execute(&result)

		if erro != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": erro.Error()})
			return
		}

		c.JSON(http.StatusOK, result)
	}
}

func EditarEstudio(client *supabase.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		var estudioAtual models.Estudio
		err := client.DB.From("estudios").Select("*").Single().Eq("id", id).Execute(&estudioAtual); 
		fmt.Println(3, estudioAtual)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"bla": err.Error()})
			return
		}
		// BindJSON tentará analisar o corpo da solicitação JSON na variável 'user'
		var estudioUpdate models.Estudio
		errDadosUpdate := c.ShouldBindJSON(&estudioUpdate);
		
		if errDadosUpdate != nil {
			c.JSON(http.StatusBadRequest, gin.H{"blabla": errDadosUpdate.Error()})
			return
		}

		// Usar reflexão para iterar sobre os campos da struct e atualizar os valores
		estudioAtualReflect := reflect.ValueOf(&estudioAtual).Elem()
		estudioUpdateReflect := reflect.ValueOf(&estudioUpdate).Elem()

		for i := 0; i < estudioAtualReflect.NumField(); i++ {
			fieldName := estudioAtualReflect.Type().Field(i).Name
			atualCampoValor := estudioAtualReflect.Field(i)
			updateCampoValor := estudioUpdateReflect.FieldByName(fieldName)

			// Verifique se o campo existe na solicitação e é diferente de zero
			if updateCampoValor.IsValid() && updateCampoValor.Interface() != reflect.Zero(updateCampoValor.Type()).Interface() {
				atualCampoValor.Set(updateCampoValor)
			}
		}

		// Atualize o registro no banco de dados
		var results []models.Tatuador
		fmt.Println(88, estudioAtual)
		fmt.Println(99, estudioUpdate)
		errUpdate := client.DB.From("estudios").Update(estudioAtual).Eq("id", id).Execute(&results)

		if errUpdate != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"blablabla": errUpdate.Error()})
			return
		}

		c.JSON(http.StatusOK, estudioAtual)
	}
}

func DeletarEstudio(client *supabase.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		var estudioId = c.Param("id")

		var estudio models.Estudio

		var errSelect = client.DB.From("estudios").Select("*").Single().Eq("id", estudioId).Execute(&estudio)
		if errSelect != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": errSelect.Error()})
			return
		}

		var listaEstudio []models.Estudio
		var err = client.DB.From("estudios").Delete().Eq("id", estudioId).Execute(&listaEstudio)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, listaEstudio)
	}
}
