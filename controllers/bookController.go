package controllers

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "1234"
	dbname   = "db-go-sql"
)

var (
	db  *sql.DB
	err error
)

func StartDB() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected to database")
}

type Book struct {
	BookID      string `json:"id"`
	Title       string `json:"title"`
	Author      string `json:"author"`
	Description string `json:"description"`
}

func CreateBook(ctx *gin.Context) {
	var book = Book{}

	type BookInput struct {
		Title       string `json:"title" binding:"required"`
		Author      string `json:"author" binding:"required"`
		Description string `json:"description" binding:"required"`
	}

	var newBook BookInput

	if err := ctx.ShouldBindJSON(&newBook); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	sqlStatement := `
	INSERT INTO books (title, author, description)
	VALUES ($1, $2, $3)
	Returning *
	`

	defer func() {
		if err := recover(); err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, err)
			return
		}
	}()

	err = db.QueryRow(sqlStatement, newBook.Title, newBook.Author, newBook.Description).
		Scan(&book.BookID, &book.Title, &book.Author, &book.Description)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"messsage": "Book created successfully",
	})
}

func UpdateBook(ctx *gin.Context) {
	// parse nilai parameter ID dari URL
	BookID, err := strconv.Atoi(ctx.Param("BookID"))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Invalid ID",
		})
		return
	}

	type BookInput struct {
		Title       string `json:"title" binding:"required"`
		Author      string `json:"author" binding:"required"`
		Description string `json:"description" binding:"required"`
	}

	var updateBook BookInput
	// bind request body ke dalam variabel updateBook
	if err := ctx.BindJSON(&updateBook); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// update data buku di dalam database
	sqlStatement := `
	UPDATE books 
	SET title=$2, author=$3, description=$4
	WHERE id=$1
	`

	defer func() {
		if err := recover(); err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, err)
			return
		}
	}()

	_, err = db.Exec(sqlStatement, BookID, updateBook.Title, updateBook.Author, updateBook.Description)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, err)
		return
	}

	// Kembalikan respon sukses
	ctx.JSON(http.StatusOK, gin.H{
		"id":      BookID,
		"message": "Book updated successfully",
	})
}

func GetBook(ctx *gin.Context) {
	// parse nilai parameter ID dari URL
	BookID, err := strconv.Atoi(ctx.Param("BookID"))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Invalid ID",
		})
		return
	}

	var book = Book{}

	sqlStatement := `SELECT * FROM books WHERE id=$1`

	defer func() {
		if err := recover(); err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, err)
			return
		}
	}()

	err = db.QueryRow(sqlStatement, BookID).
		Scan(&book.BookID, &book.Title, &book.Author, &book.Description)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, err)
		return
	}

	ctx.JSON(http.StatusOK, book)
}

func GetBooks(ctx *gin.Context) {
	// buat slice untuk menampug data buku
	var books []Book

	sqlStatement := `SELECT * FROM books`

	defer func() {
		if err := recover(); err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, err)
			return
		}
	}()

	// query semua data buku
	rows, _ := db.Query(sqlStatement)

	// iterasi setiap baris hasil query
	for rows.Next() {
		// buat variabel unntuk menampung data buku pada setiap baris
		var book = Book{}
		err := rows.Scan(&book.BookID, &book.Title, &book.Author, &book.Description)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, err)
			return
		}

		// tambahkan data buku ke dalam slice
		books = append(books, book)
	}

	ctx.JSON(http.StatusOK, books)
}

func DeleteBook(ctx *gin.Context) {
	// parse nilai parameter ID dari URL
	BookID, err := strconv.Atoi(ctx.Param("BookID"))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Invalid ID",
		})
	}

	sqlStatement := `DELETE FROM books WHERE id = $1;`

	defer func() {
		if err := recover(); err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, err)
			return
		}
	}()

	_, err = db.Exec(sqlStatement, BookID)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("Book with id %v has been deleted", BookID),
	})
}
