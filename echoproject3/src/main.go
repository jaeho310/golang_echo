package main

import (
	"html/template"
	"io"
	"log"
	"net/http"
	"sort"
	"strconv"

	"github.com/labstack/echo/v4"
)

// TemplateRenderer is a custom html/template renderer for Echo framework
type TemplateRenderer struct {
	templates *template.Template
}

// Render renders a template document
func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {

	// Add global methods if data is a map
	if viewContext, isMap := data.(map[string]interface{}); isMap {
		viewContext["reverse"] = c.Echo().Reverse
	}

	return t.templates.ExecuteTemplate(w, name, data)
}

type (
	user struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}
)

var (
	userDB = map[int]*user{}
	seq    = 1
)

func getAllUsers(c echo.Context) error {
	log.Println("select all users")
	var users = []user{}
	for _, val := range userDB {
		users = append(users, *val)
	}

	sort.Slice(users, func(i, j int) bool {
		return users[i].ID < users[j].ID
	})
	return c.JSON(http.StatusOK, users)
	// return c.JSON(http.StatusOK, userDB)
}

func createUser(c echo.Context) error {
	u := &user{
		ID: seq,
	}
	if err := c.Bind(u); err != nil {
		return err
	}
	userDB[u.ID] = u
	seq++
	log.Println("insert new user : ", u.Name)
	return c.JSON(http.StatusCreated, u)
}

func deleteUsers(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	log.Println("delete user : ", id)
	delete(userDB, id)
	return c.NoContent(http.StatusNoContent)
}

func main() {
	e := echo.New()
	renderer := &TemplateRenderer{
		templates: template.Must(template.ParseGlob("resources/templates/*.html")),
	}
	e.Renderer = renderer
	e.Static("/static", "resources/static")

	e.GET("/", func(c echo.Context) error {
		return c.Render(http.StatusOK, "main.html", map[string]interface{}{})
	})
	e.GET("/list", func(c echo.Context) error {
		return c.Render(http.StatusOK, "list.html", map[string]interface{}{})
	})

	e.POST("/api/users", createUser)
	e.GET("/api/users", getAllUsers)
	e.DELETE("/api/users/:id", deleteUsers)

	e.Logger.Fatal(e.Start(":8395"))
}
