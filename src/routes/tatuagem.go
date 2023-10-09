package routes

import (
    "github.com/gin-gonic/gin"
    "tatovering/src/controllers"
    supabase "github.com/nedpals/supabase-go"
    "tatovering/src/middlewares"
)

func SetupTatuagemRoutes(router *gin.Engine, client *supabase.Client) {
    tatuagemGroup := router.Group("/tatuagens")
    {
        tatuagemGroup.GET(
            "/",
            controllers.ListagemTatuagem(client),
        )
        tatuagemGroup.GET("/:id",
            controllers.GetByIdTatuagem(client),
        )
        tatuagemGroup.GET("/tatuador/:tatuador_id",
            controllers.GetArtByTatuadorId(client),
        )
        tatuagemGroup.POST(
            "/",
            middlewares.JwtTokenCheck(client),
            controllers.CadastrarTatuagem(client),
        )
        tatuagemGroup.PATCH(
            "/:id",
            middlewares.JwtTokenCheck(client),
            controllers.EditarTatuagem(client),
        )
        tatuagemGroup.DELETE(
            "/:id",
            middlewares.JwtTokenCheck(client),
            controllers.DeletarTatuagem(client),
        )
    }
}