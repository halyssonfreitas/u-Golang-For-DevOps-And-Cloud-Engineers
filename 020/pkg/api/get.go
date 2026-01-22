package api

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

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

func (a api) DoGetRequest(requestURL string) (Response, error) {

	resp, err := a.Client.Get(requestURL)
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

	if !json.Valid(body) {
		return nil, RequestError{
			HTTPCode: resp.StatusCode,
			Body:     string(body),
			Err:      fmt.Sprintf("No valid JSON returned"),
		}
	}

	var page Page

	err = json.Unmarshal(body, &page)
	if err != nil {
		return nil, RequestError{
			HTTPCode: resp.StatusCode,
			Body:     string(body),
			Err:      fmt.Sprintf("Page unmarshal error: %s", err),
		}
	}

	prettyJSON, _ := json.MarshalIndent(page, "", "  ")

	fmt.Println(string(prettyJSON))

	switch page.Page {
	case "occurrence":
		var occurrence Occurrence

		err = json.Unmarshal(body, &occurrence)
		if err != nil {
			return nil, RequestError{
				HTTPCode: resp.StatusCode,
				Body:     string(body),
				Err:      fmt.Sprintf("Ocurrences unmarshal error: %s", err),
			}
		}
		return occurrence, nil
	case "words":
		var words Words

		err = json.Unmarshal(body, &words)
		if err != nil {
			return nil, RequestError{
				HTTPCode: resp.StatusCode,
				Body:     string(body),
				Err:      fmt.Sprintf("Words unmarshal error: %s", err),
			}
		}
		return words, nil
	}

	return nil, nil
}
