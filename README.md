# API Health Checker 
[![Go](https://img.shields.io/badge/Go-1.25.5-blue)](https://golang.org)
[![License:MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT) 
[![Docker](https://img.shields.io/badge/Docker-%2300D7DF.svg?&logo=docker&logoColor=white)](https://www.docker.com/)


Lightweight service for checking health and availability of HTTP APIs. It is intended to monitor internal and external services and can be used both locally and as part of infrastructure.

## Features
- Periodic health checks of configured services
- Status change notifications via Telegram
- REST API for managing services
- Persistent storage of service configurations
- Detailed logging


## Configuration
Create a `.env` file in the project root with the following variables:

```bash
cd path_to_project
cp .env.example .env
```

### Environment Variables

- `TG_BOT_TOKEN`: Telegram bot token for sending notifications
- `TG_CHAT_ID`: Telegram chat ID to send notifications to
- `SERVICES_DURATION`: Interval between service checks (default: 60s)
- `SERVICES_FILE`: Path to the JSON file for storing service configurations (default: ./data/services.json)
- `HTTP_ADDR`: HTTP server address (default: :8081)
- `CHECKER_ADDR`: Health checker address (default: http://localhost:8081/services) 

## API Endpoints

### GET /services
List all monitored services

### POST /services
Add a new service to monitoring

**Request body:**
```json
{
  "Name": "Service name",
  "URL": "http://example.com"
}
```

### DELETE /services
Remove a service from monitoring

**Request body:**
```json
{
  "Name": "Service name"
}
```

## Service Structure

```go
type Service struct {
    Name     string    // Service name
    URL      string    // Service URL to check
    IsUp     bool      // Current status
    LastDown time.Time // Timestamp of last downtime
}
```

## Architecture

The application follows a modular architecture with the following components:

- `app/`: Main application code
- `cli/`: Command-line interface
- `internal/common/`: Shared data structures
- `internal/logs/`: Logging functionality
- `internal/notifier/`: Telegram notification service
- `internal/services/`: Service monitoring logic
- `internal/storage/`: Data persistence

## CLI Usage

A command-line interface is available for managing monitored services:

```bash
# Build the CLI
go build -o checker-cli cli/main.go

# List all services
./checker-cli list

# Add a new service
./checker-cli add -n "Service Name" -url "http://example.com"

# Delete a service
./checker-cli delete -n "Service Name"

# Show help
./checker-cli help
```

### CLI Commands

- `list`: Display all monitored services in formatted JSON
- `add`: Add a new service to monitoring
- `delete`: Remove a service from monitoring
- `help`: Show available commands and usage

The CLI reads configuration from the same `.env` file as the main application and uses the `CHECKER_ADDR` environment variable to determine the API endpoint (default: http://localhost:8081/services).

## Getting Started
```bash
# 1. Clone the repository
mkdir healthchecker && cd healthchecker
git clone https://github.com/mksmin/api-health-checker.git .

# 2. Create .env file
cp .env.example .env
# edit .env (TG_BOT_TOKEN and TG_CHAT_ID is required)

# 3. Create data directory
mkdir -p data logs 

# 4. Run locally (optional)
go run app/main.go

```


## Docker
The application can be run in a container. This is the recommended approach.

```bash
# Build the image
docker build -t ghcr.io/yourusername/api-health-checker:latest .

# Run the container
docker run -d \
  --name healthchecker \
  -v $(pwd)/data:/app/data \
  --env-file .env \
  -p 8081:8081 \
  ghcr.io/yourusername/api-health-checker:latest

```

## Logs

Application logs are written to:
- Console output
- `logs/healthchecker.log` file

## License

[MIT](LICENSE)