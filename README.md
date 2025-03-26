# Overview

A system that emulates a distributed manufacturing pipeline simulation system running on Kubernetes. It has a User Interface with authentication and authorization and provides real-time control and status monitoring. The system simulates the creation of a pipeline with multiple stages that can run sequentially or in parallel based on the user's choice. The backend is built using Go (Golang) and deployed on Kubernetes, utilizing REST and gRPC APIs for communication.

## Backend Architecture

The backend is structured into modular components for scalability and maintainability. It consists of:

### Backend Tech Stack:
- **Golang** (Primary Backend Language)
- **GORM** (ORM for Database Management)
- **Supabase** (Managed Backend Platform using PostgreSQL)
- **Gin Framework** (REST API Development)
- **gRPC** (Google Remote Procedure Call for High-Performance API Communication)

### Core Interfaces:
- **Pipeline Interface**: Orchestrates multiple stages in either sequential or parallel execution.
- **Stage Interface**: Defines the structure of each pipeline stage, handling execution, error management, and rollback.
- **Database Interface**: Defines the database functions.

## Database & Supabase Details

The system uses **Supabase** as the backend database, which internally runs on **PostgreSQL**. Supabase also provides authentication (user sign-up, login, and session management). The database is initialized in Golang using **GORM**, which migrates tables from the `model.go` file (structured as structs). The migrations include **users**, **pipeline_executions**, and **execution_logs** tables. The **Supabase Go SDK** is initialized to interact with authentication and database services. User authentication and metadata updates are handled via **Supabase's Auth**.

<p align="center">
  <img src="https://github.com/user-attachments/assets/2dd201e9-b89d-4d90-8f9a-0201af011f37" alt="Database & Supabase Architecture">
</p>


## Frontend Architecture

The frontend is built using **React.js** and serves as the user interface for interaction. It provides a seamless experience for users to authenticate, create pipelines, start pipelines, and monitor the status of stages and pipelines.

### Frontend Tech Stack:
- **React.js** (Frontend framework)
- **Axios** (API Communication)
- **Material UI** (Styling)
- **React Router** (Navigation & Routing)
- **SSE** (Server-Sent Events for real-time updates)

## Key Features

- **User Authentication**: Users can register, log in, and manage sessions using **Supabase authentication**. Authentication is handled via **JWT tokens**, securely stored in the browser.
- **Pipeline Management**: Users can create pipelines, select the number of stages and execution type, and get real-time updates on **pipeline status, logs, and execution progress**.

---

# Running the Project using Docker and Kubernetes

## Prerequisites
1. Install **Docker Desktop**: [Download](https://www.docker.com/get-started/) and verify with:
   ```sh
   docker --version
   ```
2. Enable **Kubernetes** from Docker Desktop and verify:
   ```sh
   kubectl version --client
   kubectl get nodes
   ```

## Deploying Docker Images on Kubernetes
1. Pull the Docker Images:
   ```sh
   docker pull hsri/frontend:v1.0.0
   docker pull hsri/rest-api:v1.0.0
   docker pull hsri/grpc-server:v1.0.0
   ```
2. Download deployment files from `distributed-manufacturing-pipeline/deploy/kubernetes/deployment`. It contains:
   - `deployment-frontend.yaml`
   - `deployment-rest.yaml`
   - `deployment-grpc.yaml`
3. Create `config-db.yaml` and `secret-db.yaml`:

### `config-db.yaml`
```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: supabase-config
data:
  SUPABASE_DB: "XXXXXXXXXXX"
```

### `secret-db.yaml`
```yaml
apiVersion: v1
kind: Secret
metadata:
  name: supabase-secrets
type: Opaque
data:
  SUPABASE_URL: XXXXXXXXXXXXXXXXX
  SUPABASE_API_KEY: XXXXXXXXXXXXXXX
```
4. Deploy all configurations:
   ```sh
   kubectl apply -f config-db.yaml
   kubectl apply -f secret-db.yaml
   kubectl apply -f deployment-frontend.yaml
   kubectl apply -f deployment-rest.yaml
   kubectl apply -f deployment-grpc.yaml
   ```
5. Verify Kubernetes pods and services:
   ```sh
   kubectl get pods
   kubectl get svc
   ```
6. Access the frontend via `http://localhost:30080` and check logs with:
   ```sh
   kubectl logs -l app=rest-api -n default -f
   ```

---

# Running the Project Locally

## Prerequisites

### Install Required Tools:
- **Golang (v1.21+)**: [Download](https://go.dev/dl/) and verify:
  ```sh
  go version
  ```
  Install dependencies:
  ```sh
  go mod tidy
  ```
- **Node.js (v18+)**: [Download](https://nodejs.org/en) or install via CLI:
  ```sh
  curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.40.2/install.sh | bash
  \ "$HOME/.nvm/nvm.sh"
  nvm install 22
  ```
  Verify installation:
  ```sh
  node -v  # Should print "v22.14.0"
  npm -v   # Should print "10.9.2"
  ```
- Install frontend dependencies:
  ```sh
  cd web/frontend
  npm install
  ```

## Steps to Start the Project Locally

1. Clone the repository:
   ```sh
   git clone https://github.com/hsri-pf9/distributed-manufacturing-pipeline-simulation-system.git
   cd distributed-manufacturing-pipeline
   ```
2. Create a `.env` file in the root directory:
   ```sh
   SUPABASE_DB = XXXXXXXXXXXXXXX
   SUPABASE_URL = XXXXXXXXXXXXXX
   SUPABASE_API_KEY= XXXXXXXXXXXXXXX
   ```
3. Modify `db.go` and `supabase_client.go` to use environment variables (see full instructions above).
4. Initialize the backend:
   ```sh
   cd cmd/api-server
   go run main.go
   ```
5. Start the frontend:
   ```sh
   cd web/frontend
   npm start
   ```
6. Access the UI at `http://localhost:3000`

---

# Using `democtl` CLI with gRPC Server on Kubernetes

1. Clone the repository:
   ```sh
   git clone https://github.com/hsri-pf9/distributed-manufacturing-pipeline-simulation-system.git
   cd distributed-manufacturing-pipeline
   ```
2. Build `democtl`:
   ```sh
   go build -o democtl ./cmd/democtl/main.go
   sudo mv democtl /usr/local/bin/
   ```
3. Set environment variables:
   ```sh
   export DEMOCTL_GRPC_URL=localhost:50051
   ```
4. Forward gRPC service port:
   ```sh
   kubectl port-forward svc/grpc-server 50051:50051
   ```
5. Use `democtl` commands:
   ```sh
   democtl register --email --password
   democtl login --email --password
   democtl pipeline create --user --stages --parallel
   democtl pipeline start --pipeline-id "XXXXXX" --user-id "XXXXXX" --input "test_input" --parallel
   democtl pipeline status --pipeline-id "XXXXX" --parallel
   democtl pipeline cancel --pipeline-id "XXXXX" --user-id "XXXXXXX" --parallel
   ```
---
# REST Endpoints Overview

The backend REST API provides endpoints for user authentication, pipeline management, and real-time updates using Gin and SSE (Server-Sent Events). Below is a structured table of the available API endpoints.

## User Authentication Endpoints

| Method | Endpoint    | Description          | Auth Required | Request Body | Response |
|--------|------------|----------------------|---------------|--------------|----------|
| POST   | /register  | Register a new user  | ❌ No         | `{ "email": "user@example.com", "password": "password" }` | `{ "message": "Check your email for verification" }` |
| POST   | /login     | User login           | ❌ No         | `{ "email": "user@example.com", "password": "password" }` | `{ "token": "jwt-token", "user_id": "uuid" }` |
| GET    | /user/:id  | Fetch user profile   | ✅ Yes        | N/A          | `{ "user_id": "uuid", "email": "user@example.com", "role": "worker" }` |
| PUT    | /user/:id  | Update user profile  | ✅ Yes        | `{ "name": "New Name", "role": "admin" }` | `{ "message": "User updated successfully" }` |

## Pipeline Management Endpoints

| Method | Endpoint              | Description                  | Auth Required | Request Body | Response |
|--------|------------------------|------------------------------|---------------|--------------|----------|
| GET    | /pipelines             | Get all pipelines for a user | ✅ Yes        | N/A          | `[ { "pipeline_id": "uuid", "status": "Running" } ]` |
| GET    | /pipelines/:id/stages  | Get pipeline stages          | ✅ Yes        | N/A          | `{ "stages": [ ... ] }` |
| POST   | /createpipelines       | Create a new pipeline        | ✅ Yes        | `{ "user_id": "uuid", "stage_count": 5, "is_parallel": true }` | `{ "pipeline_id": "uuid" }` |
| POST   | /pipelines/:id/start   | Start a pipeline execution   | ✅ Yes        | `{ "user_id": "uuid" }` | `{ "status": "Running" }` |
| GET    | /pipelines/:id/status  | Get pipeline execution status | ✅ Yes        | N/A          | `{ "pipeline_id": "uuid", "status": "Running" }` |
| POST   | /pipelines/:id/cancel  | Cancel a pipeline execution  | ✅ Yes        | `{ "user_id": "uuid" }` | `{ "status": "Cancelled" }` |

## Real-Time Updates & SSE

| Method | Endpoint              | Description                  | Auth Required | Response |
|--------|------------------------|------------------------------|---------------|----------|
| GET    | /pipelines/:id/stream  | Subscribe to real-time updates | ✅ Yes | Event Stream (SSE) |

## API Flow Diagrams

### User management

<p align="center">
  <img src="https://github.com/user-attachments/assets/c32d6075-78d2-4bd8-bd38-3b3a934cce39" alt="Database & Supabase Architecture">
</p>


### Pipeline Management

<p align="center">
  <img src="https://github.com/user-attachments/assets/0f860e34-6d4d-457d-ac57-5b6379bfd9a7" alt="Database & Supabase Architecture">
</p>

## Deployment Architecture


<p align="center">
  <img src="https://github.com/user-attachments/assets/9a4b7b8a-3e3a-41eb-abe9-06ed5494a55c" alt="Deployment Architecture">
</p>

The system consists of three main components, each running in a separate pod inside the Kubernetes cluster:

### Frontend (React + Nginx)
- Serves the user interface.
- Exposes port 80 via a NodePort service (`frontend-service`).
- Uses Nginx as a reverse proxy to forward API requests (`/api/`) to the REST API.

### REST API (Golang)
- Handles authentication, user management, and pipeline execution.
- Exposes port 8080 via a NodePort service (`rest-api-service`).
- Calls the gRPC Server for specific operations.
- Uses Kubernetes ConfigMaps & Secrets for environment variables.

### gRPC Server (Golang)
- Provides high-performance operations for internal use.
- Exposes port 50051 via a ClusterIP service (`grpc-server`).
- Only accessible within the cluster, meaning users cannot directly call it.
- Uses Kubernetes ConfigMaps & Secrets for configuration.

## Step-by-Step Execution Flow

### 1. User Accesses the Frontend (NodePort 30080)
- The React app is served by Nginx.
- When a user performs an action (like logging in or starting a pipeline), Nginx forwards API requests to the REST API via /api/*.

### 2. REST API Processes Requests (NodePort 30081)
- The REST API handles authentication, user management, and pipeline execution.
- It interacts directly with Supabase for database operations.

### 3. CLI Calls gRPC Server for High-Performance Tasks
- The democtl CLI tool interacts with the gRPC Server via port 50051.
- The gRPC Server does not expose a NodePort, meaning it is not accessible externally.
- This design ensures better security and performance since gRPC is optimized for fast internal communication.

## Configuration & Secrets Management
- Kubernetes ConfigMaps and Secrets store database connection details and API keys.
- Both REST API and gRPC Server pods use these environment variables for secure communication.

## Deployment Architecture Components

### Frontend Deployment (`frontend`)
- **Pod:** `frontend-pod`
- **Image:** `hsri/frontend:v1.0.0`
- **Service Type:** NodePort
- **Port Exposed:** 80 (external), mapped to `30080`
- **Reverse Proxy:** Nginx forwards `/api/` requests to the REST API.

### REST API Deployment (`rest-api`)
- **Pod:** `rest-api-pod`
- **Image:** `hsri/rest-api:v1.0.0`
- **Service Type:** NodePort
- **Port Exposed:** 8080, mapped to `30081`
- **Communicates with:** gRPC Server for backend processing.

### gRPC Server Deployment (`grpc-server`)
- **Pod:** `grpc-server-pod`
- **Image:** `hsri/grpc-server:v1.0.0`
- **Service Type:** ClusterIP (only accessible within the cluster)
- **Port Exposed:** 50051
- **Communicates with:** REST API internally.

### Configuration & Secrets
- **ConfigMaps (`supabase-config`)**
  - Stores database connection details (`SUPABASE_DB`).
- **Secrets (`supabase-secrets`)**
  - Stores API keys (`SUPABASE_URL`, `SUPABASE_API_KEY`).
  - Used by both REST API and gRPC Server.

## Deployment Summary

| Component       | Pod Name         | Service Type | Exposed Port | Purpose |
|----------------|-----------------|--------------|--------------|---------|
| **Frontend**   | frontend-pod     | NodePort     | 80 → 30080   | Serves React UI via Nginx |
| **REST API**   | rest-api-pod     | NodePort     | 8080 → 30081 | Handles API requests and business logic |
| **gRPC Server** | grpc-server-pod  | ClusterIP    | 50051        | Processes backend operations for REST API |
| **ConfigMaps & Secrets** | supabase-config & supabase-secrets | - | - | Stores database credentials and API keys |

---

# Thank You for visiting my Repo
