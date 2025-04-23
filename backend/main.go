package main

import (
    "fmt"
    "log"
    "os"
    "github.com/parasmotwani/student-course-tracker/backend/models"
	"github.com/parasmotwani/student-course-tracker/backend/routes"

    "gorm.io/driver/postgres"
    "gorm.io/gorm"
    "github.com/joho/godotenv"
)

var DB *gorm.DB

func initDatabase() {
    err := godotenv.Load()
    if err != nil {
        log.Fatal("Error loading .env file")
    }

    dsn := fmt.Sprintf(
        "host=%s user=%s password=%s dbname=%s port=%s sslmode=require",
        os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"),
        os.Getenv("DB_NAME"), os.Getenv("DB_PORT"),
    )

    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatal("Failed to connect to database: ", err)
    }

    db.AutoMigrate(&models.Student{}, &models.Course{}, &models.Enrollment{})
    DB = db
}

func main() {
    initDatabase()
	controllers.SetDB(DB)
    r := routes.SetupRoutes(DB)
    r.Run(":8080")
}
