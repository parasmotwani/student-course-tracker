package controllers

import (
    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
    "student-course-tracker/models"
    "net/http"
    "student-course-tracker/utils" 
	"log"
)

// Global DB instance (will be passed from main)
var DB *gorm.DB

func SetDB(database *gorm.DB) {
	log.Println("✔️ controllers.SetDB was called")
    DB = database
}

// Register a new student
func Register(c *gin.Context) {

	if DB == nil {
		log.Println("ERROR: Database is not initialized")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database not initialized"})
		return
	}
	
    var student models.Student
    if err := c.ShouldBindJSON(&student); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
        return
    }

    // Hash password before saving (use bcrypt for password hashing)
    hashedPassword, err := utils.HashPassword(student.Password)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
        return
    }
    student.Password = hashedPassword

	// Ensure DB is properly initialized before calling Create
	if DB == nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Database not initialized"})
        return
    }

    // Save student in database
    if result := DB.Create(&student); result.Error != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Registration successful"})
}

// Login a student and return JWT
func Login(c *gin.Context) {
    var loginRequest struct {
        Email    string `json:"email"`
        Password string `json:"password"`
    }

    if err := c.ShouldBindJSON(&loginRequest); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
        return
    }

    var student models.Student
    if err := DB.Where("email = ?", loginRequest.Email).First(&student).Error; err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
        return
    }

    // Compare the hashed password
    if err := utils.CheckPasswordHash(loginRequest.Password, student.Password); err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
        return
    }

    // Create JWT token
    token, err := utils.GenerateJWT(student.ID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate JWT"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"token": token})
}

// Get all courses a student is enrolled in
func GetEnrollments(c *gin.Context) {
	tokenString, err := utils.ExtractTokenFromHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	claims, err := utils.VerifyJWT(tokenString)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
		return
	}

	var student models.Student
	if err := DB.Preload("Enrollments.Course").First(&student, claims.Subject).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get enrollments"})
		return
	}

	var courses []models.Course
	for _, enrollment := range student.Enrollments {
		courses = append(courses, enrollment.Course)
	}

	c.JSON(http.StatusOK, courses)
}

// Rate a course
func RateCourse(c *gin.Context) {
	var rateRequest struct {
		CourseID uint `json:"course_id"`
		Rating   uint `json:"rating"`
	}

	if err := c.ShouldBindJSON(&rateRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	tokenString, err := utils.ExtractTokenFromHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	claims, err := utils.VerifyJWT(tokenString)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
		return
	}

	var enrollment models.Enrollment
	if err := DB.Where("student_id = ? AND course_id = ?", claims.Subject, rateRequest.CourseID).First(&enrollment).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Course enrollment not found"})
		return
	}

	course := models.Course{ID: rateRequest.CourseID}
	//course.SetID(rateRequest.CourseID)

	if err := DB.Model(&course).UpdateColumn("rating", gorm.Expr("rating + ? / ?", rateRequest.Rating, 1)).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to rate course"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Course rated successfully"})
}

// Enroll a student in a course
func EnrollStudent(c *gin.Context) {
    var enrollmentRequest struct {
        StudentID uint `json:"student_id"`
        CourseID  uint `json:"course_id"`
    }

    // Bind the request body to the struct
    if err := c.ShouldBindJSON(&enrollmentRequest); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
        return
    }

    // Check if the course exists before enrolling
    var course models.Course
    if err := DB.First(&course, enrollmentRequest.CourseID).Error; err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Course not found"})
        return
    }

    // Create the enrollment record
    enrollment := models.Enrollment{
        StudentID: enrollmentRequest.StudentID,
        CourseID:  enrollmentRequest.CourseID,
    }

    // Insert the enrollment record into the database
    if err := DB.Create(&enrollment).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create enrollment"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Student enrolled successfully"})
}