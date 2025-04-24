package controllers

import (
    "net/http"
    "student-course-tracker/models"

    "github.com/gin-gonic/gin"
)

// GetCourses returns all courses
func GetCourses(c *gin.Context) {
    var courses []models.Course
    if err := DB.Find(&courses).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch courses"})
        return
    }
    c.JSON(http.StatusOK, courses)
}
