package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
)

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "ウェルカムトゥーホー<br />ムページ！")
	fmt.Println("Endpoint Hit : homepage")
}

type Article struct {
	Id      string `json:"Id"`
	Title   string `json:"Title"`
	Desc    string `json:"desc"`
	Content string `json:"content"`
}

func returnSingleArticle(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit : returnSingleArticle")

	vars := mux.Vars(r)
	key := vars["id"]
	// fmt.Fprint(w, "Key: "+key)
	for _, article := range articles {
		if article.Id == key {
			json.NewEncoder(w).Encode(article)
		}
	}
}

var articles []Article

func returnAllArticles(w http.ResponseWriter, r *http.Request) {
	// json形式でarticlesをレスポンスとして書き込む
	fmt.Println("Endpoint Hit : returnAllArticles")
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(articles)
}

func createNewArticle(w http.ResponseWriter, r *http.Request) {
	// readerである(readcloser)httpリクエストのボディをまるっと読み込み
	reqbody, _ := ioutil.ReadAll(r.Body)

	// レスポンスライターと標準出力に書いてみる log出力の方がよいきがする
	// fmt.Fprintf(w, "%+v", string(reqbody))
	// fmt.Printf("%+v", string(reqbody))

	// まるっと読んだやつを[]byteからunmarshalizeして構造体に戻す
	var art Article
	json.Unmarshal(reqbody, &art)

	// articlesに追加したら/artcle/{id}でアクセスするもよし、allで読んでもよし
	// idがダブってても弾かないから同じidのarticleがarticlesに存在しうる
	articles = append(articles, art)
	json.NewEncoder(w).Encode(art)
}
func deleteArticle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	for i, article := range articles {
		if article.Id == id {
			// delete all article id from articles
			articles = append(articles[:i], articles[i+1:]...)
		}
	}

}

func updateArticle(w http.ResponseWriter, r *http.Request) {
	reqbody, _ := ioutil.ReadAll(r.Body)
	var newArt Article
	json.Unmarshal(reqbody, &newArt)
	for i, art := range articles {
		if art.Id == newArt.Id {
			articles[i] = newArt
		}
	}
}
func handleReq() {
	// http.HandleFunc("/", homePage)
	// http.HandleFunc("/articles", returnAllArticles)

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", homePage)
	router.HandleFunc("/all", returnAllArticles)
	router.Handle("/articles", http.RedirectHandler("/all", http.StatusMovedPermanently))

	router.HandleFunc("/article", createNewArticle).Methods("POST")
	router.HandleFunc("/article", updateArticle).Methods("PUT")

	router.HandleFunc("/article/{id}", deleteArticle).Methods("DELETE")
	router.HandleFunc("/article/{id}", returnSingleArticle)

	log.Fatal(http.ListenAndServe(":8080", router))
}

func main() {
	fmt.Println("Rest API v2.0 - mux router")
	articles = []Article{
		Article{Id: "1", Title: "Hello", Desc: "説明", Content: "中身"},
		Article{Id: "2", Title: "Hello2", Desc: "2の説明", Content: "2の中身"},
	}
	handleReq()
}
