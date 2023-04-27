package routers

import (
	"DTSGolang/Kelas2/Assignment7/controllers"

	"github.com/gin-gonic/gin"
)

func StartServer() *gin.Engine {
	router := gin.Default()
	router.POST("/books", controllers.CreateBook)
	router.PUT("/books/:BookID", controllers.UpdateBook)
	router.GET("/books/:BookID", controllers.GetBook)
	router.GET("/books/all", controllers.GetBooks)
	router.DELETE("/books/:BookID", controllers.DeleteBook)

	return router
}
