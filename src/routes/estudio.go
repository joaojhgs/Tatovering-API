package routes

import (
    "github.com/gin-gonic/gin"
    "tatovering/src/controllers"
    supabase "github.com/nedpals/supabase-go"
)

func SetupEstudiosRoutes(router *gin.Engine, client *supabase.Client) {
    estudioGroup := router.Group("/estudios")
    {
        estudioGroup.GET("/", controllers.ListagemEstudio(client))
        estudioGroup.GET("/:id", controllers.GetByIdEstudio(client))
        estudioGroup.POST("/", controllers.CadastrarEstudio(client))
        estudioGroup.PATCH("/:id", controllers.EditarEstudio(client))
        estudioGroup.DELETE("/:id", controllers.DeletarEstudio(client))
    }
}