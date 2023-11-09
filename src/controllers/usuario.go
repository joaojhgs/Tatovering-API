package controllers

import (
	"reflect"
	"fmt"
	"net/http"
	"github.com/gin-gonic/gin"
	supabase "github.com/nedpals/supabase-go"
	"tatovering/src/models"
)

func GetById(client *supabase.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		usuarioId := c.Param("id")

		var existingUser []models.Usuario

		err := client.DB.From("usuarios").Select("*").Eq("id", usuarioId).Execute(&existingUser);

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, existingUser)
	}
}

func CadastrarUsuario(client *supabase.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		var usuario models.Usuario

		if errBind := c.ShouldBindJSON(&usuario); errBind != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": errBind.Error()})
			return
		}

		var results []models.Usuario
		errInsert := client.DB.From("usuarios").Insert(usuario).Execute(&results)

		if errInsert != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": errInsert.Error()})
			return
		}

		c.JSON(http.StatusCreated, results)
	}
}

func DeletarUsuario(client *supabase.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		var usuario models.Usuario
		err := client.DB.From("usuarios").Delete().Eq("id", id).Execute(&usuario)
		if err != nil {
			panic(err)
		}

		fmt.Println(usuario) // Empty - nothing returned from delete
		c.JSON(http.StatusOK, usuario)
	}
}

func EditarUsuario(client *supabase.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var existingUser models.Usuario
		err := client.DB.From("usuarios").Select("*").Single().Eq("id", id).Execute(&existingUser); 
		
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		// BindJSON tentará analisar o corpo da solicitação JSON na variável 'user'
		var user models.Usuario
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
		var results []models.Usuario
		errUpdate := client.DB.From("usuarios").Update(existingUser).Eq("id", id).Execute(&results)

		if errUpdate != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, existingUser)
	}
}