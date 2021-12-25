package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github/disorn-inc/Developing-Robust-API-Services-with-Go/Todo/auth"
	"github/disorn-inc/Developing-Robust-API-Services-with-Go/Todo/todo"
)

func main() {
	err := godotenv.Load("local.env")
	if err != nil {
		log.Panicln("please consider environment variable: %s", err)
	}

	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&todo.Todo{})
	r := gin.Default()

	r.GET("/tokenz", auth.AccessToken(os.Getenv("SIGN")))

	protected := r.Group("", auth.Protect([]byte(os.Getenv("SIGN"))))
	
	handler := todo.NewTodoHandler(db)
	protected.POST("/todos", handler.NewTask)
	r.Run()
}