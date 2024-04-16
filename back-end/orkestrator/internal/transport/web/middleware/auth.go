package middleware

import (
	"log"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/Conty111/SuperCalculator/back-end/orkestrator/internal/clierrs"
	"github.com/Conty111/SuperCalculator/back-end/orkestrator/internal/transport/web/helpers"
	"github.com/gin-gonic/gin"
)

var notRequireAuthenticationEndpoints = []string{
	"/api/v1/users/login",
	"/api/v1/users/create",
	"/api/v1/status",
}

func TokenAuthMiddleware(publicKeyPath string) gin.HandlerFunc {
	jwtVerifier, err := helpers.NewJWTVerifier(publicKeyPath)
	if err != nil {
		log.Fatal(err)
	}

	return func(c *gin.Context) {
		setAccessHeaders(c)

		url := c.Request.URL.String()

		if !requiresAuth(url) {
			c.Next()
			return
		}

		token := getTokenFromHeader(&c.Request.Header)
		if token == "" {
			respondWithError(c, http.StatusUnauthorized, clierrs.ErrAuthTokenWasNotProvided)
			return
		}

		claims, err := jwtVerifier.Verify(token)
		if err != nil {
			respondWithError(c, http.StatusUnauthorized, clierrs.ErrInvalidAuthToken)
			return
		}

		if time.Now().After(claims.Expires) {
			respondWithError(c, http.StatusUnauthorized, clierrs.ErrTokenExpired)
			return
		}

		c.Set("callerID", claims.UserID)

		c.Next()
	}
}

func setAccessHeaders(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Credentials", "true")
	c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Auth, Authorization, accept, origin, Cache-Control, X-Requested-With")
	c.Header("Access-Control-Allow-Methods", "POST, PATCH, OPTIONS, GET, PUT, DELETE")

	if c.Request.Method == "OPTIONS" {
		c.AbortWithStatus(http.StatusNoContent)
		return
	}
}

func getTokenFromHeader(header *http.Header) string {
	tokenData := strings.Split(header.Get("Authorization"), " ")
	return tokenData[len(tokenData)-1]
}

func respondWithError(c *gin.Context, code int, err error) {
	c.AbortWithStatusJSON(code, gin.H{"error": err.Error()})
}

func requiresAuth(url string) bool {
	if strings.Contains(url, "swagger") {
		return false
	}
	return !slices.Contains(notRequireAuthenticationEndpoints, url)
}
