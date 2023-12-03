package routes

import (
	"tatovering/src/controllers"
	"tatovering/src/middlewares"

	"github.com/gin-gonic/gin"
	supabase "github.com/nedpals/supabase-go"
)

func SetupTatuagemRoutes(router *gin.Engine, client *supabase.Client) {
	tatuagemGroup := router.Group("tatuagens")
	{
		tatuagemGroup.GET(
			"",
			controllers.ListagemTatuagem(client),
		)
		tatuagemGroup.GET(":id",
			controllers.GetByIdTatuagem(client),
		)
		tatuagemGroup.GET("tatuador/:tatuador_id",
			controllers.GetArtByTatuadorId(client),
		)
		tatuagemGroup.GET("favoritos/:id",
			middlewares.JwtTokenCheck(client),
			controllers.GetFavoritos(client),
		)
		tatuagemGroup.POST(
			"favoritar",
			controllers.FavoritarTatuagem(client),
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
			"favoritos",
			middlewares.JwtTokenCheck(client),
			controllers.DeleteFavorito(client),
		)
		tatuagemGroup.DELETE(
			"/:id",
			middlewares.JwtTokenCheck(client),
			controllers.DeletarTatuagem(client),
		)
	}
}
