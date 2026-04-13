package engine

import (
	"backend/src/utils"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type ollamaChatResponse struct {
	Model           string    `json:"model"`
	CreatedAt       time.Time `json:"created_at"`
	Message         Message   `json:"message"`
	Done            bool      `json:"done"`
	DoneReason      string    `json:"done_reason"`
	TotalDuration   int64     `json:"total_duration"`
	PromptEvalCount int       `json:"prompt_eval_count"`
	EvalCount       int       `json:"eval_count"`
}

type ChatRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
	Stream   bool      `json:"stream"`
	Options  struct {
		Num_ctx int
	}
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

func llmService(systemPrompt string, userPrompt string, timeout float64) (ollamaChatResponse, error) {
	timeout = min(timeout, 1800)
	requestData := ChatRequest{
		Model: "qwen3:4b",
		Messages: []Message{
			{Role: "system", Content: systemPrompt},
			{Role: "user", Content: userPrompt},
		},
		Stream: false,
		Options: struct{ Num_ctx int }{
			Num_ctx: 4096,
		},
	}

	client := &http.Client{
		Timeout: time.Duration(timeout) * time.Second,
	}

	jsonData, err := json.Marshal(requestData)
	if err != nil {
		fmt.Println("Error marshaling json :", err)
		return ollamaChatResponse{}, err
	}

	fmt.Printf(" sent to llm : \n %+v \n", requestData)

	resp, err := client.Post(
		"http://localhost:11434/api/chat",
		"application/json",
		bytes.NewReader(jsonData),
	)
	if err != nil {
		fmt.Println("Error making request:", err)
		return ollamaChatResponse{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {

		fmt.Println("Llm Request rejected by ollamam", err)
		return ollamaChatResponse{}, err
	}

	var llmResponse ollamaChatResponse
	err = utils.JsonResponseRead(resp, &llmResponse)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return ollamaChatResponse{}, err
	}

	fmt.Printf(" received from llm : \n %+v \n", llmResponse)

	return llmResponse, nil
}
