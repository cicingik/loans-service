# Loans Service

## Database Designs

```mermaid
erDiagram

roles {
    int id PK
    string desctiption
}

users {
    int id PK
    int role_id FK
    string name
}

user_login_data {
    int id PK
    int user_id FK
    string login_name
    string password_hash
}

loans {
    int id PK
    int borrower_id FK
    float principal
    float rate
    float rate_percentage
    float funding_remaining
    float roi
    string status
    datetime disburse_at
    datetime approved_at
    string visited_by
    string disburse_by
    string visit_document
    string aggrement_letter
}

loan_fundings {
    int id PK
    int lender_id FK
    int loan_id FK
    float funding_amount 
    float roi
    datetime funding_at
    string aggrement_letter

}

users }o--|| roles : is
users ||--|| user_login_data : have
users ||--o{ loans : have
users ||--o{ loan_fundings : have
loans }o--o{ loan_fundings: funded

```

## API

### How to Run
```
mv .env.example .envrc
direnv allow .

make serve
```
### List API

```txt
POST http://0.0.0.0:2727/v1/user/login
{
  "user_name": "lender",
  "password": "12345"
}


GET http://0.0.0.0:2727/v1/assessment

PUT http://0.0.0.0:2727/v1/assessment/{loan_id}/{status}
{
  "document": "doc_visit",
  "employee_id": "12345",
  "execute_at": "2025-06-27T07:20:50.52Z"  
}


POST http://0.0.0.0:2727/v1/funding/{loan_id}
{
  "funding_amount": 1000000
}
```

## Cron

This cron will change status loan to `invested` when loan can not be fund anymore and will generate letter agreement for lender when loan already `invested`

```txt
mv .env.example .envrc
direnv allow .

make cron
```


