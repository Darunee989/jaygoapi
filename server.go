package main

import (
	"log"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
)

func helloHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{
		"message": "hello Darunee",
	})

}

func main() {
	e := echo.New()
	e.GET("/hello", helloHandler)
	port := os.Getenv("PORT")
	log.Println("port:", port)
	e.Start(":" + port)

}
