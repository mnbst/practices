package model

import (
	"godemo/database"
	"sync"
	"time"
)

// TODOを保持するモデル
type Todo struct {
	ID          uint      `gorm:"primary_key" sql:"not null"`
	UserID      uint      `sql:"not null"`
	Title       string    `sql:"not null" json:"title"`
	Text        string    `sql:"not null" json:"text"`
	Completed   bool      `json:"completed"`
	CompletedAt time.Time `sql:"not null"`
	CreatedAt   time.Time `sql:"not null"`
}

//todoリストを返却
func GetList(userID uint) ([]Todo, []Todo) {
	wg := new(sync.WaitGroup)
	ch := make(chan []Todo)

	var todo []Todo
	var done []Todo
	var doneResult []Todo

	if userID != 0 {
		db := database.GetDB()
		wg.Add(1)

		go func(userID uint) {
			defer wg.Done()
			db := database.GetDB()
			db.Model(&done).Where(map[string]interface{}{"user_id": userID, "completed": true}).Find(&done)
			ch <- done
		}(userID)

		db.Model(&todo).Where(map[string]interface{}{"user_id": userID, "completed": false}).Find(&todo)

		doneResult = <-ch
		wg.Wait()
	}

	return todo, doneResult
}

func SortList(userID uint, input string) ([]Todo, []Todo) {
	wg := new(sync.WaitGroup)
	ch1 := make(chan []Todo)
	var todo []Todo
	var todo1 []Todo

	if userID != 0 {
		db := database.GetDB()
		wg.Add(1)

		go func(userID uint, input string) {
			defer wg.Done()
			if userID != 0 {
				db := database.GetDB()
				db.Model(&todo1).Where("user_id = $1 AND completed = $2 AND (title LIKE $3 OR text LIKE $3)", userID, true, "%"+input+"%").Find(&todo1)
			}
			ch1 <- todo1
		}(userID, input)

		db.Model(&todo).Where("user_id = $1 AND completed = $2 AND (title LIKE $3 OR text LIKE $3)", userID, false, "%"+input+"%").Find(&todo)
	}

	done := <-ch1
	wg.Wait()

	return todo, done
}
