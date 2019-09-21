package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"

	_ "github.com/mattn/go-sqlite3"
)

var Db1 *gorm.DB
var Db2 *gorm.DB
var Db3 *sql.DB
var Db4 *sql.DB

//Thumbnail return top page thubmnails
type Thumbnail struct {
	Video_href  string
	Video_img   string
	Video_title string
	Video_time  string
}

//Dic return words in dic page
type Dic struct {
	Word_ini string
	Word     string
	Word_imi string
}

type TextTokenized struct {
	Word []string `json:"word"`
	Imi  []string `json:"imi"`
}

type Caption struct {
	Index         int           `json:"index, string"`
	Start_time    int           `json:"start_time, string"`
	End_time      int           `json:"end_time, string"`
	Text          string        `json:"text"`
	TextTokenized TextTokenized `json:"textTokenized"`
}

type Video struct {
	Video_href               string
	Video_img                string
	Video_title              string
	Video_caption            string
	Word_element_appear_time string
}

func getthumbnails(Db1 *gorm.DB) []Thumbnail {
	var thumbnail []Thumbnail
	Db1.Table("video").Find(&thumbnail)
	return thumbnail
}

func readdicDb(ini string, Db2 *gorm.DB) []Dic {
	var dic []Dic
	Db2.Order("word asc").Where("word_ini=?", ini).Find(&dic)
	return dic
}

func getcaption(href string, Db4 *sql.DB) []Caption {
	var cap string
	err := Db4.QueryRow("SELECT video_caption FROM video where video_href=?", href).Scan(&cap)
	if err != nil {
		panic(err)
	}
	var caption []Caption
	err = json.Unmarshal([]byte(cap), &caption)
	if err != nil {
		panic(err)
	}
	return caption
}

type CaptionText struct {
	Index         int           `json:"index, string"`
	StartTime     int           `json:"start_time, string"`
	EndTime       int           `json:"end_time, string"`
	Text          string        `json:"text"`
	TextTokenized TextTokenized `json:"textTokenized"`
}

type ResultContent struct {
	Index       int
	Video_href  string
	Text        string
	StartTime   string
	EndTime     string
	Video_title string
}

func getVideoContents(word string, Db1 *gorm.DB, Db3 *sql.DB) ([]ResultContent, string) {
	wg := new(sync.WaitGroup)
	c := make(chan string)

	wg.Add(1)
	go func(word string) {
		defer wg.Done()
		row := Db3.QueryRow("SELECT word_imi FROM dic where word=?", word)
		var imi string
		err := row.Scan(&imi)
		if err != nil {
			imi = "\"" + word + "\"は辞書にありません"
		}
		c <- imi
	}(word)

	var allvideo []Video
	var result []ResultContent
	Db1.Find(&allvideo)
	for _, v := range allvideo {
		var wordappeartime map[string][]int //mapping video_caption
		var caption []CaptionText
		var r ResultContent

		err := json.Unmarshal([]byte(v.Word_element_appear_time), &wordappeartime)
		if err != nil {
			panic(err)
		}

		if _, ok := wordappeartime[word]; ok {
			err = json.Unmarshal([]byte(v.Video_caption), &caption)
			if err != nil {
				panic(err)
			}
			for _, i := range wordappeartime[word] {
				times := caliculateMilsec([]int{caption[i].StartTime, caption[i].EndTime})
				r = ResultContent{
					Index:       caption[i].Index,
					Video_href:  v.Video_href,
					Text:        caption[i].Text,
					StartTime:   times[0],
					EndTime:     times[1],
					Video_title: v.Video_title,
				}
				result = append(result, r)
			}
		} else {
			continue
		}
	}
	imi := <-c
	wg.Wait()
	return result, imi
}

func caliculateMilsec(times []int) []string {
	var result []string
	for _, i := range times {
		min := i / 1000 / 60
		sec := i / 1000 % 60
		r := strconv.Itoa(min) + ":" + strconv.Itoa(sec)
		result = append(result, r)
	}
	return result
}

type Yid struct {
	Id string `json:"id"`
}

func main() {

	Db1, err := gorm.Open("sqlite3", "./db/video.db")
	if err != nil {
		panic(err)
	}
	defer Db1.Close()
	Db1.SingularTable(true)
	Db2, err := gorm.Open("sqlite3", "./db/dictionary.db")
	if err != nil {
		panic(err)
	}
	defer Db2.Close()
	Db2.SingularTable(true)
	Db3, err := sql.Open("sqlite3", "./db/dictionary.db")
	if err != nil {
		panic(err)
	}
	defer Db3.Close()
	Db4, err := sql.Open("sqlite3", "./db/video.db")
	if err != nil {
		panic(err)
	}
	defer Db4.Close()

	router := gin.Default()
	router.LoadHTMLGlob("templates/*.gotpl")

	router.Static("static/", "./static")

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.gotpl", gin.H{
			"result": getthumbnails(Db1),
		})
	})

	router.GET("/video/:href", func(c *gin.Context) {
		href := c.Param("href")
		c.HTML(http.StatusOK, "video.gotpl", gin.H{
			"id":      href,
			"caption": getcaption(href, Db4),
		})
	})

	router.GET("dic/:initial", func(c *gin.Context) {
		ini := c.Param("initial")
		c.HTML(http.StatusOK, "dic1.gotpl", gin.H{
			"rows": readdicDb(ini, Db2),
		})
	})

	router.GET("dic/:initial/:word", func(c *gin.Context) {
		ini := c.Param("initial")
		word := c.Param("word")
		yid := c.Query("id")
		if yid != "" {
			var cap string
			var caption []Caption
			err = Db4.QueryRow("SELECT video_caption FROM video where video_href=?", yid).Scan(&cap)
			if err != nil {
			}
			err = json.Unmarshal([]byte(cap), &caption)
			c.JSON(200, caption)
			return
		}
		contents, imi := getVideoContents(word, Db1, Db3)
		c.HTML(http.StatusOK, "dic2.gotpl", gin.H{
			"word":     word,
			"contents": contents,
			"ini":      ini,
			"imi":      imi,
		})
	})

	router.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{"code": "PAGE_NOT_FOUND", "message": "Page not found"})
	})

	router.Run("127.0.0.1:8081")
}
