package main

import (
	"fmt"
	"os"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	supabase "github.com/nedpals/supabase-go"
	"tatovering/src/routes"
	"tatovering/src/middlewares"
	"tatovering/configs"
)

func main() {
	if err := godotenv.Load(); err != nil {
		fmt.Printf("Erro ao carregar arquivo .env: %v\n", err)
		os.Exit(1)
	}

	supabaseURL := os.Getenv("SUPABASE_URL")
	supabaseKey := os.Getenv("SUPABASE_KEY")

	// Create a Gin router
	router := gin.Default()
	router.Use(configs.Cors())
	
	//Initialize a single supabase client instead of one for each query received
	client := supabase.CreateClient(supabaseURL, supabaseKey)

	routes.SetupUsuariosRoutes(router, client)
	routes.SetupTatuadoresRoutes(router, client)
	routes.SetupTatuagemRoutes(router, client)
	routes.SetupEstudiosRoutes(router, client)

	private := router.Group("/private", middlewares.JwtTokenCheck(client))
	fmt.Println(88, private)

	port := 8080 // Change to the desired port
	router.Run(fmt.Sprintf(":%d", port))
}
