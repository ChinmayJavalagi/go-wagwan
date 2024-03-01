package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	dotenv "github.com/joho/godotenv"
	openai "github.com/sashabaranov/go-openai"
)

type QuestionRequest struct {
	Question string `json:"question"`
}

type OpenApiResponse struct {
	ApiCompletion string `json:"api_completion"`
}

type AddRequest struct {
	A string `json:"a"`
	B string `json:"b"`
}

type AddResponse struct {
	Total int `json:"total"`
}

type FrequencyRequest struct {
	String string `json:"string"`
}

type FrequencyResponse map[string]int

func main() {
	e := dotenv.Load()
	if e != nil {
		log.Fatal("Error loading .env file")
	}

	//endpoints
	http.HandleFunc("/open_api_completion", handleOpenApiEndPoint)
	http.HandleFunc("/add", handleAdd)
	http.HandleFunc("/total", hadleTotal)

	//web host
	e2 := http.ListenAndServe(":8080", nil)

	if e2 != nil {
		log.Fatal(e2)
	}
}

// Task 3
func handleOpenApiEndPoint(w http.ResponseWriter, r *http.Request) {
	log.Println("Endpoint for OpenApi Completion has been called...")

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusForbidden)
	}

	var reqBody QuestionRequest

	e := json.NewDecoder(r.Body).Decode(&reqBody)
	if e != nil {
		http.Error(w, e.Error(), http.StatusBadRequest)
		return
	}

	msg := reqBody.Question
	apiKey := os.Getenv("openapi")
	response := OpenApiResponse{
		ApiCompletion: openApiMessage(apiKey, msg),
	}

	jsonResponse, e := json.Marshal(response)
	if e != nil {
		http.Error(w, e.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}

// Task 4
func handleAdd(w http.ResponseWriter, r *http.Request) {
	log.Println("Endpoint for Addition has been called..")

	if r.Method != http.MethodPost {
		http.Error(w, "Wrong Method", 404)
	}

	var reqBody AddRequest

	e := json.NewDecoder(r.Body).Decode(&reqBody)
	if e != nil {
		http.Error(w, e.Error(), http.StatusBadRequest)
		return
	}

	response := AddResponse{
		Total: addTwoString(reqBody.A, reqBody.B),
	}

	jsonResponse, e := json.Marshal(response)
	if e != nil {
		http.Error(w, e.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}

// Task 5
func hadleTotal(w http.ResponseWriter, r *http.Request) {
	log.Println("Frequency count endpoint has been called...")
	if r.Method != http.MethodPost {
		http.Error(w, "Wrong Method", 404)
	}

	var reqBody FrequencyRequest

	e := json.NewDecoder(r.Body).Decode(&reqBody)

	if e != nil {
		http.Error(w, e.Error(), http.StatusBadRequest)
	}

	x := reqBody.String
	mp := make(map[rune]int)
	for _, v := range x {
		mp[v] += 1
	}

	response := make(FrequencyResponse)

	for key, value := range mp {
		response[string(key)] = value
	}

	jsonResponse, e := json.Marshal(response)
	if e != nil {
		http.Error(w, e.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)

}


//helper function for addition 
func addTwoString(a, b string) int {
	first, e1 := strconv.Atoi(a)
	second, e2 := strconv.Atoi(b)

	if e1 != nil && e2 != nil {
		panic("Please enter a valid number")
	}
	result := first + second
	return result
}

//helper function chat completion
func openApiMessage(secretKey, content string) string {
	client := openai.NewClient(secretKey)
	resp, e := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: content,
				},
			},
		},
	)

	if e != nil {
		fmt.Printf("Error in ChatCompletion: %v\n", e)
		return "Please fix the error"
	}
	fmt.Println(resp.Choices[0].Message.Content)
	return resp.Choices[0].Message.Content
}
