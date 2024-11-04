Recruitment Management System
Overview
The Recruitment Management System is a backend application designed to manage recruitment processes. Built with Go and the Gin framework, the application leverages PostgreSQL for data storage and AWS services (EC2 and RDS) for deployment. This project is designed for handling recruitment processes efficiently and securely.

Tech Stack
Backend: Go (Gin Framework)
Database: PostgreSQL
Deployment: AWS EC2, AWS RDS
Prerequisites
Go 1.23 or higher
PostgreSQL (locally for development)
Docker (optional for containerization)
Environment Configuration
Database Configuration (For Local Development)
Update your environment file (prod.env for deployment or .env for local) with the following database configuration:

Copy code
DB_HOST=127.0.0.1
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=0000
DB_NAME=rms
DB_SSLMODE=require
JWT_SECRET=your_jwt_secret_key
Local Setup Instructions
Clone the Repository:

bash
Copy code
git clone https://github.com/yourusername/recruitment-management-system.git
cd recruitment-management-system
Install Dependencies:

bash
Copy code
go mod download
Set Up Database: Ensure PostgreSQL is running locally with a database named rms. Use the provided database configuration above for connecting to the database.

Run the Server Locally: In main.go, update the server run command to use localhost for local development.

go
Copy code
server.Run("localhost:8080")
Start the server:

bash
Copy code
go run main.go
Accessing the API: Once running, the API can be accessed at http://localhost:8080.

Deployment Setup (AWS EC2 and RDS)
Database Setup on AWS RDS:

Create a PostgreSQL instance on RDS.
Update prod.env with the RDS database endpoint.
Deploying on EC2:

Launch an EC2 instance and configure security groups to allow inbound traffic on port 8080.(0.0.0.0:8080 TCP)
SSH into the EC2 instance and pull the Docker image (or set up a Go environment to run the application directly).
Running with Docker (Optional): If you prefer Docker, ensure that Docker is installed on your EC2 instance and build/pull the Docker image.

bash
Copy code
docker build -t recruitment-management-system:latest .
docker run -d -p 8080:8080 --env-file prod.env recruitment-management-system:latest
