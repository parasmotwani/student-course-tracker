package utils

import (
    "golang.org/x/crypto/bcrypt"
    "github.com/golang-jwt/jwt/v5"
    "time"
    "os"
)

func HashPassword(password string) (string, error) {
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    return string(hashedPassword), err
}

func CheckPasswordHash(password, hash string) error {
    return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

func GenerateJWT(studentID uint) (string, error) {
    secret := os.Getenv("JWT_SECRET")

    claims := &jwt.RegisteredClaims{
        Issuer:    "student-course-tracker",
        Subject:   string(studentID), // Store student ID in token
        ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)), // token expires in 24 hours
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(secret))
}

func VerifyJWT(tokenString string) (*jwt.RegisteredClaims, error) {
    secret := os.Getenv("JWT_SECRET")

    token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
        return []byte(secret), nil
    })

    if err != nil {
        return nil, err
    }

    if claims, ok := token.Claims.(*jwt.RegisteredClaims); ok && token.Valid {
        return claims, nil
    } else {
        return nil, err
    }
}
