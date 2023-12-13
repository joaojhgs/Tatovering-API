package controllers

import (
	"fmt"
	"net/http"
	"reflect"
	"tatovering/src/models"

	"github.com/gin-gonic/gin"
	supabase "github.com/nedpals/supabase-go"
)

func GetByIdTatuador(client *supabase.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		var tatuador models.Tatuador

		err := client.DB.From("tatuadores").Select("*").Single().Eq("id", id).Execute(&tatuador)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"caiu": err.Error()})
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

		var tatuador models.TatuadorPost

		if err := c.BindJSON(&tatuador); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Erro ao decodificar JSON"})
			return
		}

		var results []models.TatuadorPost

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
		err := client.DB.From("tatuadores").Select("*").Single().Eq("id", id).Execute(&tatuadorAtual)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"bla": err.Error()})
			return
		}
		// BindJSON tentará analisar o corpo da solicitação JSON na variável 'user'
		var tatuadorUpdate models.Tatuador
		errDadosUpdate := c.ShouldBindJSON(&tatuadorUpdate)

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

func TatuadoresEstudio(client *supabase.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		estudioId := c.Param("estudio_id")

		var tatuadores []models.Tatuador

		errSelectTatuadores := client.DB.From("tatuadores").Select("*").Eq("estudio_id", estudioId).Execute(&tatuadores)

		if errSelectTatuadores != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"bla": errSelectTatuadores.Error()})
			return
		}

		var estudio []models.Estudio

		errSelectEstudios := client.DB.From("estudios").Select("*").Eq("id", estudioId).Execute(&estudio)

		if errSelectEstudios != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"blabla": errSelectEstudios.Error()})
			return
		}

		dadosRetorno := gin.H{
			"tatuadores": tatuadores,
			"estudio":    estudio,
		}

		c.JSON(http.StatusOK, dadosRetorno)
	}
}
