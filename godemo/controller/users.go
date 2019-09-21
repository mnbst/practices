package controller

import (
	"godemo/database"
	"godemo/model"
	"godemo/session"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/justinas/nosurf"
)

var Users users = users{}

type users struct{}

func (u *users) Top(c *gin.Context) {
	q := c.Query("q")
	if q != "" {
		Sort(c)
		return
	}

	user := session.GetCurrentUser(c.Request)
	todo, done := model.GetList(user.ID)

	c.HTML(http.StatusOK, "index.tpl", gin.H{
		"user":       user,
		"todo":       todo,
		"count":      len(todo),
		"done":       done,
		"count_done": len(done),
	})
}

//部分検索結果
func Sort(c *gin.Context) {
	text := strings.TrimSpace(c.Query("q"))
	user := session.GetCurrentUser(c.Request)

	todo, done := model.SortList(user.ID, text)

	c.HTML(http.StatusFound, "index.tpl", gin.H{
		"user":       user,
		"todo":       todo,
		"count":      len(todo),
		"done":       done,
		"count_done": len(done),
	})
}

func (u *users) Login(c *gin.Context) {
	csrfToken := nosurf.Token(c.Request)
	c.HTML(http.StatusOK, "user_form.tpl", gin.H{
		"csrfToken": csrfToken,
	})
}

func (u *users) Register(c *gin.Context) {
	csrfToken := nosurf.Token(c.Request)
	c.HTML(http.StatusOK, "user_form.tpl", gin.H{
		"new":       true,
		"csrfToken": csrfToken,
	})
}

func (u *users) Create(c *gin.Context) {
	email := strings.TrimSpace(c.PostForm("email"))
	password := strings.TrimSpace(c.PostForm("password"))
	if email == "" || password == "" {
		c.Redirect(http.StatusMovedPermanently, "/register")
		return
	}
	user := model.User{
		Email:    email,
		Password: model.PasswordHash(password),
	}
	db := database.GetDB()
	db.Table("users").Create(&user)

	c.Redirect(http.StatusMovedPermanently, "/login")
}

func (u *users) Authenticate(c *gin.Context) {
	user := model.User{
		Email:    c.PostForm("email"),
		Password: c.PostForm("password"),
	}

	id, err := user.Auth()
	if err != nil {
		c.Redirect(http.StatusMovedPermanently, "/login")
		return
	}

	s := session.GetSession(c.Request)
	s.Values["userId"] = id
	session.Save(c.Request, c.Writer)

	c.Redirect(http.StatusMovedPermanently, "/")
}

func (u *users) Logout(c *gin.Context) {
	session.Destroy(c.Request, c.Writer)
	c.Redirect(http.StatusFound, "/")
}

func (u *users) Deluser(c *gin.Context) {
	var user model.User
	var todo model.Todo

	id := session.GetCurrentUser(c.Request).ID
	db := database.GetDB()

	if err := db.Model(&user).Where("id = ?", id).Delete(&user).Error; err != nil {
		panic(err)
	} else {
		if err = db.Model(&todo).Where("user_id = ?", id).Delete(&todo).Error; err != nil {
			panic(err)
		}
	}

	session.Destroy(c.Request, c.Writer)
	c.Redirect(http.StatusFound, "/")
}
