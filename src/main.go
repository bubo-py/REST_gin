package main

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "github.com/jinzhu/gorm"
    _ "github.com/jinzhu/gorm/dialects/sqlite"
)

// gorm variable to initialize the database
var DB *gorm.DB

type Book struct {
    ID uint `json:"id" gorm:"primary_key"`
    Title string `json:"title"`
    Author string `json:"author"`
}

type BookCreationSchema struct { // fields are required
    Title string `json:"title" binding:"required"`
    Author string `json:"author" binding:"required"`
}

type BookUpdateSchema struct { // fields don't need to be required
    Title string `json:"title"`
    Author string `json:"author"`
}

func main()  {
    // Setting router
    r := gin.Default()

    // API routes with corresponding functions
    r.GET("/api/books", getBooks)
    r.GET("/api/book/:id", getBook)
    r.POST("/api/create-book", createBook)
    r.PATCH("/api/update-book/:id", updateBook)
    r.DELETE("/api/delete-book/:id", deleteBook)

    // Database connection via function below
    connectDB()

    r.Run()
}

// Mock database connection
func connectDB() {
    // making db with gorm
    database, err := gorm.Open("sqlite3", "test.db")

    if err != nil {
        panic("Failed to connect.")
    }

    database.AutoMigrate(&Book{})

    DB = database
}

// Getting all books
func getBooks(c *gin.Context) {
    var books []Book
    DB.Find(&books)

    c.JSON(http.StatusOK, gin.H{"respond": books})
}

// Getting single, specified book
func getBook(c *gin.Context) {
    var book Book

    if err := DB.Where("id = ?", c.Param("id")).First(&book).Error; err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "No item with given ID"})
            return
    }

    c.JSON(http.StatusOK, gin.H{"respond": book})
}

// Adding book function
func createBook(c *gin.Context) {
    // Input validation
    var input BookCreationSchema

    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Book creation
    book := Book{Title: input.Title, Author: input.Author}
    DB.Create(&book)

    c.JSON(http.StatusOK, gin.H{"respond": book})
}

// Updating book function
func updateBook(c *gin.Context) {
    var book Book

    // Checks if given ID exists, if doesn't - return error
    if err := DB.Where("id = ?", c.Param("id")).First(&book).Error; err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "No item with given ID"})
            return
    }

    // Input validation
    var input BookUpdateSchema

    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    DB.Model(&book).Updates(input)
    c.JSON(http.StatusOK, gin.H{"response": book})
}

// Deleting book function
func deleteBook(c *gin.Context) {
    var book Book

    if err := DB.Where("id = ?", c.Param("id")).First(&book).Error; err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "No item with given ID"})
            return
    }

    DB.Delete(&book)
    c.JSON(http.StatusOK, gin.H{"response": "Item has been successfully deleted"})
}
