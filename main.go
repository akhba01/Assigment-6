package main

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "postgres"
	dbname   = "db-go-sql"
)

var (
	db  *sql.DB
	err error
)

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Succesfully connected to database")
	r := gin.Default()
	r.GET("/books", GetAllBook)
	r.GET("/books/:idBook", GetBookById)
	r.POST("/books", CreateBook)
	r.PUT("/books/:idBook", UpdatedBookById)
	r.DELETE("/books/:idBook", DeleteBookById)
	r.Run(":4000")
}

type Books struct {
	BookId int    `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
	Desc   string `json:"desc"`
}

func GetAllBook(ctx *gin.Context) {
	var results = []Books{}
	sqlStatement := `SELECT * FROM books`

	rows, err := db.Query(sqlStatement)

	if err != nil {
		panic(err)
	}

	defer rows.Close()

	for rows.Next() {
		var book = Books{}

		err := rows.Scan(&book.BookId, &book.Title, &book.Author, &book.Desc)

		if err != nil {
			panic(err)
		}

		results = append(results, book)
	}
	ctx.JSON(http.StatusOK, results)
	fmt.Println(results)
}

func GetBookById(ctx *gin.Context) {
	id := ctx.Param("idBook")
	sqlStatement := `SELECT * FROM books WHERE id = $1`

	rows, err := db.Query(sqlStatement, id)

	if err != nil {
		panic(err)
	}

	var book = Books{}

	if rows.Next() {
		err = rows.Scan(&book.BookId, &book.Title, &book.Author, &book.Desc)

		if err != nil {
			panic(err)
		}
	} else {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
		return
	}

	fmt.Println(book)
	ctx.JSON(http.StatusOK, book)
}

type CreateBookInput struct {
	Title  string `json:"title"`
	Author string `json:"author"`
	Desc   string `json:"desc"`
}

func CreateBook(ctx *gin.Context) {
	input := CreateBookInput{}
	var newBook Books

	sqlStatement := `
	INSERT INTO books (title, author, description)
	VALUES ($1,$2,$3)
	Returning *
	`
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	err = db.QueryRow(sqlStatement, input.Title, input.Author, input.Desc).Scan(&newBook.BookId, &newBook.Title, &newBook.Author, &newBook.Desc)

	if err != nil {
		fmt.Println(err)
		panic(err)

	}

	ctx.JSON(http.StatusCreated, newBook)
	fmt.Println(newBook)
}

type UpdateBookInput struct {
	Title  string `json:"title"`
	Author string `json:"author"`
	Desc   string `json:"desc"`
}

func UpdatedBookById(ctx *gin.Context) {
	input := UpdateBookInput{}
	id := ctx.Param("idBook")
	sqlStatement := `
	UPDATE books
	SET title = $2, author = $3, description = $4
	WHERE id = $1
	`
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	res, err := db.Exec(sqlStatement, &id, &input.Title, &input.Author, &input.Desc)
	if err != nil {
		panic(err)
	}
	count, err := res.RowsAffected()
	if err != nil {
		panic(err)
	}
	fmt.Println("Updated data amount", count)
	ctx.JSON(http.StatusOK, "Updated")

}

func DeleteBookById(ctx *gin.Context) {
	id := ctx.Param("idBook")
	fmt.Println(id)
	// idBook, _ := strconv.Atoi(id)
	sqlStatement := `
	DELETE FROM books
	WHERE id = $1;
	`
	res, err := db.Exec(sqlStatement, id)
	if err != nil {
		panic(err)
	}
	count, err := res.RowsAffected()
	if err != nil {
		panic(err)
	}
	fmt.Println("Deleted data amount ", count)

	ctx.JSON(http.StatusOK, "Deleted")

}
