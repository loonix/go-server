package bookTypes

// Book represents a book object
type Book struct {
    ID          string      `json:"id"`
    Title       string      `json:"title"`
    Author      string      `json:"author"`
    Quantity    int         `json:"quantity"`
}