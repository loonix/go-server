package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Book represents a book object
type book struct {
	ID     		string 	`json:"id"`
	Title  		string 	`json:"title"`
	Author 		string 	`json:"author"`
	Quantity 	int 	`json:"quantity"`
}

// Mock database (temporary)
var books []book = []book{
	{ID: "1", Title: "The Alchemist", Author: "Paulo Coelho", Quantity: 10},
	{ID: "2", Title: "The Monk Who Sold His Ferrari", Author: "Robin Sharma", Quantity: 5},
	{ID: "3", Title: "The Secret", Author: "Rhonda Byrne", Quantity: 7},
	{ID: "4", Title: "The Power of Your Subconscious Mind", Author: "Joseph Murphy", Quantity: 3},
	{ID: "5", Title: "The Power of Now", Author: "Eckhart Tolle", Quantity: 2},
}

// Returns a list of all books
func getBooks(c *gin.Context) {
	c.JSON(http.StatusOK, books)
}

// Returns a single book by ID
func getBook(c *gin.Context) {
	id := c.Param("id")

	for _, b := range books {
		if b.ID == id {
			c.JSON(http.StatusOK, b)
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"message": "Book not found"})
}

// Creates a new book and adds it to the books slice (temporary)
func addBook(c *gin.Context) {
	var newBook book

	if err := c.ShouldBindJSON(&newBook); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	books = append(books, newBook)
	c.JSON(http.StatusCreated, newBook)
}

// Deletes a book from the books slice (temporary)
func deleteBook(c *gin.Context) {
	id := c.Param("id")

	for i, b := range books {
		if b.ID == id {
			books = append(books[:i], books[i+1:]...)
			c.JSON(http.StatusOK, gin.H{"message": "Book deleted"})
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"message": "Book not found"})
}

// Updates a book from the books slice (temporary)
func updateBook(c *gin.Context) {
	id := c.Param("id")

	var updatedBook book
	if err := c.ShouldBindJSON(&updatedBook); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	for i, b := range books {
		if b.ID == id {
			books[i].Title = updatedBook.Title
			books[i].Author = updatedBook.Author
			books[i].Quantity = updatedBook.Quantity
			c.JSON(http.StatusOK, books[i])
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"message": "Book not found"})
}


func main() {
	router := gin.Default()

	router.GET("/books", getBooks)
	router.GET("/books/:id", getBook)
	router.POST("/books", addBook)
	router.PUT("/books/:id", updateBook)
	router.DELETE("/books/:id", deleteBook)

	router.Run( "localhost:8080")
}