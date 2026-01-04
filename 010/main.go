// buid : go build http-get
// execute 1 :
//   - ./http-get
//   - echo $?
//
// execute 2 :
//   - ./http-get localhost:8080/words (inicialize o test-server)
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

type Words struct {
	Page  string   `json:"page"`
	Input string   `json:"input"`
	Words []string `json:"words"`
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

	var words Words

	err = json.Unmarshal(body, &words)
	if err != nil {
		log.Fatal(err)
	}

	prettyJSON, err := json.MarshalIndent(words, "", "  ")
	
	fmt.Println(string(prettyJSON))
}
