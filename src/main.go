package main

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "github.com/jinzhu/gorm"
    _ "github.com/jinzhu/gorm/dialects/sqlite"
)

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
    r := gin.Default()

    r.GET("/api/books", getBooks)
    r.GET("/api/book/:id", getBook)
    r.POST("/api/create-book", createBook)
    r.PATCH("/api/update-book/:id", updateBook)
    r.DELETE("/api/delete-book/:id", deleteBook)

    connectDB()

    r.Run()
}

// mock database connection
func connectDB() {
    database, err := gorm.Open("sqlite3", "test.db")

    if err != nil {
        panic("Failed to connect.")
    }

    database.AutoMigrate(&Book{})

    DB = database
}

// get all books
func getBooks(c *gin.Context) {
    var books []Book
    DB.Find(&books)

    c.JSON(http.StatusOK, gin.H{"respond": books})
}

// get single book
func getBook(c *gin.Context) {
    var book Book

    if err := DB.Where("id = ?", c.Param("id")).First(&book).Error; err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "No item with given ID"})
            return
    }

    c.JSON(http.StatusOK, gin.H{"respond": book})
}


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

// updsssssss
func updateBook(c *gin.Context) {
    var book Book

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

func deleteBook(c *gin.Context) {
    var book Book

    if err := DB.Where("id = ?", c.Param("id")).First(&book).Error; err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "No item with given ID"})
            return
    }

    DB.Delete(&book)
    c.JSON(http.StatusOK, gin.H{"response": true})
}
