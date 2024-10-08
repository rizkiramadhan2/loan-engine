
# Loan Engine

## Stack
1. Go 1.21 (inside docker) with Gin HTTP Framework
2. Postgresql 11 (inside docker)


## Pre-requisite
1. Docker Compose v2.20+

## Run Server
```sh
make dev
```

The server will run on port 4040.

Test: [http://localhost:4040/ping](http://localhost:4040/ping)

## API Design
Postman collection: [Postman collection](https://github.com/rizkiramadhan2/loan-engine/blob/main/docs/loan_engine.postman_collection.json)
### GET /loans
Get list of loans

**Response:**
```json
[
    {
        "id": 4,
        "borrower_id": 1,
        "principal_amount": 1500,
        "rate": 0.2,
        "roi": 0.1,
        "status": 1,
        "status_str": "proposed",
        "agreement_letter_url": "http://example-of-agreement-letter.com"
    },
    {
        "id": 2,
        "borrower_id": 1,
        "principal_amount": 1500,
        "rate": 0.2,
        "roi": 0.1,
        "status": 3,
        "status_str": "invested",
        "agreement_letter_url": "http://example-of-agreement-letter.com",
        "picture_proof_url": "http://example-of-proof",
        "approver_id": 2,
        "approval_date": "2024-06-25T11:16:12.533823+07:00"
    },
    {
        "id": 1,
        "borrower_id": 1,
        "principal_amount": 1500,
        "rate": 0.2,
        "roi": 0.1,
        "status": 4,
        "status_str": "disbursed",
        "agreement_letter_url": "http://example-of-agreement-letter.com",
        "picture_proof_url": "http://example-of-proof",
        "approver_id": 2,
        "approval_date": "2024-06-25T12:01:38.413757+07:00"
    }
]
```

### GET /loans/:id/detail
Get detail of loans

**Response:**
```json
{
    "id": 2,
    "borrower_id": 1,
    "principal_amount": 1500,
    "rate": 0.2,
    "roi": 0.1,
    "status": 3,
    "status_str": "invested",
    "agreement_letter_url": "http://example-of-agreement-letter.com",
    "picture_proof_url": "http://example-of-proof",
    "approver_id": 2,
    "approval_date": "2024-06-25T11:16:12.533823+07:00",
    "investors": [
        {
            "id": 1,
            "loan_id": 2,
            "investor_id": 3,
            "amount": 1500
        }
    ],
    "disbursement": null
}
```

### POST /loans
Request a loan and give default status of 'proposed'

**Request:**
```json
{
    "borrower_id": 1,
    "principal_amount": 1500,
    "roi": 0.1,
    "rate": 0.2
}
```

**Response:**
```json
{
    "id": 2,
    "borrower_id": 1,
    "principal_amount": 1500,
    "rate": 0.2,
    "roi": 0.1,
    "status": 1,
    "status_str": "proposed",
    "agreement_letter_url": "http://example-of-agreement-letter.com"
}
```

### PATCH /loans/:id/approve
Approve a loan, will update status to 'approved'

**Request:**
```json
{
    "picture_proof_url": "http://example-of-proof",
    "approver_id": 2
}
```

**Response:**
```json
{
    "loan_id": 1,
    "status": "approved"
}
```

### POST /loans/:id/invest
Invest a loan, will update status to 'invested' if the total amount of invested is equal to principal amount

**Request:**
```json
{
    "investor_id": 3,
    "amount": 1500
}
```

**Response:**
```json
{
    "loan_id": 2,
    "status": "invested",
    "total_of_invested": 1500
}
```

### POST /loans/:id/disburse
Disburse a loan, will update status to 'disbursed'

**Request:**
```json
{
    "signed_agreement_url": "http://example-of-agreement-url",
    "disburser_employee_id": 1
}
```

**Response:**
```json
{
    "loan_id": 4,
    "status": "disbursed"
}
```

## DB Design
There are 3 tables that hold data of loan

1. **loan**: loan request will be stored here, along with the approval (picture_proof_url, approver_id, approval_date)
```sql
CREATE TABLE IF NOT EXISTS public.loan (
    id SERIAL PRIMARY KEY,
    borrower_id INT NOT NULL,
    principal_amount FLOAT NOT NULL,
    rate FLOAT NOT NULL,
    roi FLOAT NOT NULL,
    status INT NOT NULL DEFAULT 1,
    agreement_letter_url TEXT,
    picture_proof_url TEXT,
    approver_id INT,
    approval_date TIMESTAMPTZ
);
```

2. **loan_investment**: this table holds the information about the investment
```sql
CREATE TABLE IF NOT EXISTS public.loan_investment (
    ID SERIAL PRIMARY KEY,
    loan_id INTEGER NOT NULL,
    investor_id VARCHAR NOT NULL,
    amount FLOAT NOT NULL
);
```

3. **loan_disbursement**: this table holds the information about the disbursement
```sql
CREATE TABLE IF NOT EXISTS public.loan_disbursement (
    ID SERIAL PRIMARY KEY,
    loan_id INTEGER NOT NULL,
    signed_agreement_url TEXT NOT NULL,
    disburser_employee_id VARCHAR NOT NULL,
    disbursement_date TIMESTAMPTZ NOT NULL
);
```


## Function Implementation
### Handler
location: app/api/http/handler/loan.go
To validate incoming request before forward to process
- **CreateLoan():** Create a new loan application.
- **ApproveLoan():** Approve a loan application.
- **InvestLoan():** Invest in a loan opportunity.
- **DisburseLoan():** Transfer approved loan amounts.
- **GetDetail():** Retrieve detailed information about a loan.
- **GetList():** Retrieve a list of loans.

### UseCase
**Location: `internal/usecase/loan/loan.go`**
Handles the logic and rules of loan operations.

- **CreateLoan():** Logic to create a new loan application.
- **ApproveLoan():** Logic to approve a loan application.
- **InvestLoan():** Logic to invest in a loan.
- **DisburseLoan():** Logic to disburse approved loan amounts.
- **GetDetail():** Logic to retrieve detailed information about a loan.
- **GetList():** Logic to retrieve a list of loans.

### Repository
Handles data retrieval and modification from the database.

#### Fetch: Getter Data
**Location: `internal/repository/loan/fetch.go`**

- **GetByID():** Retrieve loan details by ID.
- **GetList():** Retrieve a list of loans.
- **GetInvestByID():** Retrieve investment details by loan ID.
- **GetDisburseByID():** Retrieve disbursement details by loan ID.

#### Mutation: Setter Data
**Location: `internal/repository/loan/mutation.go`**

- **Create():** Create a new loan record.
- **Approve():** Approve a loan application in the database.
- **UpdateStatus():** Update the status of a loan.
- **Invest():** Record an investment in a loan.
- **Disburse():** Record the disbursement of loan funds.




## Code Layout

```
├── Makefile
├── app
│   ├── api
│   │   └── http
│   │       ├── handler
│   │       │   ├── error.go
│   │       │   ├── init.go
│   │       │   └── loan.go
│   │       └── server.go
│   ├── interface.go
├── cmd
│   ├── http-api
│   │   └── main.go
├── config
│   ├── config.go
│   └── type.go
├── docs
│   └── loan_engine.postman_collection.json
├── files
│   └── etc
│       └── simple-app
│           └── config.yaml
├── internal
│   ├── model
│   │   ├── loan.go
│   ├── pkg
│   │   ├── agreementLetter
│   │   │   └── agreementLetter.go
│   ├── repository
│   │   └── loan
│   │       ├── fetch.go
│   │       ├── init.go
│   │       └── mutation.go
│   └── usecase
│       └── loan
│           └── loan.go
```

### Root Directory
- **Makefile**: Contains a set of directives used by the make build automation tool to compile and manage the project.

### app Directory
- **api**: Likely contains the API layer of the application.
- **http**: Specifically for HTTP-related components.
  - **handler**: Contains HTTP request handlers.
    - **error.go**: Manages error responses and error handling logic.
    - **init.go**: Initialization logic for the handlers.
    - **loan.go**: Handler logic for loan-related HTTP endpoints.
  - **server.go**: Sets up and runs the HTTP server, likely configuring routes and middleware.
- **interface.go**: Could define interfaces for the application, perhaps for dependency injection or defining contracts between layers.

### cmd Directory
- **http-api**: Typically contains the entry point for the HTTP API service.
  - **main.go**: The main file that starts the HTTP API server.

### config Directory
- **config.go**: Logic for loading and managing configuration settings.
- **type.go**: Defines configuration-related types and structures.

### docs Directory
- **loan_engine.postman_collection.json**: loan engine postman collection.

### files Directory
- **etc**: Likely for configuration and environment-specific files.
  - **simple-app**: Contains specific configurations for the application.
    - **config.yaml**: YAML file containing configuration settings.

### internal Directory
- **model**: Contains the data models or domain entities.
  - **loan.go**: Defines the loan data model.
- **pkg**: Holds package-specific code that is used internally.
  - **agreementLetter**: Contains code related to agreement letters.
    - **agreementLetter.go**: Logic for managing agreement letters.
- **repository**: Data access layer, handling interactions with the database or data sources.
  - **loan**: Repository logic specific to loans.
    - **fetch.go**: Logic to fetch loan data from the data source.
    - **init.go**: Initialization logic for the loan repository.
    - **mutation.go**: Logic for modifying loan data (e.g., create, update, delete).
- **usecase**: Business logic layer, containing the core functionality and rules.
  - **loan**: Use cases related to loans.
    - **loan.go**: Business logic for loan operations.