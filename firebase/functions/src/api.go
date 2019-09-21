package functions

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"strconv"
	"sync"

	"cloud.google.com/go/firestore"
)

var ctx context.Context
var client *firestore.Client

func init() {
	var err error
	ctx = context.Background()
	client, err = firestore.NewClient(ctx, os.Getenv("GCP_PROJECT"))
	if err != nil {
		panic(err)
	}
}

type TextTokenized struct {
	Words    []string `firestore:"word"`
	Meanings []string `firestore:"imi"`
}

type SubCaption struct {
	Index         int           `firestore:"index"`
	StartTime     int           `firestore:"start_time"`
	EndTime       int           `firestore:"end_time"`
	Text          string        `firestore:"text"`
	TextTokenized TextTokenized `firestore:"textTokenized"`
}

type Caption struct {
	Subcaps []SubCaption `firestore:"video_caption"`
}

type ResultContent struct {
	Index      int
	VideoID    string
	Text       string
	StartTime  string
	EndTime    string
	VideoTitle string
}

type Result struct {
	Meaning        string
	ResultContents []ResultContent
}

type Capdoc struct {
	VideoID    string           `firestore:"caption_id"`
	Title      string           `firestore:"caption_title"`
	Caption    []SubCaption     `firestore:"video_caption"`
	Appeartime map[string][]int `firestore:"word_element_appear_time"`
}

type Meaning struct {
	Meaning string `firestore:"word_imi"`
}

func usecase(w http.ResponseWriter, r *http.Request) {
	wg := new(sync.WaitGroup)
	c := make(chan string)
	word := r.URL.Query().Get("word")
	if word == "" {
		return
	}
	wg.Add(1)
	go func(word string) {
		var imi Meaning
		meaning := ""
		defer wg.Done()
		rawdata, err := client.Collection("dictionary").Doc(word).Get(ctx)
		if err != nil {
			panic(err)
		}
		if err = rawdata.DataTo(&imi); err != nil {
			meaning = "\"" + word + "\"は辞書にありません"
		} else {
			meaning = imi.Meaning
		}
		c <- meaning
	}(word)
	var result []ResultContent
	capdocs, err := client.Collection("captions").Documents(ctx).GetAll()
	if err != nil {
		panic(err)
	}
	for _, rowdoc := range capdocs {
		var r ResultContent
		var capdoc Capdoc
		if err = rowdoc.DataTo(&capdoc); err != nil {
			panic(err)
		}
		if _, ok := capdoc.Appeartime[word]; ok {
			for _, i := range capdoc.Appeartime[word] {
				times := caliculateMilsec([]int{capdoc.Caption[i].StartTime, capdoc.Caption[i].EndTime})
				r = ResultContent{
					Index:      i,
					VideoID:    capdoc.VideoID,
					Text:       capdoc.Caption[i].Text,
					StartTime:  times[0],
					EndTime:    times[1],
					VideoTitle: capdoc.Title,
				}
				result = append(result, r)
			}
		} else {
			continue
		}
	}
	wordimi := <-c
	wg.Wait()
	result2 := Result{
		Meaning:        wordimi,
		ResultContents: result,
	}
	res, err := json.Marshal(result2)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(res)
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
