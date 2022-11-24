package utils

import (
	"context"
	"github.com/gin-gonic/gin"
	"google.golang.org/api/idtoken"
)

func ValidateIdToken(c *gin.Context) (map[string]interface{}, error) {
	idToken := extractToken(c)
	payload, err := idtoken.Validate(context.Background(), idToken, "")
	if err != nil {
		return nil, err
	}
	return payload.Claims, nil
}
