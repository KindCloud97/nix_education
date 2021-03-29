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

var url string = "https://jsonplaceholder.typicode.com"

func GetComments(postId int, c chan Comments, wg *sync.WaitGroup) {
	time.Sleep(time.Second)
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
		fmt.Println(jsonErr)
	}

	for elem := range dataComments {
		c <- dataComments[elem]
	}
	wg.Done()
}

func main() {
	/*dbP, err := sql.Open("mysql",
		"root:password@/post")
	if err != nil {
		fmt.Print(err)
	}
	defer dbP.Close()*/
	//получаем посты
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
		fmt.Println(jsonErr)
	}
	////получаем комменты и записываем посты в бд
	c := make(chan Comments, 100)
	var wg sync.WaitGroup
	for _, elem := range dataPosts {
		wg.Add(1)
		go GetComments(elem.Id, c, &wg)
		go func(post Posts) {
			time.Sleep(time.Millisecond)
			fmt.Println("[WriteToDB--Posts]\t", post)
		}(elem)
	}
	//закрываем канал после записи всех комментов
	go func() { //?????????
		wg.Wait()
		close(c)
	}()
	//записываем комменты в бд
	for elem := range c {
		go func(comment Comments) {
			time.Sleep(time.Millisecond)
			fmt.Println("[WriteToDB--Comments]\t", comment)
		}(elem)
	}

}

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
