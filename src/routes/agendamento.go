package routes

import (
	"tatovering/src/controllers"
	"tatovering/src/middlewares"

	"github.com/gin-gonic/gin"
	supabase "github.com/nedpals/supabase-go"
)

func SetupAgendamentosRoutes(router *gin.Engine, client *supabase.Client) {
	estudioGroup := router.Group("agendamentos")
	{
		estudioGroup.GET(
			"tatuador/:tatuador_id",
			controllers.ObterDisponibilidadeTatuador(client),
		)
		estudioGroup.GET(
			"",
			middlewares.JwtTokenCheck(client),
			controllers.ObterAgendamentosUsuario(client),
		)
		estudioGroup.GET(
			"tatuador",
			middlewares.JwtTokenCheck(client),
			controllers.ObterAgendamentosTatuador(client),
		)
		estudioGroup.POST(
			"",
			middlewares.JwtTokenCheck(client),
			controllers.EfetuarAgendamentoUsuario(client),
		)
	}
}
