package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"log"
	"net/http"
	"strings"
)

const (
	endPointUrl = "https://coding-challenge.xeptore.me/verify"
	secretCode  = "mZjMs7-Ci3wqXaFtI5FdhEqAb8Z8YkeYOOmmorinEHVf0bZHn_DCnM7oItT"
)

type requester struct {
	url    string
	method string
	secret string
	body   interface{}
}

type VerifyRequest struct {
	Lastname  string `json:"lastName"`
	Firstname string `json:"firstName"`
}

type externalRequestBody struct {
	Name string `json:"name"`
}

func handleVerify(w http.ResponseWriter, r *http.Request) {
	var request VerifyRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		responseJson(w, err)
		return
	}
	fmt.Println(request)
	err = request.validation()
	if err != nil {
		errorResponse := map[string]any{"ok": false}
		responseJson(w, errorResponse)
		return
	}
	resClient, err := requestClient(requester{
		url:    endPointUrl,
		secret: secretCode,
		method: http.MethodPost,
		body: externalRequestBody{
			Name: concatStringPlusOperation(request.Firstname, request.Lastname),
		},
	})
	if err != nil {
		// Return general error response
		errorResponse := map[string]any{"ok": false}
		responseJson(w, errorResponse)
		return
	}

	switch resClient[3] {
	case 431:
		// Return error response with "You are not allowed"
		errorResponse := map[string]any{"error": "You are not allowed"}
		responseJson(w, errorResponse)
		return
	case 360:
		message := fmt.Sprintf("Request ID %s is usable since %s", resClient[4], resClient[5])
		successResponse := map[string]any{"message": message}
		responseJson(w, successResponse)
		return
	case 200:
		// Return error response with the value of the "message" key
		errorResponse := map[string]any{"error": resClient[5]}
		responseJson(w, errorResponse)
		return
	default:
		// Return general error response
		errorResponse := map[string]any{"ok": false}
		responseJson(w, errorResponse)
		return
	}
}

func responseJson(w http.ResponseWriter, body any) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(body)
}

func (v VerifyRequest) validation() error {
	if v.Lastname == "" || v.Firstname == "" {
		return errors.New("lastname and firstname dose not empty")
	}
	return nil
}

func concatWhitFmt(firstname, lastname string) string {
	return fmt.Sprintf("%v %v", firstname, lastname)
}

func concatStringJoin(strs ...string) string {
	return strings.Join(strs, " ")
}

func concatStringPlusOperation(firstname, lastname string) string {
	return firstname + " " + lastname
}

func requestClient(requester requester) ([]any, error) {
	payload, err := json.Marshal(requester.body)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	req, err := http.NewRequest(requester.method, requester.url, bytes.NewBuffer(payload))
	if err != nil {
		log.Println(err)
		return nil, err
	}
	req.Header.Add("Secret", requester.secret)
	client := &http.Client{
		Timeout: 160,
	}

	res, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	var body []byte
	body, err = io.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	var data []any
	err = json.Unmarshal(body, &data)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return data, nil
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/verify", handleVerify).Methods(http.MethodPost)
	log.Fatal(http.ListenAndServe("127.0.0.1:9009", router))
}
