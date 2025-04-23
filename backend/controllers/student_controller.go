// Global DB instance (will be passed from main)
var DB *gorm.DB

func SetDB(database *gorm.DB) {
    DB = database
}

// Register a new student
func Register(c *gin.Context) {
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
    // Assuming student ID is in the JWT token
    tokenString := c.GetHeader("Authorization")
    claims, err := utils.VerifyJWT(tokenString)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or missing token"})
        return
    }

    var student models.Student
    if err := DB.Preload("Courses").First(&student, claims.Subject).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get enrollments"})
        return
    }

    c.JSON(http.StatusOK, student.Courses)
}

// Rate a course
func RateCourse(c *gin.Context) {
    var rateRequest struct {
        CourseID uint   `json:"course_id"`
        Rating   uint   `json:"rating"`
    }

    if err := c.ShouldBindJSON(&rateRequest); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
        return
    }

    // Assuming student ID is in the JWT token
    tokenString := c.GetHeader("Authorization")
    claims, err := utils.VerifyJWT(tokenString)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or missing token"})
        return
    }

    var enrollment models.Enrollment
    if err := DB.Where("student_id = ? AND course_id = ?", claims.Subject, rateRequest.CourseID).First(&enrollment).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Course enrollment not found"})
        return
    }

    // Update the course rating
    course := models.Course{ID: rateRequest.CourseID}
    if err := DB.Model(&course).UpdateColumn("rating", gorm.Expr("rating + ? / ?", rateRequest.Rating, 1)).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to rate course"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Course rated successfully"})
}