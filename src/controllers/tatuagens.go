package controllers

import (
	"fmt"
	"reflect"
	"net/http"
	"github.com/gin-gonic/gin"
	supabase "github.com/nedpals/supabase-go"
	"tatovering/src/models"
)

// Tatuagens Ok

func CadastrarTatuagem(client *supabase.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		var requestBody models.Tatuagem

		if err := c.ShouldBindJSON(&requestBody); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var results []models.Tatuagem
		fmt.Println(requestBody)
		err := client.DB.From("tatuagens").Insert(requestBody).Execute(&results)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, results)
	}
}

func GetByIdTatuagem(client *supabase.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		tatuagemId := c.Param("id")

		var results []models.Tatuagem

		err := client.DB.From("tatuagens").Select("*").Eq("id", tatuagemId).Execute(&results)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, results)
	}
}

func ListagemTatuagem(client *supabase.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		var listaTatuagens []models.Tatuagem
		err := client.DB.From("tatuagens").Select("*").Execute(&listaTatuagens)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
			}
	
		c.JSON(http.StatusOK, listaTatuagens)
	}
}

func DeletarTatuagem(client *supabase.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		tatuagemId := c.Param("id")

		var results models.Tatuagem
		err := client.DB.From("tatuagens").Delete().Eq("id", tatuagemId).Execute(&results)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, results)
	}
}

func EditarTatuagem(client *supabase.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		var tatuagemAtual models.Tatuagem
		err := client.DB.From("tatuagens").Select("*").Single().Eq("id", id).Execute(&tatuagemAtual); 
		
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"bla": err.Error()})
			return
		}
		// BindJSON tentará analisar o corpo da solicitação JSON na variável 'user'
		var tatuagemUpdate models.Tatuagem
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
		var results []models.Tatuagem
	
		errUpdate := client.DB.From("tatuagens").Update(tatuagemAtual).Eq("id", id).Execute(&results)

		if errUpdate != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"blablabla": errUpdate.Error()})
			return
		}

		c.JSON(http.StatusOK, tatuagemAtual)
	}
}

func GetArtByTatuadorId(client *supabase.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		tatuadorId := c.Param("tatuador_id")

		var tatuagens []models.Tatuagem

		errSelectTatuagen := client.DB.From("tatuagens").Select("*").Eq("tatuador_id", tatuadorId).Execute(&tatuagens)

		if errSelectTatuagen != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"bla": errSelectTatuagen.Error()})
			return
		}

		var tatuador []models.Tatuador

		errSelectTatuador := client.DB.From("tatuadores").Select("*").Eq("id", tatuadorId).Execute(&tatuador)

		if errSelectTatuador != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"blabla": errSelectTatuador.Error()})
			return
		}

		dadosRetorno := gin.H{
			"tatuador":  tatuador,
			"tatuagens": tatuagens,
		}

		c.JSON(http.StatusOK, dadosRetorno)

	}
}