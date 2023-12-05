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
			"",
			middlewares.JwtTokenCheck(client), 
			controllers.VerAgendamentosUsuario(client),
		)
		estudioGroup.POST(
			"",
			middlewares.JwtTokenCheck(client),
			controllers.EfetuarAgendamentoUsuario(client),
		)
	}
}
