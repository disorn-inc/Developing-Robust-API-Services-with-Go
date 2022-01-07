package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"golang.org/x/time/rate"
	"github.com/gin-contrib/cors"

	"github/disorn-inc/Developing-Robust-API-Services-with-Go/Todo/auth"
	"github/disorn-inc/Developing-Robust-API-Services-with-Go/Todo/todo"
	"github/disorn-inc/Developing-Robust-API-Services-with-Go/Todo/store"
	"github/disorn-inc/Developing-Robust-API-Services-with-Go/Todo/router"
)

var (
	buildcommit = "dev"
	buildtime = time.Now().String()
)

func main() {
	_, err := os.Create("/tmp/live")

	if err != nil {
		log.Fatal(err)
	}

	defer os.Remove("/tmp/live")

	err = godotenv.Load("local.env")
	if err != nil {
		log.Println("please consider environment variable: %s", err)
	}

	db, err := gorm.Open(mysql.Open(os.Getenv("DB_CONN")), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&todo.Todo{})
	r := gin.Default()
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{
		"http://localhost:8080",
	}
	config.AllowHeaders = []string{
		"Origin",
		"Authorization",
		"TransactionID",
	}
	r.Use(cors.New(config))

	r.GET("/healthz", func(c *gin.Context) {
		c.Status(200)
	})
	r.GET("/limitz", limitedHandler)
	r.GET("/x", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"buildcommit": buildcommit,
			"buildtime": buildtime,
		})
	})

	r.GET("/tokenz", auth.AccessToken(os.Getenv("SIGN")))

	protected := r.Group("", auth.Protect([]byte(os.Getenv("SIGN"))))

	gormStore := store.NewGormStore(db)
	
	handler := todo.NewTodoHandler(gormStore)
	protected.POST("/todos", router.NewGinHandler(handler.NewTask))
	protected.GET("/todos", router.NewGinHandler(handler.List))
	protected.DELETE("/todos/:id", router.NewGinHandler(handler.Remove))

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	s := &http.Server{
		Addr: ":" + os.Getenv("PORT"),
		Handler: r,
		ReadTimeout: 10 * time.Second,
		WriteTimeout: 10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	go func ()  {
		if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	<-ctx.Done()
	stop()
	fmt.Println("shutting down gracefully, press ctrl+c again to force")

	timeoutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.Shutdown(timeoutCtx); err != nil {
		fmt.Println(err)
	}
}

var limiter = rate.NewLimiter(5, 5)

func limitedHandler(c *gin.Context) {
	if !limiter.Allow() {
		c.AbortWithStatus(http.StatusTooManyRequests)
		return
	}
	c.JSON(200, gin.H{
		"message": "pong",
	})
}