package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"strconv"
	"time"
)

type Article struct {
	Id int `json:id`
	Title string `json:"title"`
	Subtitle string `json:"subtitle"`
	Content string `json:"content"`
	Timestamp int64 `json:"timestamp"`
}

// let's declare a global Articles array
// that we can then populate in our main function
// to simulate a database
type Articles []Article

var articles = Articles{
	Article{Id: 0, Title: "Hello0", Subtitle: "Sub0", Content: "Article0", Timestamp: makeTimestamp()},
	Article{Id: 1, Title: "Hello1", Subtitle: "Sub1", Content: "Article1", Timestamp: makeTimestamp()},
	Article{Id: 2, Title: "Hello2", Subtitle: "Sub2", Content: "Article2", Timestamp: makeTimestamp()},
	Article{Id: 3, Title: "Hello3", Subtitle: "Sub3", Content: "Article3", Timestamp: makeTimestamp()},
	Article{Id: 4, Title: "Hello4", Subtitle: "Sub4", Content: "Article4", Timestamp: makeTimestamp()},
	Article{Id: 5, Title: "Hello5", Subtitle: "Sub5", Content: "Article5", Timestamp: makeTimestamp()},
	Article{Id: 6, Title: "Hello6", Subtitle: "Sub6", Content: "Article6", Timestamp: makeTimestamp()},
	Article{Id: 7, Title: "Hello7", Subtitle: "Sub7", Content: "Article7", Timestamp: makeTimestamp()},
}

var totalArticles = 8

func makeTimestamp() int64 {
	return time.Now().UnixNano()
}

func returnPaginationResult(offset int, limit int) Articles {
	response := Articles{}

	fmt.Println(offset, limit)

	if offset >= totalArticles {
		return response
	}

	var count = 0

	for id := offset; id < offset + limit; id++ {
		if id >= totalArticles {
			break
		}
		fmt.Println(id)
		response = append(response, articles[id])
		count++
	}

	return response
}

func homePage(w http.ResponseWriter, r *http.Request){
	fmt.Fprintf(w, "Welcome to the HomePage!")
	fmt.Println("Endpoint Hit: homePage")
}

func returnAllArticles(w http.ResponseWriter, r *http.Request){
	fmt.Println("Endpoint Hit: articles")
	switch r.Method {
	case http.MethodGet:
		// GET
		type Response struct {
			Response Articles `json:"response"`
		}

		var response Response

		fmt.Println("GET")
		offsets, ok1 := r.URL.Query()["o"]
		limits, ok2 := r.URL.Query()["l"]

		if !ok1 || len(offsets[0]) < 1 || !ok2 || len(limits[0]) < 1 {
			response.Response = articles
			json.NewEncoder(w).Encode(response)
			return
		}
		offset, err1 := strconv.Atoi(offsets[0])
		limit, err2 := strconv.Atoi(limits[0])
		if err1 != nil || err2 != nil {
			response.Response = nil
			json.NewEncoder(w).Encode(response)
			return
		}

		response.Response = returnPaginationResult(offset, limit)

		json.NewEncoder(w).Encode(response)
	case http.MethodPost:
		// POST
		decoder := json.NewDecoder(r.Body)

		var newArticle Article
		err := decoder.Decode(&newArticle)

		if err != nil {
			fmt.Println("error")
		}
		newArticle.Id = totalArticles
		newArticle.Timestamp = makeTimestamp()
		totalArticles++
		articles = append(articles, newArticle)
		fmt.Println("Success!")
	}
}

func returnArticleById(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: articleById")
	resultArticle := Article{}
	q := r.URL.Path[len("/articles/"):]
	id, err := strconv.Atoi(q)
	if err != nil {
		json.NewEncoder(w).Encode(nil)
		return
	}
	fmt.Println(id)
	if id >= totalArticles {
		json.NewEncoder(w).Encode(nil)
		return
	}
	resultArticle = articles[id]
	json.NewEncoder(w).Encode(resultArticle)
}

func queryArticle(w http.ResponseWriter, r *http.Request)  {
	fmt.Println("Endpoint Hit: queryArticle")
	keys, val := r.URL.Query()["q"]

	if !val || len(keys[0]) < 1 {
		log.Println("Invalid query")
		return
	}

	// Query()["key"] will return an array of items,
	// we only want the single item.
	query := keys[0]

	var searchResult = Articles{}

	for id := range articles{
		var title = articles[id].Subtitle
		var subtitle = articles[id].Title
		var content = articles[id].Content
		if strings.Contains(title, query) ||
			strings.Contains(subtitle, query) ||
			strings.Contains(content, query){
			searchResult = append(searchResult, articles[id])
		}
	}

	json.NewEncoder(w).Encode(searchResult)

}

func handleRequests() {
	http.HandleFunc("/", homePage)
	http.HandleFunc("/articles", returnAllArticles)
	http.HandleFunc("/articles/", returnArticleById)
	http.HandleFunc("/articles/search", queryArticle)
	log.Fatal(http.ListenAndServe(":5000", nil))
}

func main() {
	handleRequests()
}
