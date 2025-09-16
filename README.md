# VIT 2026 Capstone Internship Hiring Task

This project is a full-stack web application with a Go backend, a React frontend, and a PostgreSQL database.

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes.

### Prerequisites

You need to have Docker and Docker Compose installed on your machine.

*   [Docker](https://docs.docker.com/get-docker/)
*   [Docker Compose](https://docs.docker.com/compose/install/)

### Installation

1.  Clone the repo
    ```sh
    git clone https://github.com/joel2607/vit-2026-capstone-internship-hiring-task-joel2607.git
    ```
2.  Navigate to the project directory
    ```sh
    cd vit-2026-capstone-internship-hiring-task-joel2607
    ```

## Environment Variables

Create a `.env` file in the root of the project and add the following environment variables:

```
POSTGRES_USER=user
POSTGRES_PASSWORD=password
POSTGRES_DB=file_vault
```

## Running the application

To run the application, run the following command from the root of the project:

```sh
docker-compose up
```

This will start the following services:

*   **backend**: Go API running on `http://localhost:8080`
*   **frontend**: React app running on `http://localhost:3000`
*   **db**: PostgreSQL database running on `localhost:5432`



[![Review Assignment Due Date](https://classroom.github.com/assets/deadline-readme-button-22041afd0340ce965d47ae6ef1cefeee28c7c493a6346c4f15d667ab976d596c.svg)](https://classroom.github.com/a/2xw7QaEj)