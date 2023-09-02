package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"github.com/joho/godotenv"
	"io"
	"log"
	"net/http"
	"os"
	"time"
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

func getQuote(url string) Quote {
	var quotes Quotes

	err := getJson(url, &quotes.Quotes)

	if err != nil {
		fmt.Printf("%v", err.Error())
		return Quote{}
	}

	fmt.Println(quotes)

	return quotes.Quotes[0]
}

func generate(quote Quote) {
	path := "README.md"

	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)

	if err != nil {
		fmt.Printf("%v", err.Error())
		return
	}

	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	data := []byte(fmt.Sprintf("_**%s**_\n\n%s", quote.Content, quote.Author))

	_, err = file.Write(data)
	if err != nil {
		fmt.Printf("%v", err.Error())
		return
	}

	fmt.Println("Write file successfully!")
}

func getQuoteRandom(c *gin.Context, quote Quote) {
	c.JSONP(http.StatusOK, quote)
}

func restApi(port string, quote Quote) {
	router := gin.Default()

	router.GET("/quotes/random", func(c *gin.Context) {
		getQuoteRandom(c, quote)
	})

	// goroutine
	go func() {
		if err := router.Run(fmt.Sprintf(":%s", port)); err != nil {
			fmt.Println("Server error:", err)
		}
	}()

	duration := 10 * time.Second
	fmt.Printf("Server will automatically close after %s...\n", duration)
	<-time.After(duration)
	fmt.Println("Server is shutting down!")
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	url := os.Getenv("URL")
	port := os.Getenv("PORT")
	quote := getQuote(url)

	if quote != (Quote{}) {
		generate(quote)
		restApi(port, quote)
	}
}
