package models

import (
    "gorm.io/gorm"
    "time"
)

type Student struct {
    ID       uint   `gorm:"primaryKey"`
    Name     string
    Email    string `gorm:"unique"`
    Password string
    Enrollments []Enrollment
}

type Course struct {
    ID        uint    `gorm:"primaryKey"`
    Name      string
    CreatedAt time.Time
    Rating    float64
    Enrollments []Enrollment
}

type Enrollment struct {
    ID        uint `gorm:"primaryKey"`
    StudentID uint
    CourseID  uint
    Rating    *int
}
