package middlewares

import (
	"context"
	"net/http"

	"github.com/aungmyozaw92/go-restapi-sta/utils"
	"github.com/gin-gonic/gin"
)


type authString string

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.Request.Header.Get("Authorization")

		if auth == "" {
			c.Next()
			return
		}

		bearer := "Bearer "
		auth = auth[len(bearer):]

		// Check if the token is blacklisted
		// _, exists, err := config.GetRedisValue(auth)
		// if err == nil && exists {
		// 	c.JSON(http.StatusUnauthorized, gin.H{"error": "token is invalidated"})
		// 	c.Abort()
		// 	return
		// }

		// Validate the token
		validate, err := utils.JwtValidate(auth)
		
		if err != nil || !validate.Valid {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "unauthorized",
				"error" : err.Error(),
			})
			c.Abort()
			return
		}

		customClaim, _ := validate.Claims.(*utils.JwtCustomClaim)

		ctx := context.WithValue(c.Request.Context(), authString("auth"), customClaim)
		ctx = context.WithValue(ctx, utils.ContextKeyUserId, customClaim.ID)
		ctx = context.WithValue(ctx, utils.ContextKeyToken, auth)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

func CtxValue(ctx context.Context) *utils.JwtCustomClaim {
	raw, _ := ctx.Value(authString("auth")).(*utils.JwtCustomClaim)
	return raw
}