package middleware

import (
    "net/http"
    "strings"

    "github.com/gin-gonic/gin"
    "student-course-tracker/utils"
)

func JWTAuth() gin.HandlerFunc {
    return func(c *gin.Context) {
        auth := c.GetHeader("Authorization")
        if auth == "" {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error":"authorization header missing"})
            return
        }
        parts := strings.Fields(auth)
        if len(parts)!=2 || strings.ToLower(parts[0])!="bearer" {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error":"invalid auth format"})
            return
        }
        token := parts[1]
        _, err := utils.VerifyJWT(token)
        if err!=nil {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error":"invalid or expired token"})
            return
        }
        c.Next()
    }
}
