# student-course-tracker
A Student Course Tracker web app using Go, PostgreSQL, and Next.js.

## Features

- Student registration & login
- Course listing & enrollment
- Course rating system
- Token-based authentication using JWT

## Getting Started

### Prerequisites

Make sure you have the following installed:

- [Go](https://golang.org/dl/) (v1.20+ recommended)
- [PostgreSQL](https://www.postgresql.org/download/)
- [Git](https://git-scm.com/)


---

###  Setup Instructions

#### 1. Clone the Repository

In you terminal write:<br>

```bash
git clone https://github.com/yourusername/student-course-tracker.git 
cd student-course-tracker/backend
```

#### 2. Create a PostgreSQL database (host for free on aiven)

Link - https://aiven.io/

#### 3. Generate a JWT Secret key by running generate_jwt.py

In terminal navigate to the directory where generate_jwt.py is saved and write the following command<br>

```bash
python generate_jwt.py
```

OR<br>
Simply run the script

#### 4. Create a .env file

Create a .env file inside the backend/ folder with the following variables:<br>

```
DB_HOST=your_aiven_host
DB_PORT=your_aiven_port
DB_USER=your_aiven_user
DB_NAME=your_aiven_database_name
DB_PASSWORD=your_aiven_database_password
JWT_SECRET=your_jwt_secret_generated_by_script
```

#### 5. Initialize and install Go dependencies

In the terminal of your working directory write the following:-<br>

Initializing Go module:<br>

```bash
go mod init student-course-tracker
```

Installing dependencies:<br>

```bash
go get github.com/gin-gonic/gin
go get gorm.io/gorm
go get gorm.io/driver/postgres
go get github.com/joho/godotenv
go get github.com/golang-jwt/jwt/v5
go get golang.org/x/crypto/bcrypt
```

#### 6. Run the backend server

```bash
go run main.go
```




