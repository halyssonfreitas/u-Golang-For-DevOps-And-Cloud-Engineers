package api

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"
)

type MockClient struct {
	ResponseOutput *http.Response
}

func (m *MockClient) Get(url string) (resp *http.Response, err error) {
	return m.ResponseOutput, nil
}

func TestDoGetResquet(t *testing.T) {
	words := Words{
		Page: Page{
			Page: "words",
		},
		Input: "Pluma",
		Words: []string{"Pluma", "Vaquinha", "Pluma"},
	}

	wordsBytes, err := json.Marshal(words)
	if err != nil {
		t.Errorf("marshal error: %s", err)
	}

	apiInstance := api{
		Options: Options{},
		Client: &MockClient{
			ResponseOutput: &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(bytes.NewReader(wordsBytes)),
			},
		},
	}
	response, err := apiInstance.DoGetRequest("http://localhost/words")
	if err != nil {
		t.Errorf("DoGetRequest error: %s", err)
	}
	if response == nil {
		t.Fatalf("response is empty", )
	}
	if response.GetResponse() != strings.Join([]string{"Pluma", "Vaquinha", "Pluma"}, ", "){
		t.Errorf("Unexpected response: %s", response.GetResponse())
	}
}
