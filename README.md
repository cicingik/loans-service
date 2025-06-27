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


## Assumption
- Loan Proposed by borrower
    - Borrower can propose loan after previous loan disburse. 
- Rounding value for loan, funding amount, roi is 2 digits
-


## TODO
- jwt with role
- route loan -> create and update
- route loan/funding -> create

