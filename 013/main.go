// buid : go build http-get
// execute 1 :
//   - ./http-get
//   - echo $?
//
// execute 2 :
//   - # run test-server
//   - curl http://localhost:8080/words?input=halysson # 2x
//   - ./http-get http://localhost:8080/words
//   - ./http-get http://localhost:8080/occurrence
//   - echo $?
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
)

// {"page":"words","input":"marieli","words":["marieli"]}

type Response interface {
	GetResponse() string
}

type Page struct {
	Page string `json:"page"`
}

type Words struct {
	Page
	Input string   `json:"input"`
	Words []string `json:"words"`
}

func (w Words) GetResponse() string {
	return fmt.Sprintf("%s", strings.Join(w.Words, ", "))
}

type Occurrence struct {
	Page
	Words map[string]int `json:"words"`
}

func (o Occurrence) GetResponse() string {
	out := []string{}
	for word, ocurrence := range o.Words {
		out = append(out, fmt.Sprintf("%s (%d)", word, ocurrence))
	}
	return fmt.Sprintf("%s", strings.Join(out, ", "))
}

func main() {
	args := os.Args

	if len(args) < 2 {
		fmt.Println("Usage: ./http-get <url>")
		os.Exit(1)
	}

	res, err := doRequest(args[1])
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}

	if res == nil {
		fmt.Printf("No response\n")
		os.Exit(1)
	}

	fmt.Printf("Response: %s\n", res.GetResponse())
}

func doRequest(requestURL string) (Response, error) {
	if _, err := url.ParseRequestURI(requestURL); err != nil {
		return nil, fmt.Errorf("URL is in invalida format: %s\n", err)
	}

	resp, err := http.Get(requestURL)
	if err != nil {
		return nil, fmt.Errorf("Http get error: %s", err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return nil, fmt.Errorf("ReadAll error: %s", err)
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Invalid output (HTTP Code %d): %s\n", resp.StatusCode, body)
	}

	var page Page

	err = json.Unmarshal(body, &page)
	if err != nil {
		return nil, fmt.Errorf("Unmarshal error: %s", err)
	}

	prettyJSON, _ := json.MarshalIndent(page, "", "  ")

	fmt.Println(string(prettyJSON))

	switch page.Page {
	case "occurrence":
		var occurrence Occurrence

		err = json.Unmarshal(body, &occurrence)
		if err != nil {
			return nil, fmt.Errorf("Unmarshal error: %s", err)
		}
		return occurrence, nil
	case "words":
		var words Words

		err = json.Unmarshal(body, &words)
		if err != nil {
			return nil, fmt.Errorf("Unmarshal error: %s", err)
		}
		return words, nil
	}

	return nil, nil
}
