# Library Management System API

A RESTful API service for managing library operations including books, customers, borrowing journals, and late fee charges. Built with Go and Fiber framework.

## Description

This is a comprehensive library management system that handles:
- Book inventory and stock management
- Customer registration and management
- Book borrowing and return operations
- Automated late fee calculation and charging
- Media file management for book covers
- JWT-based authentication and authorization

## Tech Stack

- **Go** 1.25.6
- **Fiber** v2.52.10 - Web framework
- **PostgreSQL** - Database
- **JWT** (golang-jwt/jwt v5.3.1) - Authentication
- **Goqu** v9.19.0 - SQL query builder
- **Validator** v10.30.1 - Request validation
- **UUID** v1.6.0 - Unique identifier generation

## Project Structure

```
golang-perpustakaan/
├── domain/              # Domain entities and interfaces
│   ├── auth.go
│   ├── book.go
│   ├── book_stock.go
│   ├── charge.go
│   ├── customer.go
│   ├── journal.go
│   ├── media.go
│   ├── user.go
│   └── error.go
├── dto/                 # Data Transfer Objects
│   ├── auth_data.go
│   ├── book_data.go
│   ├── book_stock_data.go
│   ├── charge_data.go
│   ├── customer_data.go
│   ├── journal_data.go
│   ├── media_data.go
│   └── response.go
├── internal/
│   ├── api/            # HTTP handlers
│   │   ├── auth.go
│   │   ├── book.go
│   │   ├── book_stock.go
│   │   ├── customer.go
│   │   ├── journal.go
│   │   └── media.go
│   ├── config/         # Configuration management
│   │   ├── loader.go
│   │   └── model.go
│   ├── connection/     # Database connection
│   │   └── database.go
│   ├── repository/     # Data access layer
│   │   ├── book.go
│   │   ├── book_stock.go
│   │   ├── charge.go
│   │   ├── customer.go
│   │   ├── journal.go
│   │   ├── media.go
│   │   └── user.go
│   ├── service/        # Business logic layer
│   │   ├── auth.go
│   │   ├── book.go
│   │   ├── book_stock.go
│   │   ├── customer.go
│   │   ├── journal.go
│   │   └── media.go
│   └── util/           # Utilities
│       ├── claim.go
│       └── validation.go
├── storage/            # File storage directory
├── main.go             # Application entry point
└── go.mod
```

## Setup

### Prerequisites

- Go 1.25.6 or higher
- PostgreSQL
- Git

### Environment Variables

Create a `.env` file in the root directory:

```env
# Server Configuration
SERVER_HOST=localhost
SERVER_PORT=8080
SERVER_ASSET=http://localhost:8080/media

# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_NAME=perpustakaan
DB_USER=your_db_user
DB_PASS=your_db_password
DB_TZ=Asia/Jakarta
DB_SCHEMA=perpustakaan

# JWT Configuration
JWT_KEY=your_secret_jwt_key
JWT_EXP=3600

# Storage Configuration
STORAGE_BASE_PATH=./storage
```

## Installation

1. Clone the repository
```bash
git clone https://github.com/yusriltakeuchi/golang-perpustakaan.git
cd golang-perpustakaan
```

2. Install dependencies
```bash
go mod download
```

3. Create database and schema
```bash
createdb perpustakaan
psql -d perpustakaan -c "CREATE SCHEMA perpustakaan;"
```

4. Run database migrations (DDL - see Database Schema section below)
```bash
psql -d perpustakaan -f schema.sql
```

5. Create storage directory
```bash
mkdir -p storage
```

6. Run the application
```bash
go run main.go
```

The server will start on `http://localhost:8080`

## API Endpoints

### Authentication
- `POST /auth` - User login

### Books
- `GET /books` - Get all books (requires authentication)
- `POST /books` - Create a new book (requires authentication)
- `GET /books/:id` - Get book by ID (requires authentication)
- `PUT /books/:id` - Update book (requires authentication)
- `DELETE /books/:id` - Delete book (requires authentication)

### Book Stocks
- `GET /book-stocks` - Get all book stocks (requires authentication)
- `POST /book-stocks` - Create book stock (requires authentication)
- `PUT /book-stocks/:code` - Update book stock (requires authentication)
- `DELETE /book-stocks/:code` - Delete book stock (requires authentication)

### Customers
- `GET /customers` - Get all customers (requires authentication)
- `POST /customers` - Create a new customer (requires authentication)
- `GET /customers/:id` - Get customer by ID (requires authentication)
- `PUT /customers/:id` - Update customer (requires authentication)
- `DELETE /customers/:id` - Delete customer (requires authentication)

### Journals (Borrowing Records)
- `GET /journals` - Get all journals with optional filters (requires authentication)
  - Query params: `customer_id`, `status`
- `POST /journals` - Create borrowing record (requires authentication)
- `PUT /journals/:id` - Return book (requires authentication)

### Media
- `GET /media/:id` - Get media file (requires authentication)
- `POST /media` - Upload media file (requires authentication)

## API Documentation

Complete API documentation with examples is available on Postman:

[View API Documentation](https://documenter.getpostman.com/view/3808786/2sBXc7L4h4)

## Database Schema

### Tables

```sql
-- Users Table
CREATE TABLE perpustakaan.users (
	id varchar(36) DEFAULT gen_random_uuid() NOT NULL,
	email varchar(255) NOT NULL,
	password varchar(255) NOT NULL
);

-- Customers Table
CREATE TABLE perpustakaan.customers (
	id uuid DEFAULT gen_random_uuid() NOT NULL,
	code varchar(50) NULL,
	name varchar(100) NULL,
	created_at timestamp NULL,
	updated_at timestamp NULL,
	deleted_at timestamp NULL,
	CONSTRAINT customers_pkey PRIMARY KEY (id)
);

-- Books Table
CREATE TABLE perpustakaan.books (
	id varchar(36) DEFAULT gen_random_uuid() NOT NULL,
	title varchar(255) NOT NULL,
	description text NULL,
	isbn varchar(100) NOT NULL,
	created_at timestamp(6) NULL,
	updated_at timestamp(6) NULL,
	deleted_at timestamp(6) NULL,
	cover_id varchar(36) NULL,
	CONSTRAINT books_pk PRIMARY KEY (id)
);

-- Book Stocks Table
CREATE TABLE perpustakaan.book_stocks (
	book_id varchar(36) NOT NULL,
	code varchar(50) NOT NULL,
	status varchar(50) NOT NULL,
	borrower_id varchar(36) NULL,
	borrowed_at timestamp(6) NULL,
	CONSTRAINT book_stocks_pk PRIMARY KEY (code)
);

-- Journals Table
CREATE TABLE perpustakaan.journals (
	id varchar(36) DEFAULT gen_random_uuid() NOT NULL,
	book_id varchar(36) NOT NULL,
	stock_code varchar(255) NOT NULL,
	customer_id varchar(36) NOT NULL,
	status varchar(50) NOT NULL,
	borrowed_at timestamp(6) NOT NULL,
	returned_at timestamp(6) NULL,
	due_at timestamp(6) NULL,
	CONSTRAINT journals_pk PRIMARY KEY (id)
);

-- Charges Table
CREATE TABLE perpustakaan.charges (
	id varchar(36) NOT NULL,
	journal_id varchar(36) NOT NULL,
	days_late int4 DEFAULT 1 NOT NULL,
	daily_late_fee int4 NOT NULL,
	total int4 NOT NULL,
	user_id varchar(36) NOT NULL,
	created_at timestamp(6) NULL,
	CONSTRAINT charges_pk PRIMARY KEY (id)
);

-- Media Table
CREATE TABLE perpustakaan.media (
	id varchar(36) DEFAULT gen_random_uuid() NOT NULL,
	path text NULL,
	created_at timestamp(6) NOT NULL,
	CONSTRAINT media_pk PRIMARY KEY (id)
);

-- Sample users
INSERT INTO perpustakaan.users (id,email,"password") VALUES
	 ('4246fb58-ff45-4d2a-8946-93e541fc39fd','admin@perpustakaan.id','$2a$12$Rvslxj25D4OU7w3Ercz/IucMiDkEp1dOCSwq902oWpy0mqcUx2GAq');

```

## License

This project is licensed under the MIT License.
