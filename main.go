package main

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

// DB holds the database connection
var DB *sql.DB


// Book represents a book object
type book struct {
	ID     		string 	`json:"id"`
	Title  		string 	`json:"title"`
	Author 		string 	`json:"author"`
	Quantity 	int 	`json:"quantity"`
}

// Returns a list of all books
func getBooks(c *gin.Context) {
	rows, err := DB.Query("SELECT id, title, author, quantity FROM books")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var books []book
	for rows.Next() {
		var b book
		err := rows.Scan(&b.ID, &b.Title, &b.Author, &b.Quantity)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		books = append(books, b)
	}

	c.JSON(http.StatusOK, books)
}

func getBook(c *gin.Context) {
    id := c.Param("id")
    
    var b book
    
    err := DB.QueryRow("SELECT id, title, author, quantity FROM books WHERE id=$1", id).Scan(&b.ID, &b.Title, &b.Author, &b.Quantity)
    if err != nil {
        if err == sql.ErrNoRows {
            c.JSON(http.StatusNotFound, gin.H{"message": "Book not found"})
            return
        }
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(http.StatusOK, b)
}

// Creates a new book and adds it to the database
func addBook(c *gin.Context) {
	var newBook book

	if err := c.ShouldBindJSON(&newBook); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	sqlStatement := `INSERT INTO books (id, title, author, quantity) VALUES ($1, $2, $3, $4)`
	_, err := DB.Exec(sqlStatement, newBook.ID, newBook.Title, newBook.Author, newBook.Quantity)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, newBook)
}

func viewAllTables(c *gin.Context) {
	rows, err := DB.Query("SELECT table_name FROM information_schema.tables WHERE table_schema = 'public'")
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	
	var tables []string
	for rows.Next() {
		var tableName string
		err := rows.Scan(&tableName)
		if err != nil {
			panic(err)
		}
		tables = append(tables, tableName)
	}
	
	if err = rows.Err(); err != nil {
		panic(err)
	}
	
	fmt.Println(tables)
}

// Deletes a book from the database
func deleteBook(c *gin.Context) {
	id := c.Param("id")

	// Delete book from the database
	res, err := DB.Exec("DELETE FROM books WHERE id=$1", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "Book not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Book deleted"})
}

func updateBook(c *gin.Context) {
	id := c.Param("id")

	var updatedBook book
	if err := c.ShouldBindJSON(&updatedBook); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var b book
	row := DB.QueryRow("SELECT id, title, author, quantity FROM books WHERE id = $1", id)
	err := row.Scan(&b.ID, &b.Title, &b.Author, &b.Quantity)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"message": "Book not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	b.Title = updatedBook.Title
	b.Author = updatedBook.Author
	b.Quantity = updatedBook.Quantity

	_, err = DB.Exec("UPDATE books SET title=$1, author=$2, quantity=$3 WHERE id=$4", b.Title, b.Author, b.Quantity, b.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, b)
}

func main() {
	
	router := gin.Default()
	dsn := "postgres://gouser:admingres@localhost:5432/godb?sslmode=disable"

	var err error
	DB, err = sql.Open("postgres", dsn)
	if err != nil {
		panic(err)
	}

	router.GET("/tables", viewAllTables)
	router.GET("/books", getBooks)
	router.GET("/books/:id", getBook)
	router.POST("/books", addBook)
	router.PUT("/books/:id", updateBook)
	router.DELETE("/books/:id", deleteBook)

	router.Run( "localhost:8080")
}
