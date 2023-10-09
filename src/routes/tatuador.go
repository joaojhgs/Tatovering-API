package routes

import (
    "github.com/gin-gonic/gin"
    "tatovering/src/controllers"
    supabase "github.com/nedpals/supabase-go"
    "tatovering/src/middlewares"
)

func SetupTatuadoresRoutes(router *gin.Engine, client *supabase.Client) {
    tatuadoreGroup := router.Group("/tatuadores")
    {
        tatuadoreGroup.GET(
            "/",
            controllers.ListagemTatuador(client),
        )
        tatuadoreGroup.GET(
            "/:id",
            controllers.GetByIdTatuador(client),
        )
        tatuadoreGroup.GET(
            "/estudios/:estudio_id",
            controllers.TatuadoresEstudio(client),
        )
        tatuadoreGroup.POST(
            "/",
            middlewares.JwtTokenCheck(client),
            controllers.CadastroTatuador(client),
        )
        tatuadoreGroup.PATCH(
            "/:id",
            middlewares.JwtTokenCheck(client),
            controllers.EdicaoTatuador(client),
        )
    }
}