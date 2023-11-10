package routes

import (
	"tatovering/src/controllers"
	"tatovering/src/middlewares"

	"github.com/gin-gonic/gin"
	supabase "github.com/nedpals/supabase-go"
)

func SetupTatuadoresRoutes(router *gin.Engine, client *supabase.Client) {
	tatuadoreGroup := router.Group("tatuadores")
	{
		tatuadoreGroup.GET(
			"",
			controllers.ListagemTatuador(client),
		)
		tatuadoreGroup.GET(
			":id",
			controllers.GetByIdTatuador(client),
		)
		tatuadoreGroup.GET(
			"estudios/:estudio_id",
			controllers.TatuadoresEstudio(client),
		)
		tatuadoreGroup.POST(
			"",
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
