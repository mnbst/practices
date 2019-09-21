package main

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"godemo/controller"
	"godemo/database"
	"godemo/model"
	"gopkg.in/bluesuncorp/validator.v5"
	"net/http"
	"os"
	"unicode"
)

const defaultPort = "8080"

var (
	msgInvalidJSON     = "Invalid JSON format"
	msgInvalidJSONType = func(e *json.UnmarshalTypeError) string {
		return "Expected " + e.Value + " but given type is " + e.Type.String() + " in JSON"
	}
	msgValidationFailed = func(e *validator.FieldError) string {
		switch e.Tag {
		case "required":
			return toSnakeCase(e.Field) + ": required"
		case "max":
			return toSnakeCase(e.Field) + ": too_long"
		default:
			return e.Error()
		}
	}
)

func main() {
	migrate()

	router := gin.Default()
	router.Static("/css", "./assets/dist/css")
	router.Static("/js", "./assets/dist/js")
	router.LoadHTMLGlob("templates/*")

	router.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	router.GET("/", controller.Users.Top)
	router.GET("/login", controller.Users.Login)
	router.GET("/logout", controller.Users.Logout)
	router.GET("/register", controller.Users.Register)
	router.POST("/authenticate", controller.Users.Authenticate)
	router.POST("/users/create", controller.Users.Create)

	router.POST("/postTodo", controller.Todos.Create)
	router.PATCH("/done", controller.Todos.Done)
	router.DELETE("/delete", controller.Todos.Delete)
	router.GET("/delregistration", controller.Users.Deluser)



	http.ListenAndServe(":"+port(), router)
}


// https://gist.github.com/elwinar/14e1e897fdbe4d3432e1
func toSnakeCase(in string) string {
	runes := []rune(in)
	length := len(runes)

	var out []rune
	for i := 0; i < length; i++ {
		if i > 0 && unicode.IsUpper(runes[i]) && ((i+1 < length && unicode.IsLower(runes[i+1])) || unicode.IsLower(runes[i-1])) {
			out = append(out, '_')
		}
		out = append(out, unicode.ToLower(runes[i]))
	}

	return string(out)
}

func port() string {
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = defaultPort
	}

	return port
}

func migrate() {
	db := database.GetDB()

	db.AutoMigrate(&model.User{}, &model.Todo{})
}
