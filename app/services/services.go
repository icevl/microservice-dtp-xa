package services

import (
	"app/types"
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"
)

type Service interface {
	Create(resultCh chan<- types.Result)
	Commit(resultCh chan<- types.Result)
	Rollback()
}

func HttpRequest(url string, payload interface{}) (int, error) {
	jsonData, err := json.Marshal(payload)

	if err != nil {
		return http.StatusInternalServerError, errors.New("Payload conversion error")
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonData))

	if err != nil {
		return http.StatusInternalServerError, errors.New("Error creating HTTP request")
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: time.Second * 10}
	response, err := client.Do(req)
	if err != nil {
		return http.StatusInternalServerError, errors.New("Error executing HTTP request")
	}

	defer response.Body.Close()

	var result map[string]interface{}
	var body, _ = io.ReadAll(response.Body)

	err = json.Unmarshal([]byte(body), &result)

	if err != nil {
		return http.StatusInternalServerError, errors.New("Error converting response to JSON")
	}

	if result["success"] == false {
		return response.StatusCode, errors.New(result["message"].(string))
	}

	return http.StatusAccepted, nil
}
