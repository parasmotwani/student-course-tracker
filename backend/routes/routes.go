package routes

import (
    "github.com/gin-gonic/gin"
    "student-course-tracker/controllers" 
	"gorm.io/gorm"
	"student-course-tracker/middleware"
)

var DB *gorm.DB

func SetDB(database *gorm.DB) {
    DB = database
}
func SetupRouter() *gin.Engine {
    r := gin.Default()

    // public
    r.POST("/register", controllers.Register)
    r.POST("/login",    controllers.Login)

    // protected by JWT
    auth := r.Group("/").Use(middleware.JWTAuth())
    {
        auth.GET( "/enrollments", controllers.GetEnrollments)
        auth.POST("/enroll",      controllers.EnrollStudent)
        auth.POST("/rate-course", controllers.RateCourse)
		auth.GET( "/courses",     controllers.GetCourses)
    }

    return r
}