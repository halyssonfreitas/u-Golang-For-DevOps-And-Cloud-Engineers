// buid : go build http-get
// execute 1 :
//   - ./http-get
//   - echo $?
//
// execute 2 :
//   - # run test-server
//   - curl http://localhost:8080/words?input=halysson # 2x
//   - ./http-get http://localhost:8080/words
//   - ./http-get http://localhost:8080/ocurrence
//   - echo $?
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
)

// {"page":"words","input":"marieli","words":["marieli"]}

type Page struct {
	Page string `json:"page"`
}

type Words struct {
	Page
	Input string   `json:"input"`
	Words []string `json:"words"`
}

type Occurrence struct {
	Page
	Words map[string]int `json:"words"`
}

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

	if resp.StatusCode != 200 {
		fmt.Printf("Invalid output (HTTP Code %d): %s\n", resp.StatusCode, body)
		os.Exit(1)
	}

	var page Page

	err = json.Unmarshal(body, &page)
	if err != nil {
		log.Fatal(err)
	}

	prettyJSON, _ := json.MarshalIndent(page, "", "  ")

	fmt.Println(string(prettyJSON))

	switch page.Page {
	case "occurrence":
		var occurrence Occurrence

		err = json.Unmarshal(body, &occurrence)
		if err != nil {
			log.Fatal(err)
		}

		prettyJSON, _ := json.MarshalIndent(occurrence, "", "  ")

		fmt.Println(string(prettyJSON))
	case "words":
		var words Words

		err = json.Unmarshal(body, &words)
		if err != nil {
			log.Fatal(err)
		}

		prettyJSON, _ := json.MarshalIndent(words, "", "  ")

		fmt.Println(string(prettyJSON))
	default:
		fmt.Println("Page not found!")
		println()
	}

}
