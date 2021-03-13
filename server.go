package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"
)

type Todo struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Status string `json:"status"`
}

var todos = map[int]*Todo{
	1: &Todo{ID: 1, Title: "I am Darunee", Status: "active"},
}

func helloHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{
		"message": "hello Darunee",
	})

}

func getTodosHandler(c echo.Context) error {
	items := []*Todo{}
	for _, item := range todos {
		items = append(items, item)

	}
	return c.JSON(http.StatusOK, items)

}

func createTodosHandler(e echo.Context) error {
	t := Todo{}

	if err := e.Bind(&t); err != nil {
		return e.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})

	}
	id := len(todos)
	id++
	t.ID = id
	todos[t.ID] = &t
	return e.JSON(http.StatusCreated, "created todo.")

}

func getTodoByIDHandler(c echo.Context) error {
	var id int
	//var title, status string
	err := echo.PathParamsBinder(c).Int("id", &id).BindError()
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal("Connect to database error", err)
	}
	defer db.Close()

	stmt, err := db.Prepare("SELECT id, title, status FROM todos where id=$1")
	if err != nil {
		log.Fatal("can'tprepare query one row statment", err)
	}

	row := stmt.QueryRow(id)
	//var id int
	var title, status string

	err = row.Scan(&id, &title, &status)
	if err != nil {
		log.Fatal("can't Scan row into variables", err)
	}

	t := Todo{
		ID:     id,
		Title:  title,
		Status: status,
	}

	return c.JSON(http.StatusOK, t)

}

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/hello", helloHandler)
	e.GET("/todos", getTodosHandler)
	e.GET("/todos/:id", getTodoByIDHandler)
	e.POST("/todos", createTodosHandler)

	port := os.Getenv("PORT")
	log.Println("port:", port)
	e.Start(":" + port)

}
