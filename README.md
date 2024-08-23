# NOTO ![Status](https://img.shields.io/badge/status-in%20development-yellow)
Noto ノート (nōto) is an app for storing your notes privately. You can host it on your localhost or your own server, ensuring that your data is entirely yours and not owned by anyone else.

# Table of Contents
- [NOTO ](#noto-)
- [Table of Contents](#table-of-contents)
  - [Development](#development)
    - [Running Local Server](#running-local-server)
    - [Generate Swagger Documentation](#generate-swagger-documentation)
  - [Deployment](#deployment)

## Development
### Running Local Server
1. **Clone the Repository**
   ```sh
   git clone https://github.com/Shiyinq/noto.git
   cd noto
   ```

2. **Install Go Modules**
   ```sh
   go mod tidy
   ```

3. **Create .env File**
   ```sh
   cp .env.example .env
   ```

4. **Install Air for Live Reloading**

   If you don't have `air` installed on your machine, install it first:
   ```sh
   go install github.com/air-verse/air@latest
   ```

5. **Run the Development Server**
   ```sh
   air
   ```

6. **Server**

    http://localhost:8080

### Generate Swagger Documentation
1. **Install Swagger for API Documentation**

   If you don't have `swag` installed on your machine, install it first:
   ```sh
   go install github.com/swaggo/swag/cmd/swag@latest
   ```

2. **Generate or Update Documentation**
    ```sh
    swag init -g ./cmd/server/main.go --parseDependency --parseInternal --output docs/swagger
    ```
    Or you can use the `swag.sh` script:

    For the first time, before running the script, execute:
    ```
    chmod +x docs.sh
    ```
    Then, run:
    ```
    ./swag.sh
    ```

3. **Swagger**

    http://localhost:8080/swagger/index.html

## Deployment

Before you begin, ensure you have [Docker](https://docs.docker.com/engine/install/) installed.

**1. Clone the Repository**
```sh
git clone https://github.com/Shiyinq/noto.git
cd noto
```

**2. Create Environment Files**

For the backend:
```sh
cp .env.example .env
```

For the frontend:
```sh
cd cmd/client
cp .env.example .env
cd ../../
```

Open each `.env` file you have created and update the values as needed.

**3. Build and Run the Docker Containers**
```sh
docker compose up --build -d
```
Wait a few minutes for the setup to complete. You can then access:
- Frontend at http://localhost:5000
- Backend at http://localhost:8080
