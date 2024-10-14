# Mykonsul

## Description

This project is a backend service written in Golang that leverages MongoDB for data storage and Redis as a caching layer. It uses Docker Compose to orchestrate the containers for the service and it's dependencies.

## Technologies Used

- **Golang**: Backend service
- **MongoDB**: Database storage
- **Redis**: Caching layer
- **Docker Compose**: Service orchestration

## Prerequisites

Make sure you have the following installed on your machine:

- [Docker](https://www.docker.com/get-started)
- [Docker Compose](https://docs.docker.com/compose/)
- [Golang](https://golang.org/dl/) (if running locally without Docker)

## Installation and Setup

1. **Clone the repository**:

   ```bash
   git clone git@github.com:fadhilsurya/mykonsul-mongo.git
   cd mykonsul-mongo
   ```

2. **Install dependencies**:

   Use the following Go commands:

   ```bash
   go mod tidy
   go mod download
   ```

3. **Set up environment variables**:

   Create a `.env` file in the root directory:

```bash
APP_NAME=
ENV=
PORT=
DB_PORT=
DB_ADDRESS=
GIN_MODE=
JWT_SECRET=
REDIS_ADDRESS=
REDIS_PORT=
```

4. **Build and run the service locally (optional)**:

   ```bash
   go build -o main .
   ./main
   ```

## Using Docker Compose

1. **Make sure Docker is running** on your system.

2. **Start the services** using Docker Compose:

   ```bash
   docker-compose up --build
   ```

3. **Verify that everything is running**:

   - Golang service: [http://localhost:8080](http://localhost:8080)
   - MongoDB: Accessible inside the network at `mongo:27017`
   - Redis: Accessible inside the network at `redis:6379`

4. **Stop the services**:

   ```bash
   docker-compose down
   ```

## Docker Compose Configuration

Here is the `docker-compose.yml` used in this project:

```yaml
services:
  mongo:
    image: mongo:latest
    container_name: mongo
    ports:
      - "27017:27017"
    volumes:
      - mongo-data:/data/db

  redis:
    image: redis:alpine
    container_name: redis
    ports:
      - "6379:6379"

  golang-service:
    build: .
    container_name: golang-service
    command: ["./main"]
    ports:
      - "8080:8080"
    depends_on:
      - mongo
      - redis
    env_file:
      - .env

volumes:
  mongo-data:
```

## Usage

- **Running the service**:
  After running `docker-compose up --build`, the service will be available at [http://localhost:8080](http://localhost:8080).

- **Access the MongoDB container**:

  ```bash
  docker exec -it mongo mongo
  ```

- **Access the Redis container**:
  ```bash
  docker exec -it redis redis-cli
  ```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

## Example Commands

- **Rebuild the service**:

  ```bash
  docker-compose up --build --force-recreate
  ```

- **Check running containers**:

  ```bash
  docker ps
  ```

- **Clean up Docker volumes**:
  ```bash
  docker volume prune
  ```
