package main

import (
	"fmt"
	"log"
	"os"
	"student-course-tracker/controllers"
	"student-course-tracker/models"
	"student-course-tracker/routes"

	"github.com/gin-contrib/cors"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func initDatabase() {
	if err := godotenv.Load(); err != nil {
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

	DB = db
	DB.AutoMigrate(&models.Student{}, &models.Course{}, &models.Enrollment{})
}

func seedCourses(db *gorm.DB) {
	sampleCourses := []models.Course{
		{Name: "Introduction to Golang", Rating: 70},
		{Name: "Machine Learning Basics", Rating: 70},
		{Name: "Frontend with React", Rating: 70},
		{Name: "Database Systems", Rating: 70},
		{Name: "Cloud Computing Fundamentals", Rating: 70},
	}

	for _, c := range sampleCourses {
		var exists models.Course
		if err := db.Where("name = ?", c.Name).First(&exists).Error; err == gorm.ErrRecordNotFound {
			db.Create(&c)
		}
	}
}

func main() {
	initDatabase()
	controllers.SetDB(DB)
	seedCourses(DB)

	r := routes.SetupRouter()

	// CORS — allow your Next front-end on port 3000
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	log.Println("Listening on :8080")
	r.Run(":8080")
}
