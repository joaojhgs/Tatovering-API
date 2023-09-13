package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	supabase "github.com/nedpals/supabase-go"
)

func main() {
	// Initialize your Supabase client
	supabaseURL := "https://pfzlboeaonsookzcnniv.supabase.co"
	supabaseKey := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzdXBhYmFzZSIsInJlZiI6InBmemxib2Vhb25zb29remNubml2Iiwicm9sZSI6ImFub24iLCJpYXQiOjE2OTIwNjI3NTYsImV4cCI6MjAwNzYzODc1Nn0.KuEEX9EBIQmLTA02iPtqqNIewDmXITDxnIfD4qEqTN8"

	// Create a Gin router
	router := gin.Default()

	// Define CRUD routes for "usuarios"
	router.POST("/usuarios", func(c *gin.Context) {
		// Create a new usuario
    row := Usuario{
      Nome: "Gabriel Medrado",
    }
    var results []Usuario
    client := supabase.CreateClient(supabaseURL, supabaseKey)
    err := client.DB.From("usuarios").Insert(row).Execute(&results)
    
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, results)
	})

	// router.GET("/usuarios/:id", func(c *gin.Context) {
	// 	// Get a usuario by ID
	// 	tatuadorID := c.Param("id")
	// 	resp, err := client.
	// 		Table("usuarios").
	// 		Select("*").
	// 		Equals("id", tatuadorID).
	// 		Single()
	// 	if err != nil {
	// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	// 		return
	// 	}

	// 	c.JSON(http.StatusOK, resp)
	// })

	// router.PUT("/usuarios/:id", func(c *gin.Context) {
	// 	// Update a usuario by ID
	// 	tatuadorID := c.Param("id")
	// 	var usuario Usuario // Define the Usuario struct as needed
	// 	if err := c.BindJSON(&usuario); err != nil {
	// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	// 		return
	// 	}

	// 	resp, err := client.
	// 		Table("usuarios").
	// 		Update(tatuadorToMap(usuario)).
	// 		Equals("id", tatuadorID).
	// 		Single()
	// 	if err != nil {
	// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	// 		return
	// 	}

	// 	c.JSON(http.StatusOK, resp)
	// })

	// router.DELETE("/usuarios/:id", func(c *gin.Context) {
	// 	// Delete a usuario by ID
	// 	tatuadorID := c.Param("id")
	// 	resp, err := client.
	// 		Table("usuarios").
	// 		Delete().
	// 		Equals("id", tatuadorID).
	// 		Single()
	// 	if err != nil {
	// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	// 		return
	// 	}

	// 	c.JSON(http.StatusOK, resp)
	// })

	// Start the Gin server
	port := 8080 // Change to the desired port
	router.Run(fmt.Sprintf(":%d", port))
}

// Define the Usuario struct to match your database structure
type Usuario struct {
  Nome    string `json:"nome"`
}

// Helper function to convert Usuario struct to map for Supabase
func tatuadorToMap(usuario Usuario) map[string]interface{} {
	return map[string]interface{}{
		"nome":      usuario.Nome,
	}
}
