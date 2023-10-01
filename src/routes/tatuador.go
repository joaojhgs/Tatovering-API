package routes

import (
    "github.com/gin-gonic/gin"
    "tatovering/src/controllers"
    supabase "github.com/nedpals/supabase-go"
)

func SetupTatuadoresRoutes(router *gin.Engine, client *supabase.Client) {
    tatuadoreGroup := router.Group("/tatuadores")
    {
        tatuadoreGroup.GET("/", controllers.ListagemTatuador(client))
        tatuadoreGroup.GET("/:id", controllers.GetByIdTatuador(client))
        tatuadoreGroup.POST("/", controllers.CadastroTatuador(client))
        tatuadoreGroup.PATCH("/:id", controllers.EdicaoTatuador(client))
    }
}