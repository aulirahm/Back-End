package main

import (
  "net/http"

  "github.com/gin-gonic/gin"
  "gorm.io/driver/sqlite"
  "gorm.io/gorm"
)

// Model = struktur tabel di database
type Book struct {
  ID     uint   `json:"id" gorm:"primaryKey"`
  Title  string `json:"title"`
  Author string `json:"author"`
}

var db *gorm.DB

func main() {
  var err error
  db, err = gorm.Open(sqlite.Open("books.db"), &gorm.Config{})
  if err != nil {
    panic("gagal konek database")
  }
  db.AutoMigrate(&Book{}) // otomatis bikin tabel dari struct Book

  r := gin.Default()

  r.GET("/books", getBooks)       // ambil semua buku
  r.GET("/books/:id", getBook)    // ambil satu buku
  r.POST("/books", createBook)    // tambah buku
  r.PUT("/books/:id", updateBook) // update buku
  r.DELETE("/books/:id", deleteBook) // hapus buku

  r.Run(":8080") // server jalan di localhost:8080
}

func getBooks(c *gin.Context) {
  var books []Book
  db.Find(&books)
  c.JSON(http.StatusOK, books)
}

func getBook(c *gin.Context) {
  var book Book
  if err := db.First(&book, c.Param("id")).Error; err != nil {
    c.JSON(http.StatusNotFound, gin.H{"error": "buku tidak ditemukan"})
    return
  }
  c.JSON(http.StatusOK, book)
}

func createBook(c *gin.Context) {
  var book Book
  if err := c.ShouldBindJSON(&book); err != nil {
    c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
    return
  }
  db.Create(&book)
  c.JSON(http.StatusCreated, book)
}

func updateBook(c *gin.Context) {
  var book Book
  if err := db.First(&book, c.Param("id")).Error; err != nil {
    c.JSON(http.StatusNotFound, gin.H{"error": "buku tidak ditemukan"})
    return
  }
  c.ShouldBindJSON(&book)
  db.Save(&book)
  c.JSON(http.StatusOK, book)
}

func deleteBook(c *gin.Context) {
  db.Delete(&Book{}, c.Param("id"))
  c.JSON(http.StatusOK, gin.H{"message": "berhasil dihapus"})
}
