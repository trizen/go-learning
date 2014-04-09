package main

import (
	"fmt"
	"net/http"
	"os"
)

func getContent(url string, cs chan string) {
	resp, err := http.Get(url)

	if err != nil {
		fmt.Fprintf(os.Stderr, "[%s]: %s\n", err, url)
	}

	defer resp.Body.Close()
	cs <- url
}

func main() {

	urls := []string{
		"http://facebook.com",
		"http://google.com",
		"http://trizen.go.ro",
		"http://github.com",
		"http://trizenx.blogspot.com",
	}

	cs := make(chan string, len(urls))

	fmt.Println("Start!")

	for _, url := range urls {
		go getContent(url, cs)
	}

	fmt.Println("All done!")

	for i := range urls {
		fmt.Printf("[%d] %s\n", i, <-cs)
	}
}
