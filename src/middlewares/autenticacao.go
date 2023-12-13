package middlewares

import (
	"errors"
	"net/http"
	"strings"
	"github.com/gin-gonic/gin"
	supabase "github.com/nedpals/supabase-go"
)

func ExtractBearerToken(header string) (string, error) {
	if header == "" {
		return "", errors.New("Missing authorization header")
	}

	jwtToken := strings.Split(header, " ")
	if len(jwtToken) != 2 {
		return "", errors.New("Incorrectly formatted authorization header")
	}

	return jwtToken[1], nil
}

func JwtTokenCheck(client *supabase.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		jwtToken, err := ExtractBearerToken(c.GetHeader("Authorization"))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		client.DB.AddHeader("Authorization", "Bearer "+jwtToken)
		c.Next()
	}
}
