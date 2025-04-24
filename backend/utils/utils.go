package utils

import (
    "fmt"
    "golang.org/x/crypto/bcrypt"
    "github.com/golang-jwt/jwt/v5"
    "time"
    "os"
    "strconv"
	"errors"
	"strings"
	"github.com/gin-gonic/gin"
)

func HashPassword(password string) (string, error) {
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        return "", err
    }
    return string(hashedPassword), nil
}

func CheckPasswordHash(password, hash string) error {
    return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

func GenerateJWT(studentID uint) (string, error) {
    secret := os.Getenv("JWT_SECRET")

    claims := &jwt.RegisteredClaims{
        Issuer:    "student-course-tracker",
        Subject:   strconv.Itoa(int(studentID)), // Convert to string
        ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)), // token expires in 24 hours
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(secret))
}

func VerifyJWT(tokenString string) (*jwt.RegisteredClaims, error) {
    secret := os.Getenv("JWT_SECRET")

	fmt.Println("Received token:", tokenString)

    token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
        return []byte(secret), nil
    })

    if err != nil {
        // Handle expired token or any other jwt error
        if err.Error() == "token is expired" {
            return nil, fmt.Errorf("token is expired")
        }

        // For all other validation errors
        return nil, fmt.Errorf("failed to parse token: %v", err)
    }

    if claims, ok := token.Claims.(*jwt.RegisteredClaims); ok && token.Valid {
        return claims, nil
    } else {
        return nil, fmt.Errorf("invalid token")
    }
}

func ExtractTokenFromHeader(c *gin.Context) (string, error) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return "", errors.New("authorization header missing")
	}

	fields := strings.Fields(authHeader)
	if len(fields) != 2 || strings.ToLower(fields[0]) != "bearer" {
		return "", errors.New("invalid authorization header format")
	}

	return fields[1], nil
}