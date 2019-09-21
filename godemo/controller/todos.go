package controller

import (
	"fmt"
	"godemo/database"
	"godemo/model"
	"godemo/session"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type todos struct{}

// TODO関連の処理を行うコントローラ
var Todos todos

// TODOを登録する
func (u *todos) Create(c *gin.Context) {
	var todo model.Todo
	var exist model.Todo

	if c.BindJSON(&todo) != nil {
		fmt.Println("binding error")
		return
	}

	db := database.GetDB()
	user := session.GetCurrentUser(c.Request)
	todo.UserID = user.ID
	if err := db.Model(&todo).Where(map[string]interface{}{"user_id": user.ID, "title": todo.Title, "text": todo.Title}).First(&exist).Error; err != nil {
		if err = db.Save(&todo).Error; err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
		}
	}
}

// TODOを完了ずみする
func (u *todos) Done(c *gin.Context) {
	var todo model.Todo
	//var todo1 model.Todo

	if c.BindJSON(&todo) != nil {
		return
	}

	user := session.GetCurrentUser(c.Request)
	id := user.ID
	title := todo.Title
	text := todo.Text
	var todo1 model.Todo

	db := database.GetDB()

	if err := db.Model(&todo1).Where(map[string]interface{}{"user_id": id, "title": title, "text": text}).First(&todo1).Error; err != nil {
		panic(err)
	} else {
		if err := db.Model(&todo1).Updates(map[string]interface{}{"completed": true, "completed_at": time.Now()}).Error; err != nil {
			panic(err)
		}
	}
}

//todo削除
func (u *todos) Delete(c *gin.Context) {
	var todo model.Todo
	//var todo1 model.Todo

	if c.BindJSON(&todo) != nil {
		return
	}

	user := session.GetCurrentUser(c.Request)
	id := user.ID
	title := todo.Title
	text := todo.Text
	var todo1 model.Todo

	db := database.GetDB()

	if err := db.Model(&todo1).Where(map[string]interface{}{"user_id": id, "title": title, "text": text}).Delete(&todo1).Error; err != nil {
		panic(err)
	}
}
