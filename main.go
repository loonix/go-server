package main

import (
	"context"
	"database/sql"
	"encoding/json"
	bookTypes "go-server/models"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	_ "github.com/lib/pq"
)

// Database is a struct for managing database connections
type Database struct {
    *sql.DB
}

// Error is a custom error type for database errors
type Error string

func (e Error) Error() string {
    return string(e)
}

// Common database errors
const (
    ErrNotFound = Error("not found")
    ErrInternal = Error("internal error")
)

// Returns a list of all books
func (db *Database) GetBooks(ctx context.Context) ([]bookTypes.Book, error) {
    rows, err := db.QueryContext(ctx, "SELECT id, title, author, quantity FROM books")
    if err != nil {
        return nil, ErrInternal
    }
    defer rows.Close()

    var books []bookTypes.Book
    for rows.Next() {
        var b bookTypes.Book
        err := rows.Scan(&b.ID, &b.Title, &b.Author, &b.Quantity)
        if err != nil {
            return nil, ErrInternal
        }
        books = append(books, b)
    }

    if err = rows.Err(); err != nil {
        return nil, ErrInternal
    }

    return books, nil
}

// Returns a single book with the given ID
func (db *Database) GetBook(ctx context.Context, id string) (*bookTypes.Book, error) {
    var b bookTypes.Book

    err := db.QueryRowContext(ctx, "SELECT id, title, author, quantity FROM books WHERE id=$1", id).Scan(&b.ID, &b.Title, &b.Author, &b.Quantity)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, ErrNotFound
        }
        return nil, ErrInternal
    }

    return &b, nil
}

// Creates a new book and adds it to the database
func (db *Database) AddBook(ctx context.Context, newBook *bookTypes.Book) (*bookTypes.Book, error) {
    _, err := db.ExecContext(ctx, "INSERT INTO books (id, title, author, quantity) VALUES ($1, $2, $3, $4)", newBook.ID, newBook.Title, newBook.Author, newBook.Quantity)
    if err != nil {
        return nil, ErrInternal
    }

    return newBook, nil
}

// Deletes a book from the database
func (db *Database) DeleteBook(ctx context.Context, id string) error {
    res, err := db.ExecContext(ctx, "DELETE FROM books WHERE id=$1", id)
    if err != nil {
        return ErrInternal
    }

    rowsAffected, err := res.RowsAffected()
    if err != nil {
        return ErrInternal
    }

    if rowsAffected == 0 {
        return ErrNotFound
    }

    return nil
}

// Updates a book in the database
func (db *Database) UpdateBook(ctx context.Context, id string, updatedBook *bookTypes.Book) (*bookTypes.Book, error) {
    var b bookTypes.Book
    row := db.QueryRowContext(ctx, "SELECT id, title, author, quantity FROM books WHERE id = $1", id)
    err := row.Scan(&b.ID, &b.Title, &b.Author, &b.Quantity)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, ErrNotFound
        }
        return nil, ErrInternal
    }

    b.Title = updatedBook.Title
    b.Author = updatedBook.Author
    b.Quantity = updatedBook.Quantity

    _, err = db.ExecContext(ctx, "UPDATE books SET title=$1, author=$2, quantity=$3 WHERE id=$4", b.Title, b.Author, b.Quantity, b.ID)
    if err != nil {
        return nil, ErrInternal
    }

    return &b, nil
}

// Views all tables in the public schema
func (db *Database) ViewAllTables(ctx context.Context) ([]string, error) {
    rows, err := db.QueryContext(ctx, "SELECT table_name FROM information_schema.tables WHERE table_schema = 'public'")
    if err != nil {
        return nil, ErrInternal
    }
    defer rows.Close()

    var tables []string
    for rows.Next() {
        var tableName string
        err := rows.Scan(&tableName)
        if err != nil {
            return nil, ErrInternal
        }
        tables = append(tables, tableName)
    }

    if err = rows.Err(); err != nil {
        return nil, ErrInternal
    }

    return tables, nil
}

func main() {

    // Initialize database
    dsn := "postgres://gouser:admingres@localhost:5432/godb?sslmode=disable"
    db, err := sql.Open("postgres", dsn)
    if err != nil {
        log.Fatal(err)
    }

    database := &Database{db}

    // Initialize router
    router := chi.NewRouter()
    router.Use(middleware.Logger)
    router.Use(middleware.Recoverer)

    // API endpoints
    router.Get("/tables", func(w http.ResponseWriter, r *http.Request) {
        tables, err := database.ViewAllTables(r.Context())
        if err != nil {
            switch err {
            case ErrNotFound:
                w.WriteHeader(http.StatusNotFound)
            case ErrInternal:
                w.WriteHeader(http.StatusInternalServerError)
            }
            return
        }

        json.NewEncoder(w).Encode(tables)
    })

    router.Get("/books", func(w http.ResponseWriter, r *http.Request) {
        books, err := database.GetBooks(r.Context())
        if err != nil {
            switch err {
            case ErrInternal:
                w.WriteHeader(http.StatusInternalServerError)
				http.Error(w, "internal error", http.StatusInternalServerError)
            }
            return
        }

        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusOK)
        json.NewEncoder(w).Encode(books)
    })

    router.Get("/books/{id}", func(w http.ResponseWriter, r *http.Request) {
        id := chi.URLParam(r, "id")

        b, err := database.GetBook(r.Context(), id)
        if err != nil {
            switch err {
			case ErrNotFound:
				w.WriteHeader(http.StatusNotFound)
				http.Error(w, "book not found", http.StatusNotFound)
            case ErrInternal:
                w.WriteHeader(http.StatusInternalServerError)
				http.Error(w, "internal error", http.StatusInternalServerError)
            }
            return
        }

        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusOK)
        json.NewEncoder(w).Encode(b)
    })

    router.Post("/books", func(w http.ResponseWriter, r *http.Request) {
        var newBook bookTypes.Book
        err := json.NewDecoder(r.Body).Decode(&newBook)
        if err != nil {
            w.WriteHeader(http.StatusBadRequest)
			http.Error(w, "error - bad request", http.StatusBadRequest)
            return
        }

        b, err := database.AddBook(r.Context(), &newBook)
        if err != nil {
            switch err {
            case ErrInternal:
                w.WriteHeader(http.StatusInternalServerError)
				http.Error(w, "internal error", http.StatusInternalServerError)

            }
            return
        }

        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusCreated)
        json.NewEncoder(w).Encode(b)
    })

    router.Put("/books/{id}", func(w http.ResponseWriter, r *http.Request) {
        id := chi.URLParam(r, "id")

        var updatedBook bookTypes.Book
        if err := json.NewDecoder(r.Body).Decode(&updatedBook); err != nil {
            w.WriteHeader(http.StatusBadRequest)
			http.Error(w, "error - bad request", http.StatusBadRequest)
            return
        }

        b, err := database.UpdateBook(r.Context(), id, &updatedBook)
        if err != nil {
            switch err {
            case ErrNotFound:
				w.WriteHeader(http.StatusNotFound)
				http.Error(w, "book not found", http.StatusNotFound)
            case ErrInternal:
                w.WriteHeader(http.StatusInternalServerError)
				http.Error(w, "internal error", http.StatusInternalServerError)
            }
            return
        }

        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusOK)
        json.NewEncoder(w).Encode(b)
    })

	router.Delete("/books/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
	
		ctx := r.Context()
		err := database.DeleteBook(ctx, id)
		if err != nil {
			switch err {
			case ErrNotFound:
				w.WriteHeader(http.StatusNotFound)
				http.Error(w, "book not found", http.StatusNotFound)
			case ErrInternal:
				w.WriteHeader(http.StatusInternalServerError)
				http.Error(w, "internal error", http.StatusInternalServerError)
			}
			return
		}
	
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "book deleted"})
	})

    log.Fatal(http.ListenAndServe(":8080", router))
}