## Student and Mentor Dashboard Setup

### Stack

- Go
- HTMX
- PostgreSQL
- Gorm

### Prerequisites

1. **Go**: Ensure you have Go installed. You can download it from [golang.org](https://golang.org/dl/).
2. **PostgreSQL**: Make sure PostgreSQL is installed and running.

### Setup Options

#### 1. Clone the Repository

1. Open your terminal and clone the repository:

   ```sh
   git clone https://github.com/U-to-E/dashboard
   cd dashboard
   ```

#### 2. Install Dependencies

1. Install the necessary Go packages:

   ```sh
   go mod tidy
   ```

#### 3. Create Environment File

1. Create a `.env` file in the root of your project directory:

   ```sh
   vim .env
   ```

2. Add the following content to the `.env` file:

   ```env
   DB_HOST=localhost
   DB_PORT=5432
   DB_USER=postgres
   DB_PASSWORD=your_password
   DB_NAME=your_db_name
   SECRET=your_secret
   ADMIN_EMAIL=hi@utoe.com
   ADMIN_PASS=your_admin_password
   EMAIL_ADDR=email
   EMAIL_PASS=pass
   EMAIL_SMTP=smtp.gmail.com
   EMAIL_PORT=587
   ```

### Running the Application

#### Option 1: Binary Setup

1. Build the binary:

   ```sh
   go build -o dashboard
   ```

2. Run the binary:

   ```sh
   ./dashboard
   ```

#### Option 2: Build from Source

1. Run the application directly from the source:

   ```sh
   go run main.go
   ```

2. Open your browser and navigate to `http://localhost:3000`.

#### Option 3: Using Docker

1. Create a `Dockerfile` in the root of your project directory with the following content:

   ```dockerfile
   FROM golang:1.18-alpine

   WORKDIR /app

   COPY . .

   RUN go mod tidy
   RUN go build -o dashboard

   EXPOSE 3000

   CMD ["./dashboard"]
   ```

2. Build the Docker image:

   ```sh
   docker build -t dashboard-app .
   ```

3. Run the Docker container:

   ```sh
   docker run --env-file .env -p 3000:3000 dashboard-app
   ```

4. Open your browser and navigate to `http://localhost:3000`.
