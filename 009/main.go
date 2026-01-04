// buid : go build http-get
// execute 1 :
//   - ./http-get
//   - echo $?
//
// execute 2 :
//   - ./http-get argument1
//   - echo $?
package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
)

func main() {
	args := os.Args

	if len(args) < 2 {
		fmt.Println("Usage: ./http-get <url>")
		os.Exit(1)
	}

	if _, err := url.ParseRequestURI(args[1]); err != nil {
		fmt.Printf("URL is in invalida format: %s\n", err)
		os.Exit(1)
	}

	resp, err := http.Get(args[1])
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("HTTP Status Code: %d\nBody: %s\n", resp.StatusCode, body)
}
