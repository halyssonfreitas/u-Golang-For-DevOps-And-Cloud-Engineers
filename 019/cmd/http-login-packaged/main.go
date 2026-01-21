// buid : go build http-get
// execute 1 :
//   - ./http-get
//   - echo $?
//
// execute 2 :
//   - # run test-server --password xyz
//   - curl http://localhost:8080/words?input=halysson # 2x
//   - ./http-get --help
//   - ./http-get http://localhost:8080/words
//   - ./http-get --url http://localhost:8080/words --password xyz
//   - ./http-get --url http://localhost:8080/occurrence --password xyz
//   - echo $?
package main

import (
	"flag"
	"fmt"
	"net/url"
	"os"

	api "github.com/halyssonfreitas/u-Golang-For-DevOps-And-Cloud-Engineers/http-login-packaged/pkg/api"
)

// {"page":"words","input":"marieli","words":["marieli"]}

func main() {
	var (
		requestURL string
		password   string
		parsedURL  *url.URL
		err        error
	)

	flag.StringVar(&requestURL, "url", "", "url to accesss")
	flag.StringVar(&password, "password", "", "use a password to accesss our api")

	flag.Parse()

	if parsedURL, err = url.ParseRequestURI(requestURL); err != nil {
		fmt.Printf("URL is in invalida format: %s\n", err)
		flag.Usage()
		os.Exit(1)
	}

	apiInstance := api.New(api.Options{
		Password: password,
		LoginURL: parsedURL.Scheme + "://" + parsedURL.Host + "/login",
	})

	res, err := apiInstance.DoGetRequest(parsedURL.String())
	if err != nil {
		if requestErr, ok := err.(api.RequestError); ok {
			fmt.Printf("Error: %s (HTTP Code: %d, Body: %s)\n", requestErr.Err, requestErr.HTTPCode, requestErr.Body)
			os.Exit(1)
		}
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}

	if res == nil {
		fmt.Printf("No response\n")
		os.Exit(1)
	}

	fmt.Printf("Response: %s\n", res.GetResponse())
}
