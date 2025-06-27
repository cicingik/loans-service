CREATE SCHEMA IF NOT EXISTS users;

CREATE TABLE IF NOT EXISTS users.users
(
    id         bigserial primary key,
    role_id    bigserial not null,
    name       varchar(255),
    created_at timestamp default now(),
    updated_at timestamp default now(),
    deleted_at timestamp
);


CREATE TABLE IF NOT EXISTS users.login_data
(
    id           bigserial primary key,
    user_id      bigserial    not null,
    user_name    varchar(255) not null,
    password     varchar(255) not null,
    created_at   timestamp default now(),
    updated_at   timestamp default now(),
    deleted_at   timestamp
);

CREATE TABLE IF NOT EXISTS users.roles
(
    id          bigserial primary key,
    description varchar(255) not null,
    created_at  timestamp default now(),
    updated_at  timestamp default now(),
    deleted_at  timestamp
);


CREATE SCHEMA IF NOT EXISTS loans;


CREATE TABLE IF NOT EXISTS loans.loans
(
    id                bigserial primary key,
    approved_at       timestamp,
    disburse_at       timestamp,
    letter_aggrement  varchar(255),
    visit_document    varchar(255),
    status            varchar(255) not null,
    disburse_by       varchar(255),
    approve_by        varchar(255),
    borrower_id       bigserial    not null,
    rate              float        not null,
    principal         float        not null,
    funding_remaining float        not null,
    roi               float        not null,
    created_at        timestamp default now(),
    updated_at        timestamp default now(),
    deleted_at        timestamp
);

CREATE TABLE IF NOT EXISTS loans.loan_fundings
(
    id               bigserial primary key,
    funding_at       timestamp,
    letter_aggrement varchar(255) not null,
    lender_id        bigserial    not null,
    loan_id          bigserial    not null,
    funding_amount   float        not null,
    created_at       timestamp default now(),
    updated_at       timestamp default now(),
    deleted_at       timestamp
);


DROP TYPE IF EXISTS loans.loan_status;
CREATE TYPE loans.loan_status AS ENUM ('proposed', 'approved', 'invested', 'disbursed');
