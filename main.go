package main

import (
	"fmt"
	"net/http"
	"context"
	"github.com/gin-gonic/gin"
	supabase "github.com/nedpals/supabase-go"
)

func main() {
	// Initialize your Supabase client
	supabaseURL := "https://pfzlboeaonsookzcnniv.supabase.co"
	supabaseKey := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzdXBhYmFzZSIsInJlZiI6InBmemxib2Vhb25zb29remNubml2Iiwicm9sZSI6ImFub24iLCJpYXQiOjE2OTIwNjI3NTYsImV4cCI6MjAwNzYzODc1Nn0.KuEEX9EBIQmLTA02iPtqqNIewDmXITDxnIfD4qEqTN8"

	// Create a Gin router
	router := gin.Default()

// Route for user sign-up
router.POST("/signup", func(c *gin.Context) {
	var requestBody struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	client := supabase.CreateClient(supabaseURL, supabaseKey)

	// Create user credentials
	credentials := supabase.UserCredentials{
		Email:    requestBody.Email,
		Password: requestBody.Password,
	}
	ctx := context.Background()
	// Sign up the user with Supabase
	user, err := client.Auth.SignUp(ctx, credentials)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"user": user})
})

// Route for user sign-in
router.POST("/signin", func(c *gin.Context) {
	var requestBody struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	client := supabase.CreateClient(supabaseURL, supabaseKey)

	// Create user credentials
	credentials := supabase.UserCredentials{
		Email:    requestBody.Email,
		Password: requestBody.Password,
	}

	ctx := context.Background()
	// Sign up the user with Supabase
	user, err := client.Auth.SignIn(ctx, credentials)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
})

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
