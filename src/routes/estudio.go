package routes

import (
    "github.com/gin-gonic/gin"
    "tatovering/src/controllers"
    supabase "github.com/nedpals/supabase-go"
    "tatovering/src/middlewares"
)

func SetupEstudiosRoutes(router *gin.Engine, client *supabase.Client) {
    estudioGroup := router.Group("/estudios")
    {
        estudioGroup.GET("/", controllers.ListagemEstudio(client))
        estudioGroup.GET("/:id", controllers.GetByIdEstudio(client))
        estudioGroup.POST(
            "/",
            middlewares.JwtTokenCheck(client),
            controllers.CadastrarEstudio(client), 
        )
        estudioGroup.PATCH(
            "/:id",
            middlewares.JwtTokenCheck(client),
            controllers.EditarEstudio(client),
        )
        estudioGroup.DELETE(
            "/:id",
            middlewares.JwtTokenCheck(client),
            controllers.DeletarEstudio(client),
        )
    }
}