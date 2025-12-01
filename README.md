# Flip Full Stack Technical Test

## Table of Contents
- [Description](#description)
- [Setup Instructions](#setup-instructions)
- [Architecture Decisions](#architecture-decisions)
- [Endpoints](#endpoints)
- [Test Files](#test-files)

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

### CI/CD (with Github Actions)
**Deployment**

1. Firstly, you will need to create a tag with the following format:
- FE: `{{ENV}}-fe-v{{VERSION}}` ex. **stg-fe-v0.01**, will deploy FE to staging environment with version 0.01
- BE: `{{ENV}}-be-v{{VERSION}}` ex. **stg-be-v0.01**, will deploy BE to staging environment with version 0.01

Valid Tags:
```
stg-fe-v*
stg-be-v*
prod-fe-V*
prod-be-v*
```
2. And then create a new github release by tagging a branch with the tag created from step 1 


**Pull Request**

When making a pull request to master branch, the pipeline will check for changes in backend and frontend folder. If there are changes in the folder it will run the respective jobs for that folder. BE job for backend folder and FE job for frontend folder.
- BE: Will run unit test and build verfication. If coverage is below 70% will return an error
- FE: Will run linter and build verfication

## Endpoints

The backend exposes the following RESTful API endpoints:

### 1. `POST /upload`
-   **Description**: Uploads a CSV file containing transaction data.
-   **Payload**: `multipart/form-data` with a `file` field.
-   **Responses**:
    -   **200 OK**:
        ```json
        {
            "data": {
                "transactions_uploaded": 28
            },
            "success": "success"
        }
        ```
    -   **400 Bad Request**:
        ```json
        {
            "metadata": {
                "path": "/upload",
                "code": "BE0002",
                "statusCode": 400,
                "status": "Bad Request",
                "message": "invalid CSV format: expected 6 columns, got 20 at line 1",
                "error": "POST /upload [400] Bad Request",
                "timestamp": "2025-12-01T14:56:19+07:00"
            },
            "success": "false"
        }
        ```
    -   **500 Internal Server Error**:
        ```json
        {
            "metadata": {
                "path": "/upload",
                "code": "BE0001",
                "statusCode": 500,
                "status": "Internal Server Error",
                "message": "failed to upload file",
                "error": "GET /upload [500] Internal Server Error",
                "timestamp": "2025-12-01T14:56:19+07:00"
            },
            "success": "false"
        }
        ```
### 2. `GET /balance`
-   **Description**: Retrieves the total balance calculated from successful transactions (Credits - Debits).
-   **Responses**:
    -   **200 OK**:
        ```json
        {
            "data": {
                "total_balance": 1000
            },
            "success": "success"
        }
        ```

    -   **500 Internal Server Error**:
        ```json
        {
            "metadata": {
                "path": "/balance",
                "code": "BE0001",
                "statusCode": 500,
                "status": "Internal Server Error",
                "message": "failed to calculate balance",
                "error": "GET /balance [500] Internal Server Error",
                "timestamp": "2025-12-01T14:56:19+07:00"
            },
            "success": "false"
        }
        ```

### 3. `GET /issues`
-   **Description**: Retrieves a paginated list of unsuccessful transactions (Failed or Pending).
-   **Query Parameters**:
    -   `page`: Page number (default: 1)
    -   `page_size`: Number of items per page (default: 10)
    -   `sorts`: Sort field and direction (e.g., `+timestamp`, `-amount`)
-   **Responses**:
    -   **200 OK**:
        ```json
        {
            "data": {
                "transactions": [
                    {
                        "transaction_date": "2025-12-01T10:00:00Z",
                        "name": "John Doe",
                        "type": "Credit",
                        "amount": 100,
                        "status": "FAILED",
                        "description": "Payment failed"
                    }
                ],
                "summary": {
                    "total_count": 23,
                    "pending_count": 11,
                    "failed_count": 12
                }
            },
            "pagination": {
                "current_page": 1,
                "current_elements": 1,
                "total_pages": 1,
                "total_elements": 1,
                "sort_by": "+timestamp"
            },
            "success": "success"
        }
        ```
    -   **400 Bad Request**:
        ```json
        {
            "metadata": {
                "path": "/issues",
                "code": "BE0002",
                "statusCode": 400,
                "status": "Bad Request",
                "message": "Invalid query param",
                "error": "POST /issues [400] Bad Request",
                "timestamp": "2025-12-01T14:56:19+07:00"
            },
            "success": "false"
        }
        ```
    -   **500 Internal Server Error**:
        ```json
        {
            "metadata": {
                "path": "/issues",
                "code": "BE0003",
                "statusCode": 500,
                "status": "Internal Server Error",
                "message": "failed to retrieve issues",
                "error": "GET /issues [500] Internal Server Error",
                "timestamp": "2025-12-01T14:56:19+07:00"
            },
            "success": "false"
        }
        ```

## Test Files

The root directory contains several CSV files that can be used for testing the application:

-   `test_bank_statement.csv`: A standard bank statement with various transaction types.
-   `test_invalid_bank_statement.csv`: Contains invalid data to test error handling.
-   `test_negative_bank_statement.csv`: Contains transactions that result in a negative balance.
