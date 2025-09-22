# BalkanID Full Stack Engineering Intern — Capstone Hiring Task

## Live Demo

The application is deployed and accessible at: **[https://filevault.ajayjoel.space](https://filevault.ajayjoel.space)**

## Overview

This project is a secure file vault system that supports efficient storage, powerful search, and controlled file sharing. It is a production-grade file vault application with both backend and frontend components, designed to demonstrate skills in API design, backend services in Go, relational data modeling with PostgreSQL, and modern frontend development with React.js and TypeScript.

## Core Features

-   **File Deduplication**: Detects duplicate uploads using SHA-256 content hashing to save storage space.
-   **File Uploads**: Supports single and multiple file uploads with a drag-and-drop interface.
-   **File Management & Sharing**:
    -   List and view files with detailed metadata.
    -   Organize files into folders.
    -   Share files and folders publicly via a link or privately with specific users.
    -   Public file statistics to track download counts.
    -   Strict delete rules to prevent accidental data loss.
-   **Search & Filtering**:
    -   Search files by filename.
    -   Filter files by MIME type, size, date range, and tags.
-   **Rate Limiting & Quotas**:
    -   Per-user API rate limits (2 calls per second).
    -   Per-user storage quotas (10 MB).
-   **Storage Statistics**:
    -   Display total, original, and saved storage usage.

## Bonus Features Implemented

-   **Real-time Updates**: Live download counts using GraphQL subscriptions.
-   **Folder Organization**: Users can create, update, and delete folders to organize their files.
-   **Role-Based Access Control (RBAC)**: Differentiates between regular users and admins. Admins can view all files using the search bar.

## Tech Stack

-   **Backend**: Go (Golang)
-   **API Layer**: GraphQL
-   **Database**: PostgreSQL
-   **Cache**: Redis
-   **Frontend**: Next.js (React) with TypeScript
-   **Containerization**: Docker Compose
-   **Reverse Proxy**: Nginx

## System Architecture

The application is built on a microservices-based architecture:

-   **Go Backend**: A robust GraphQL API that handles all business logic, including user authentication, file management, and storage operations.
-   **Next.js Frontend**: A modern, user-friendly interface that consumes the GraphQL API to provide a seamless user experience.
-   **PostgreSQL Database**: The primary data store for user information, file metadata, and folder structures.
-   **Redis**: Used for caching and to manage rate limiting.
-   **Nginx**: Acts as a reverse proxy, directing traffic to the appropriate service and handling SSL termination.

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes.

### Prerequisites

You need to have Docker and Docker Compose installed on your machine.

-   [Docker](https://docs.docker.com/get-docker/)
-   [Docker Compose](https://docs.docker.com/compose/install/)

### Installation

1.  Clone the repo
    ```sh
    git clone https://github.com/joel2607/vit-2026-capstone-internship-hiring-task-joel2607.git
    ```
2.  Navigate to the project directory
    ```sh
    cd vit-2026-capstone-internship-hiring-task-joel2607
    ```

### Environment Variables

Create a `.env` file in the root of the project and add the following environment variables:

```
# PostgreSQL Configuration
POSTGRES_USER=user
POSTGRES_PASSWORD=password
POSTGRES_DB=file_vault

# JWT Secret
JWT_AUTH_SECRET=your-secret-key

# Backend URL for the frontend
BACKEND_URL=http://localhost:8080

# Rate Limit
RATELIMIT_LIMIT=100

# GraphQL Endpoints for the frontend
NEXT_PUBLIC_GRAPHQL_ENDPOINT=http://localhost:8080/graphql
NEXT_PUBLIC_GRAPHQL_WS_ENDPOINT=ws://localhost:8080/graphql
```

### Running the application

To run the application, run the following command from the root of the project:

```sh
docker-compose up --build
```

This will start the following services:

-   **backend**: Go API running on `http://localhost:8080`
-   **frontend**: React app running on `http://localhost:3000`
-   **db**: PostgreSQL database running on `localhost:5433`
-   **redis**: Redis server running on `localhost:6379`



### API Documentation

The backend exposes a GraphQL API. The schema is defined in `backend/graphQL/schema.graphql` and provides a comprehensive overview of the available queries, mutations, and subscriptions.

## Testing with Postman

A Postman workspace has been provided to facilitate API testing.

**Workspace Link**: [https://web.postman.co/9471df3d-c48b-4b14-897c-e6cd77995772](https://web.postman.co/9471df3d-c48b-4b14-897c-e6cd77995772)

### Setup

1.  **Fork the Collection**: Fork the collection into your own Postman workspace.
2.  **Set Environment Variables**: The collection uses a variable `{{base_url}}`. It is recommended to set this to `http://localhost:8080` for local testing or `filevault.ajayjoel.space` for testing production.

### Example Usage

-   **Register a new user**: Use the `register` mutation.
-   **Login**: Use the `login` mutation to obtain a JWT token. This token will be automatically used in subsequent requests.
-   **Upload a file**: Use the `uploadFiles` mutation. You will need to attach a file to the request.
-   **Search for files**: Use the `searchFiles` query. Admins can search for any file, while regular users can only search for their own files.


## Github Classroom Link
[![Review Assignment Due Date](https://classroom.github.com/assets/deadline-readme-button-22041afd0340ce965d47ae6ef1cefeee28c7c493a6346c4f15d667ab976d596c.svg)](https://classroom.github.com/a/2xw7QaEj)