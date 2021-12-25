package main

import (
	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github/disorn-inc/Developing-Robust-API-Services-with-Go/Todo/auth"
	"github/disorn-inc/Developing-Robust-API-Services-with-Go/Todo/todo"
)

func main() {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&todo.Todo{})
	r := gin.Default()

	r.GET("/tokenz", auth.AccessToken("==signature=="))

	protected := r.Group("", auth.Protect([]byte("==signature==")))
	
	handler := todo.NewTodoHandler(db)
	protected.POST("/todos", handler.NewTask)
	r.Run()
}