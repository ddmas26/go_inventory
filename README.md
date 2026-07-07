# Go Inventory

A generic, type-safe CRUD and paginated query library for PostgreSQL, built with Go generics and reflection.

## Features

- **Generic Repository** — type-safe CRUD operations for any struct using Go 1.26 generics
- **Safe Paginated Queries** — SQL injection protected via parameterized queries and whitelist validation
- **Dynamic Filtering** — filter with `AND`/`OR` combinations, comparison operators (`=`, `>`, `<`, `>=`, `<=`)
- **Full-Text Search** — `ILIKE` search across specified columns
- **Ordered Pagination** — safe `ORDER BY` with whitelist-validated columns and directions
- **Count Query** — returns total matching rows for frontend pagination controls
- **Auto-Migration** — runs SQL migration files in order from `migrate/migrations/`

## Project Structure

```
├── cmd/
│   └── main.go                    # Entry point
├── database/
│   ├── connect.go                 # DB connection (Postgres)
│   ├── domain.go                  # Domain models (User, Inventory)
│   ├── dto.go                     # Pagination request/response DTOs
│   └── repo.go                    # Generic repository implementation
├── migrate/
│   ├── migrations/                # SQL migration files (numbered folders)
│   │   ├── 0_01072026/           # Initial users table
│   │   ├── 1_04072026/           # Inventory table
│   │   ├── 2_05072026/           # Add created_at/updated_at to users
│   │   └── 3_06072026/           # (placeholder)
│   └── run/
│       └── run.go                 # Migration runner
├── go.mod
└── README.md
```

## Getting Started

### Prerequisites

- Go 1.26+
- PostgreSQL

### Environment Variables

Create a `.env` file in the project root:

```env
HOST=localhost
PORT=5432
DB_USER=postgres
PASSWORD=your_password
DBNAME=inventory_db
```

### Run Migrations

```bash
go run ./migrate/run/
```

This executes `up.sql` files in each numbered folder sequentially, with automatic rollback on failure.

### Run the App

```bash
go run ./cmd/
```

## Usage

### Define a Model

Tag struct fields with `db` and optionally `primary`:

```go
type Inventory struct {
    ID        uuid.UUID  `db:"id" primary:"true"`
    Name      string     `db:"name"`
    Quantity  int        `db:"quantity"`
    Price     float64    `db:"price"`
    CreatedAt time.Time  `db:"created_at"`
    UpdatedAt *time.Time `db:"updated_at"`
}
```

### Create a Repository

```go
repo := database.NewGenericRepo[database.Inventory](db, "inventory")
```

### CRUD Operations

```go
// Create
created, err := repo.Create(&entity, db)

// Find All
all, err := repo.FindAll(db)

// Find By ID
item, err := repo.FindById(id, db)

// Delete
err := repo.Delete(id, db)
```

### Paginated Queries

```go
req := database.PaginationRequest{
    PageIndex:      1,
    PageSize:       10,
    Filter:         "price: > 10, quantity: < 100|quantity: > 500",
    SearchBy:       []string{"name"},
    SearchValue:    "widget",
    OrderBy:        "price",
    OrderDirection: "DESC",
}

result, err := repo.FindAllPaginated(req, db)
// result.Data        — items for current page
// result.TotalCount  — total matching rows
// result.PageIndex   — current page
// result.PageSize    — items per page
// result.TotalPages  — total number of pages
```

## Filter Syntax

The filter string supports `AND` (comma) and `OR` (pipe) grouping:

```
field:value                      → field = value
field: > value                   → field > value
field: < value                   → field < value
field: >= value                  → field >= value
field: <= value                  → field <= value
field: val1|field: val2          → (field = val1 OR field = val2)  [OR group]
field: val1 , field2: val2       → field = val1 AND field2 = val2  [AND groups]
```

All field names are validated against the model's `db` tags — SQL injection via column names is prevented.

## Security

- **Parameterized queries** — all values are passed via `$N` placeholders
- **Whitelist validation** — column names for `OrderBy`, `Filter`, and `SearchBy` are validated against known struct fields
- **No raw string interpolation** — user input never appears directly in SQL
