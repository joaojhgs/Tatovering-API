package routes

import (
    "github.com/gin-gonic/gin"
    "tatovering/src/controllers"
    supabase "github.com/nedpals/supabase-go"
)

func SetupTatuagemRoutes(router *gin.Engine, client *supabase.Client) {
    tatuagemGroup := router.Group("/tatuagens")
    {
        tatuagemGroup.GET("/", controllers.ListagemTatuagem(client))
        tatuagemGroup.GET("/:id", controllers.GetByIdTatuagem(client))
        tatuagemGroup.POST("/", controllers.CadastrarTatuagem(client))
        tatuagemGroup.PATCH("/:id", controllers.EditarTatuagem(client))
        tatuagemGroup.DELETE("/:id", controllers.DeletarTatuagem(client))
    }
}
