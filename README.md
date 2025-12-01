# Flip Full Stack Technical Test

## Description
This application allows users upload a bank statement file, view insights, and
inspect transaction issues.

## Setup Instructions

Prerequisites:
- Docker
- Docker Compose

To run the application:

1.  Navigate to the root directory of the project.
2.  Run the following command:
    ```bash
    docker-compose up
    ```
3.  Access the application at [http://localhost](http://localhost).

## Architecture Decisions

### Tech Stack & Libraries

-   **Backend**: **Go**
    -   **Gin**: Gin provides essential features like routing, middleware support, and JSON validation without unnecessary overhead.
    -   **uber-go/zap**: Blazing fast, structured, leveled logging in Go.
    -   **vektra/mockery**: A mock code autogenerator for Go interfaces, essential for testing.
    -   **spf13/viper**: A complete configuration solution for Go applications.
    -   **samber/oops**: A library for better error handling with stack traces and context.

-   **Frontend**: **React** (with **Next.js**)
    -   **TypeScript**: TypeScript provides static typing, which improves code quality, maintainability, and developer experience by catching errors early.
    -   **CSS Modules**: Scopes CSS to the component, preventing style conflicts.
    -   **classnames**: A simple JavaScript utility for conditionally joining classNames together. 

## Endpoints

The backend exposes the following RESTful API endpoints:

1.  **`POST /upload`**
    -   **Description**: Uploads a CSV file containing transaction data.
    -   **Payload**: `multipart/form-data` with a `file` field.

2.  **`GET /balance`**
    -   **Description**: Retrieves the total balance calculated from successful transactions (Credits - Debits).

3.  **`GET /issues`**
    -   **Description**: Retrieves a paginated list of unsuccessful transactions (Failed or Pending).
    -   **Query Parameters**:
        -   `page`: Page number (default: 1)
        -   `page_size`: Number of items per page (default: 10)
        -   `sorts`: Sort field and direction (e.g., `+timestamp`, `-amount`)

## Test Files

The root directory contains several CSV files that can be used for testing the application:

-   `test_bank_statement.csv`: A standard bank statement with various transaction types.
-   `test_invalid_bank_statement.csv`: Contains invalid data to test error handling.
-   `test_negative_bank_statement.csv`: Contains transactions that result in a negative balance.
