package models

import (
    //"gorm.io/gorm"
    "time"
)

type Student struct {
	ID        uint           `gorm:"primaryKey" json:"id"` //
    Name     string `json:"username"` // maps "username" to "Name"
    Email    string `json:"email" gorm:"unique"`
    Password string `json:"password"`
    Enrollments []Enrollment `gorm:"foreignKey:StudentID"`
}

type Course struct {
	ID        uint           `gorm:"primaryKey" json:"id"` //
    Name      string    `json:"name"` // maps "name" to "Name"
    CreatedAt time.Time `json:"created_at"`
    Rating    float64   `json:"rating"`
    Enrollments []Enrollment `gorm:"foreignKey:StudentID"`
}
/* Helper function
func (c *Course) SetID(id uint) {
	c.Model.ID = id
} */

type Enrollment struct {
	ID        uint           `gorm:"primaryKey" json:"id"` //
    StudentID uint     `json:"student_id"` // maps "student_id" to StudentID
    CourseID  uint     `json:"course_id"`  // maps "course_id" to CourseID
    Rating    *int     `json:"rating"`
    Course    Course   `gorm:"foreignKey:CourseID" json:"course"` // Nested Course
}
