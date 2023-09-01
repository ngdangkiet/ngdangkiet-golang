package main

import (
	"fmt"
	"github.com/goccy/go-json"
	"github.com/joho/godotenv"
	"io"
	"log"
	"net/http"
	"os"
)

type Quotes struct {
	Quotes []Quote
}

type Quote struct {
	_Id     string `json:"_id"`
	Content string `json:"content"`
	Author  string `json:"author"`
}

func getJson(url string, data interface{}) error {
	res, err := http.Get(url)
	if err != nil {
		return err
	}

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(res.Body)

	return json.NewDecoder(res.Body).Decode(data)
}

func getQuote(url string) (string, string) {
	var quotes Quotes

	err := getJson(url, &quotes.Quotes)

	if err != nil {
		fmt.Printf("%v", err.Error())
		return "", ""
	}

	fmt.Println(quotes)

	return quotes.Quotes[0].Content, quotes.Quotes[0].Author
}

func generate(quote string, author string) {
	path := "README.md"

	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)

	if err != nil {
		fmt.Printf("%v", err.Error())
		return
	}

	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	data := []byte(fmt.Sprintf("_**%s**_\n\n%s", quote, author))

	_, err = file.Write(data)
	if err != nil {
		fmt.Printf("%v", err.Error())
		return
	}

	fmt.Println("Write file successfully!")
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	url := os.Getenv("URL")

	quote, author := getQuote(url)

	if quote != "" && author != "" {
		generate(quote, author)
	}
}
