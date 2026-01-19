package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type LoginRequest struct {
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `token:"password"`
}

func doLoginRequest(client http.Client, requestUrl, password string) (string, error) {
	loginRequest := LoginRequest{
		Password: password,
	}

	body, err := json.Marshal(loginRequest)
	if err != nil {
		return "", fmt.Errorf("Marshal error: %s", err)
	}

	resp, err := client.Post(requestUrl, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return "", fmt.Errorf("Http POST error: %s", err)
	}

	defer resp.Body.Close()

	resBody, err := io.ReadAll(resp.Body)

	if err != nil {
		return "", fmt.Errorf("ReadAll error: %s", err)
	}

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("Invalid output (HTTP Code %d): %s\n", resp.StatusCode, resBody)
	}

	if !json.Valid(resBody) {
		return "", RequestError{
			HTTPCode: resp.StatusCode,
			Body:     string(resBody),
			Err:      fmt.Sprintf("No valid JSON returned"),
		}
	}

	var loginResponse LoginResponse

	err = json.Unmarshal(resBody, &loginResponse)
	if err != nil {
		return "", RequestError{
			HTTPCode: resp.StatusCode,
			Body:     string(resBody),
			Err:      fmt.Sprintf("Login response unmarshal error: %s", err),
		}
	}

	if loginResponse.Token == "" {
		return "", RequestError{
			HTTPCode: resp.StatusCode,
			Body:     string(resBody),
			Err:      fmt.Sprintf("Empty token replied"),
		}
	}

	return loginResponse.Token, nil
}
