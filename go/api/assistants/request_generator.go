package assistants

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func RequestGenerator(method string, endpoint string, pathParam string, token string, body map[string]string) (*http.Request, error) {
	host := "http://localhost:8080/api"
	url := fmt.Sprintf("%v/%v", host, endpoint)

	if pathParam != "" {
		url = fmt.Sprintf("%v/%v", url, pathParam)
	}

	var jsonBody []byte
	if body != nil {
		var err error
		jsonBody, err = json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to parse request body, error - %w", err)
		}
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("error while creating request %w", err)
	}

	req.Header.Set("accept", "application/json")

	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	return req, nil
}
