package main

import (
	//"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"sync"
	"time"

	//_ "github.com/go-sql-driver/mysql"
)

const url string = "https://jsonplaceholder.typicode.com"

type Posts struct {
	UserId int    `json:"userId"`
	Id     int    `json:"id"`
	Title  string `json:"title"`
	Body   string `json:"body"`
}

type Comments struct {
	PostId int    `json:"postId"`
	Id     int    `json:"id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	Body   string `json:"body"`
}

func GetPosts() []Posts {
	var dataPosts = []Posts{}

	resp, err := http.Get(url + "/posts?userId=7")
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}

	jsonErr := json.Unmarshal(body, &dataPosts)
	if jsonErr != nil {
		panic(jsonErr)
	}
	return dataPosts
}

func GetComments(postId int, c chan Comments, wg *sync.WaitGroup) {
	var dataComments = []Comments{}

	resp, err := http.Get(url + "/comments?postId=" + strconv.Itoa(postId))
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}

	jsonErr := json.Unmarshal(body, &dataComments)
	if jsonErr != nil {
		panic(jsonErr)
	}

	for elem := range dataComments {
		c <- dataComments[elem]
	}
	time.Sleep(time.Second * 2)
	wg.Done()
}

func main() {
	var dataPosts = GetPosts()
	chanComm := make(chan Comments)

	var wg sync.WaitGroup
	for _, elem := range dataPosts {
		wg.Add(1)

		go GetComments(elem.Id, chanComm, &wg)

		go func(post Posts) {
			time.Sleep(time.Millisecond)
			fmt.Println("[WriteToDB--Posts]\t", post.Id)
		}(elem)
	}

	go func() {
		wg.Wait()
		close(chanComm)
	}()

	for elem := range chanComm {
		go func(comment Comments) {
			time.Sleep(time.Millisecond)
			fmt.Println("[WriteToDB--Comments]\t", comment.PostId)
		}(elem)
	}
}
